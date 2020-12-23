// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
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

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())

	return false
}
