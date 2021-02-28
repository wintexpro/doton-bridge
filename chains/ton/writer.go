// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

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
const MessageHandlerContractKey = "MessageHandler"

type writer struct {
	cfg     Config
	conn    Connection
	log     log15.Logger
	kp      *ed25519.Keypair
	stop    <-chan int
	sysErr  chan<- error // Reports fatal error to core
	metrics *metrics.ChainMetrics
	abi     map[string]client.Abi
	relayer *RelayerContract
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, kp *ed25519.Keypair, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	abi := make(map[string]client.Abi)

	workchainID, err := strconv.ParseInt(cfg.workchainID, 10, 32)
	if err != nil {
		panic(err)
	}

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: client.KeyPair{
			Public: kp.PublicKey(),
			Secret: kp.SecretKey(),
		},
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      &signer,
		WorkchainID: null.Int32From(int32(workchainID)),
	}

	rootTokenContract := RootTokenContract{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}

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

	return &writer{
		cfg:     *cfg,
		conn:    conn,
		log:     log,
		kp:      kp,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
		abi:     abi,
		relayer: relayer,
	}
}

func (w *writer) MessageCallback(event *client.ProcessingEvent) {
	w.log.Debug("MessageID: %s", event.MessageID)
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

		fmt.Printf("\n\n input: %#v \n\n", map[string]interface{}{
			"to":     string(m.Payload[1].([]byte)),
			"tokens": amount.String(),
		})

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

		fmt.Printf("\n\n input: %#v \n\n", map[string]interface{}{
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

	//FIXME: Check is proposal valid
	shardBlockID, err := w.relayer.VoteThroughBridge("1", fmt.Sprint(m.Destination), messageType, fmt.Sprint(m.DepositNonce), data).Send(w.MessageCallback)
	if err != nil {
		w.log.Error("failed to construct proposal", "chainId", m.Destination, "error", err)
		return false
	}

	w.log.Info("Attemping to send proposal", "ShardBlockID", shardBlockID)

	return true
}
