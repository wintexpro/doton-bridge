// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ChainSafe/ChainBridge/tonbindings"
	. "github.com/ChainSafe/ChainBridge/tonbindings"
	"github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	metrics "github.com/wintexpro/chainbridge-utils/metrics/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var _ core.Writer = &writer{}

const RelayerContractKey = "Relayer"
const RootTokenContractKey = "RootToken"
const EpochVoteControllerContractKey = "EpochVoteController"
const MessageHandlerContractKey = "MessageHandler"

type writer struct {
	cfg             Config
	conn            Connection
	log             log15.Logger
	kp              *ed25519.Keypair
	stop            <-chan int
	sysErr          chan<- error // Reports fatal error to core
	metrics         *metrics.ChainMetrics
	abi             map[string]client.Abi
	relayer         *RelayerContract
	epochController *EpochControllerContract
	epochContract   *Epoch
	epoch           int64
	queue           []*msg.Message
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, kp *ed25519.Keypair, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	abi := make(map[string]client.Abi)

	workchainID, err := strconv.ParseInt(cfg.workchainID, 10, 32)
	if err != nil {
		panic(err)
	}

	signer := client.Signer{
		EnumTypeValue: client.KeysSigner{
			Keys: client.KeyPair{
				Public: kp.PublicKey(),
				Secret: kp.SecretKey(),
			},
		},
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      &signer,
		WorkchainID: null.Int32From(int32(workchainID)),
	}

	rootTokenContract := RootTokenContract{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}
	epochControllerContract := EpochController{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}
	epochContract := Epoch{Ctx: ctx}

	abiRootTokenContract, err := rootTokenContract.Abi()
	if err != nil {
		panic(err)
	}

	abiMessageHandler, err := messageHandlerContract.Abi()
	if err != nil {
		panic(err)
	}

	abi[RootTokenContractKey] = *abiRootTokenContract
	abi[MessageHandlerContractKey] = *abiMessageHandler

	relayer, err := relayerContract.New(cfg.from)
	if err != nil {
		panic(err)
	}

	epochController, err := epochControllerContract.New(cfg.contracts[EpochVoteControllerOpt])
	if err != nil {
		panic(err)
	}

	return &writer{
		cfg:             *cfg,
		conn:            conn,
		log:             log,
		kp:              kp,
		stop:            stop,
		sysErr:          sysErr,
		metrics:         m,
		abi:             abi,
		relayer:         relayer,
		epochController: epochController,
		epochContract:   &epochContract,
		epoch:           0,
	}
}

func (w *writer) PublicRandomness() (string, error) {
	publicRandomnessMap, err := w.epochController.PublicRandomness().Call()
	if err != nil {
		return "", err
	}

	return publicRandomnessMap.(map[string]interface{})["publicRandomness"].(string), nil
}

func (w *writer) MessageCallback(event *client.ProcessingEvent) {
	// w.log.Debug("MessageID: %s", event)
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())

	var data string

	switch m.Type {
	case msg.FungibleTransfer:
		amount := new(big.Int).SetBytes(m.Payload[0].([]byte))
		input, err := json.Marshal(map[string]interface{}{
			"to":     string(m.Payload[1].([]byte)),
			"tokens": amount.String(),
		})
		if err != nil {
			w.log.Error("failed to construct FungibleTransfer data", "chainId", m.Destination, "error", err)
			return false
		}

		paramsOfEncodeMessageBody := client.ParamsOfEncodeMessageBody{
			Abi:        w.abi[RootTokenContractKey],
			Signer:     *w.relayer.Ctx.Signer,
			IsInternal: true,
			CallSet: client.CallSet{
				FunctionName: "mint",
				Input:        input,
			},
		}

		resultOfEncodeMessageBody, err := w.conn.Client().AbiEncodeMessageBody(&paramsOfEncodeMessageBody)
		if err != nil {
			w.log.Error("failed to construct encode FungibleTransfer data", "chainId", m.Destination, "error", err)
			return false
		}

		data = resultOfEncodeMessageBody.Body
	case SimpleMessageTransfer:
		messageType := "0x" + hex.EncodeToString(m.ResourceId[:])
		dataStr := "0x" + hex.EncodeToString([]byte(m.Payload[1].(types.Text)))

		input, err := json.Marshal(map[string]interface{}{
			"chainId":     m.Destination,
			"nonce":       m.DepositNonce,
			"messageType": messageType,
			"data":        dataStr,
		})

		if err != nil {
			w.log.Error("failed to construct SimpleMessageTransfer data", "chainId", m.Destination, "error", err)
			return false
		}
		paramsOfEncodeMessageBody := client.ParamsOfEncodeMessageBody{
			Abi:    w.abi[MessageHandlerContractKey],
			Signer: *w.relayer.Ctx.Signer,
			CallSet: client.CallSet{
				FunctionName: "receiveMessage",
				Input:        input,
			},
		}

		resultOfEncodeMessageBody, err := w.conn.Client().AbiEncodeMessageBody(&paramsOfEncodeMessageBody)
		if err != nil {
			w.log.Error("failed to encode SimpleMessageTransfer data", "chainId", m.Destination, "error", err)
			return false
		}

		data = resultOfEncodeMessageBody.Body
	}

	messageType := "0x" + hex.EncodeToString(m.ResourceId[:])

	currentEpochNumberRef, err := w.GetEpochNumber()
	if err != nil {
		w.log.Error("failed to get epoch address", "chainId", fmt.Sprint(m.Destination), "error", err)
		return false
	}

	currentEpochNumber := currentEpochNumberRef

	if currentEpochNumber > 1 {
		currentEpochNumber = currentEpochNumber - 1
	}

	currentEpochNumberAsStr := strconv.FormatInt(currentEpochNumber, 10)

	epochAddressMap, err := w.epochController.GetEpochAddress(currentEpochNumberAsStr).Call()
	if err != nil {
		w.log.Error("failed to parse epoch address", "chainId", m.Destination, "error", err)
		return false
	}
	epochAddress := epochAddressMap.(map[string]interface{})["epoch"].(string)

	epoch, err := w.epochContract.New(
		epochAddress,
		&EpochInitVars{
			Number:                currentEpochNumberAsStr,
			VoteControllerAddress: w.epochController.Address,
		},
	)
	if err != nil {
		w.log.Error("failed to initialize epochContract", "chainId", m.Destination, "error", err)
		return false
	}

	resMap, err := epoch.IsChoosen(w.relayer.Address).Call()
	if err != nil {
		if err.Error() == "unexpected end of JSON input" {
			w.queue = append(w.queue, &m)
			w.log.Error("failed to get IsChoosen", "chainId", m.Destination, "relayer", w.relayer.Address, "epoch", epoch.Address, "epochNumber", currentEpochNumberAsStr, "error", errors.New("Epoch is not ready now"))
		} else {
			w.log.Error("failed to get IsChoosen", "chainId", m.Destination, "relayer", w.relayer.Address, "epoch", epoch.Address, "error", err)
		}
		return false
	}
	isChoosen := resMap.(map[string]interface{})["value0"].(bool)

	if !isChoosen {
		w.log.Info("Your relayer is not active now")
		return false
	}

	//FIXME: Check is proposal valid
	shardBlockID, err := w.relayer.VoteThroughBridge(strconv.FormatInt(currentEpochNumber, 10), "1", fmt.Sprint(m.Destination), messageType, fmt.Sprint(m.DepositNonce), data).Send(w.MessageCallback)
	if err != nil {
		w.log.Error("failed to construct proposal", "chainId", m.Destination, "error", err)
		return false
	}

	w.log.Info("Attemping to send proposal", "ShardBlockID", shardBlockID, "Epoch", currentEpochNumber)

	return true
}

func (w *writer) GetEpochAddress(currentEpochNumber string) (string, error) {
	epochAddressMap, err := w.epochController.GetEpochAddress(currentEpochNumber).Call()
	if err != nil {
		return "", err
	}
	return epochAddressMap.(map[string]interface{})["epoch"].(string), nil
}

func (w *writer) GetPublicRandomness() (string, error) {
	return w.PublicRandomness()
}

func (w *writer) GetEpochNumber() (int64, error) {
	currentEpochNumberMap, err := w.epochController.CurrentEpochNumber().Call()
	if err != nil {
		return 0, err
	}
	currentEpochNumberStr := currentEpochNumberMap.(map[string]interface{})["currentEpochNumber"].(string)

	currentEpochNumber, err := strconv.ParseInt(currentEpochNumberStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return currentEpochNumber, nil
}

func (w *writer) SendVrfPublicKey(vrfkp *Keypair) error {
	currentEpochNumber, err := w.GetEpochNumber()
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to get epoch number", "chainId", w.cfg.id, "error", err)
		return err
	}

	epochNumberAsStr := strconv.FormatInt(currentEpochNumber, 10)

	epochAddress, err := w.GetEpochAddress(epochNumberAsStr)
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to get epoch address", "chainId", w.cfg.id, "error", err)
		return err
	}

	publicRandomness, err := w.GetPublicRandomness()
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to get public randomness", "chainId", w.cfg.id, "error", err)
		return err
	}

	decodedPR, err := hex.DecodeString(publicRandomness[2:])
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to decode public randomness", "chainId", w.cfg.id, "error", err)
		return err
	}

	public := hex.EncodeToString(vrfkp.Public)

	sign, err := vrfkp.Sign(decodedPR)
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to sign public randomness", "chainId", w.cfg.id, "error", err)
		return err
	}

	_, err = w.relayer.SignUpForEpoch(
		epochAddress,
		fmt.Sprintf("0x%x", sign[:32]),
		fmt.Sprintf("0x%x", sign[32:]),
		fmt.Sprintf("0x%s", public),
	).Send(func(event *client.ProcessingEvent) {})
	if err != nil {
		w.log.Error("SendVrfPublicKey: failed to signup to epoch", "chainId", w.cfg.id, "epochNumber", currentEpochNumber, "epochAddress", epochAddress, "error", err)
		return err
	}

	return nil
}

func (w *writer) CheckEpoch() {
	ticker := time.NewTicker(6 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				currentEpochNumber, err := w.GetEpochNumber()
				if err != nil {
					w.log.Error("failed to get epoch number", "chainId", w.cfg.id, "epochNumber", currentEpochNumber, "error", err)
					continue
				}

				epochNumberForVoting := strconv.FormatInt(currentEpochNumber-1, 10)
				epochNumberForRegistration := strconv.FormatInt(currentEpochNumber, 10)

				epochAddressForVoting, _ := w.GetEpochAddress(epochNumberForVoting)
				epochAddressForRegistration, _ := w.GetEpochAddress(epochNumberForRegistration)

				epochForVoting, err := w.epochContract.New(epochAddressForVoting, &tonbindings.EpochInitVars{
					Number:                epochNumberForVoting,
					VoteControllerAddress: w.epochController.Address,
				})
				if err != nil {
					w.log.Error("failed to initialize epoch contract", "chainId", w.cfg.id, "error", err)
					continue
				}

				epochForRegistration, err := w.epochContract.New(epochAddressForRegistration, &tonbindings.EpochInitVars{
					Number:                epochNumberForRegistration,
					VoteControllerAddress: w.epochController.Address,
				})
				if err != nil {
					w.log.Error("failed to initialize epoch contract", "chainId", w.cfg.id, "error", err)
					continue
				}

				epochForVotingEndsAtMap, err := epochForVoting.FirstEraEndsAt().Call()
				if err != nil {
					w.log.Error("failed to get FirstEraEndsAt state", "chainId", w.cfg.id, "error", err)
					continue
				}
				epochForVotingEndsAt := epochForVotingEndsAtMap.(map[string]interface{})["firstEraEndsAt"].(string)

				epochForVotingEndsAtInt, err := strconv.ParseInt(epochForVotingEndsAt, 10, 32)
				if err != nil {
					w.log.Error("failed to parse endsAt to int", "chainId", w.cfg.id, "error", err)
					continue
				}

				epochForVotingSecondEraEndsAtMap, err := epochForVoting.SecondEraEndsAt().Call()
				if err != nil {
					w.log.Error("failed to get FirstEraEndsAt state", "chainId", w.cfg.id, "error", err)
					continue
				}
				epochForVotingSecondEraEndsAt := epochForVotingSecondEraEndsAtMap.(map[string]interface{})["secondEraEndsAt"].(string)

				epochForVotingSecondEraEndsAtInt, err := strconv.ParseInt(epochForVotingSecondEraEndsAt, 10, 32)
				if err != nil {
					w.log.Error("failed to parse endsAt to int", "chainId", w.cfg.id, "error", err)
					continue
				}

				epochForRegistrationEndsAtMap, err := epochForRegistration.FirstEraEndsAt().Call()
				if err != nil {
					w.log.Error("failed to get FirstEraEndsAt state", "chainId", w.cfg.id, "error", err)
					continue
				}
				epochForRegistrationEndsAt := epochForRegistrationEndsAtMap.(map[string]interface{})["firstEraEndsAt"].(string)

				epochForRegistrationEndsAtInt, err := strconv.ParseInt(epochForRegistrationEndsAt, 10, 32)
				if err != nil {
					w.log.Error("failed to parse endsAt to int", "chainId", w.cfg.id, "error", err)
					continue
				}

				epochForVotingMap, err := epochForVoting.IsChoosen(w.relayer.Address).Call()
				if err != nil {
					if err.Error() == "unexpected end of JSON input" {
						w.log.Error("failed to get IsChoosen", "chainId", w.cfg.id, "relayer", w.relayer.Address, "epoch", epochForVoting.Address, "epochNumber", epochNumberForVoting, "error", errors.New("Epoch is not ready now"))
					} else {
						w.log.Error("failed to get IsChoosen", "chainId", w.cfg.id, "relayer", w.relayer.Address, "epoch", epochForVoting.Address, "error", err)
					}
					continue
				}
				epochForVotingIsChoosen := epochForVotingMap.(map[string]interface{})["value0"].(bool)

				w.log.Info("Epoch for voting info", "IsChoosen", epochForVotingIsChoosen, "Number", epochNumberForVoting, "FirstEraEndsAtInt", epochForVotingEndsAtInt, "FirstEraEndsAtInt DELTA", time.Now().Unix()-epochForVotingEndsAtInt, "SecondEraEndsAtInt", epochForVotingSecondEraEndsAtInt, "SecondEraEndsAtInt DELTA", time.Now().Unix()-epochForVotingSecondEraEndsAtInt)

				epochForRegistrationMap, err := epochForRegistration.IsChoosen(w.relayer.Address).Call()
				if err != nil {
					if err.Error() == "unexpected end of JSON input" {
						w.log.Error("failed to get IsChoosen", "chainId", w.cfg.id, "relayer", w.relayer.Address, "epoch", epochForVoting.Address, "epochNumber", epochNumberForVoting, "error", errors.New("Epoch is not ready now"))
					} else {
						w.log.Error("failed to get IsChoosen", "chainId", w.cfg.id, "relayer", w.relayer.Address, "epoch", epochForVoting.Address, "error", err)
					}
					continue
				}
				epochForRegistrationIsChoosen := epochForRegistrationMap.(map[string]interface{})["value0"].(bool)

				w.log.Info("Epoch for registration info", "IsChoosen", epochForRegistrationIsChoosen, "Number", epochNumberForRegistration, "FirstEndsAtInt", epochForRegistrationEndsAtInt, "FirstEndsAtInt DELTA", time.Now().Unix()-epochForRegistrationEndsAtInt)

				if w.epoch == currentEpochNumber && time.Now().Unix()-epochForVotingSecondEraEndsAtInt < 0 {
					continue
				}

				for i, m := range w.queue {
					if m == nil {
						continue
					}
					w.queue[i] = nil
					w.ResolveMessage(*m)
				}

				vrfkp, err := VrfGenerateKeypair()
				if err != nil {
					w.log.Error("failed to generate keypair", "chainId", w.cfg.id, "error", err)
				}

				err = w.SendVrfPublicKey(vrfkp)
				if err != nil {
					w.log.Error("failed to send vrf public key", "chainId", w.cfg.id, "error", err)
				}

				w.epoch = currentEpochNumber
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
