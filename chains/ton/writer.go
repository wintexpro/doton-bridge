// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"fmt"

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

const Relayer = "Relayer"

type writer struct {
	cfg     Config
	conn    Connection
	log     log15.Logger
	kp      *ed25519.Keypair
	stop    <-chan int
	sysErr  chan<- error // Reports fatal error to core
	metrics *metrics.ChainMetrics
	abi     map[string]client.Abi
	tvc     map[string]string
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, kp *ed25519.Keypair, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	abi := make(map[string]client.Abi)
	tvc := make(map[string]string)

	abi[Relayer] = LoadAbi(cfg.contractsPath, Relayer)
	tvc[Relayer] = LoadTvc(cfg.contractsPath, Relayer)

	return &writer{
		cfg:     *cfg,
		conn:    conn,
		log:     log,
		kp:      kp,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
		abi:     abi,
		tvc:     tvc,
	}
}

func MessageCallback(event *client.ProcessingEvent) {
	fmt.Printf("MessageID: %s", event.MessageID)
}

var SimpleMessageResourceIDAsSrt = "0x000000000000000000000053696d706c654d6573736167655265736f75726365"

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())

	address := null.NewString(w.cfg.from, true)

	keys := client.KeyPair{
		Public: w.kp.PublicKey(),
		Secret: w.kp.SecretKey(),
	}

	relayerABI := LoadAbi(w.cfg.contractsPath, "Relayer")

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: keys,
	}

	randomKeys, err := ed25519.GenerateKeypair()
	if err != nil {
		return false
	}

	callSet := client.CallSet{
		FunctionName: "voteThroughBridge",
		Input: map[string]interface{}{
			"choice":            uint8(1),
			"chainId":           m.Destination,
			"messageType":       null.StringFrom(SimpleMessageResourceIDAsSrt),
			"nonce":             m.DepositNonce,
			"data":              null.StringFrom("0x" + hex.EncodeToString([]byte(m.Payload[0].(types.Text)))),
			"proposalPublicKey": null.StringFrom("0x" + randomKeys.PublicKey()),
		},
	}

	params := client.ParamsOfEncodeMessage{
		Abi:     relayerABI,
		Address: address,
		CallSet: &callSet,
		Signer:  signer,
	}

	res, err := w.conn.Client().AbiEncodeMessage(&params)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.sysErr <- fmt.Errorf("failed to construct proposal (chain=%d) Error: %w", m.Destination, err)
		return false
	}

	sparams := client.ParamsOfSendMessage{
		Message: res.Message,
		Abi:     &relayerABI,
	}

	result, err := w.conn.Client().ProcessingSendMessage(&sparams, MessageCallback)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.sysErr <- fmt.Errorf("failed to send proposal (chain=%d) Error: %w", m.Destination, err)
		return false
	}

	fmt.Printf("Result: %v\n", result)

	return false
}
