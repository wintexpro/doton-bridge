// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	metrics "github.com/wintexpro/chainbridge-utils/metrics/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var _ core.Writer = &writer{}

type writer struct {
	cfg     Config
	conn    Connection
	log     log15.Logger
	kp      *ed25519.Keypair
	stop    <-chan int
	sysErr  chan<- error // Reports fatal error to core
	metrics *metrics.ChainMetrics
	abi     map[string]client.Abi
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, kp *ed25519.Keypair, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	senderABI := LoadAbi(cfg.contractsPath, "Sender")
	abi := make(map[string]client.Abi)
	abi["Sender"] = senderABI

	return &writer{
		cfg:     *cfg,
		conn:    conn,
		log:     log,
		kp:      kp,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
		abi:     abi,
	}
}

func MessageCallback(event *client.ProcessingEvent) {
	fmt.Printf("MessageID: %s", event.MessageID)
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())

	// w.abi["Sender"]

	return false
}

func (w *writer) CheckAndDeploySender() error {
	keys := client.KeyPair{
		Public: w.kp.PublicKey(),
		Secret: w.kp.SecretKey(),
	}

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: keys,
	}

	SenderABI := LoadAbi(w.cfg.contractsPath, "Sender")
	SenderTVC := LoadTvc(w.cfg.contractsPath, "Sender")

	deploySet := client.DeploySet{
		Tvc:         SenderTVC,
		WorkchainID: null.NewInt32(0, true),
	}

	callSet := client.CallSet{
		FunctionName: "constructor",
	}

	paramsOfEncodeMsg := client.ParamsOfEncodeMessage{
		Abi:       SenderABI,
		DeploySet: &deploySet,
		CallSet:   &callSet,
		Signer:    signer,
	}

	encodedMessage, err := w.conn.Client().AbiEncodeMessage(&paramsOfEncodeMsg)
	if err != nil {
		return err
	}

	paramsOfQuery := client.ParamsOfQueryCollection{
		Collection: "accounts",
		Filter:     json.RawMessage(`{ "id": {"eq": "` + encodedMessage.Address + `"} }`),
		Result:     "id balance",
	}

	res, err := w.conn.Client().NetQueryCollection(&paramsOfQuery)
	if err != nil {
		return err
	}

	if len(res.Result) == 0 {
		return errors.New("Your sender address: " + encodedMessage.Address + " does not have balance for deploy the contract code")
	}

	paramsOfSendMessage := client.ParamsOfSendMessage{
		Message: encodedMessage.Message,
		Abi:     &SenderABI,
	}

	_, err = w.conn.Client().ProcessingSendMessage(&paramsOfSendMessage, MessageCallback)

	return nil
}
