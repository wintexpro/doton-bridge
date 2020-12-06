// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"fmt"
	"math/big"

	"github.com/ByKeks/chainbridge-utils/core"
	metrics "github.com/ByKeks/chainbridge-utils/metrics/types"
	"github.com/ByKeks/chainbridge-utils/msg"
	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
)

var _ core.Writer = &writer{}

type sendDataInput struct {
	Destination string   `json:"destination"`
	Bounce      bool     `json:"bounce"`
	Data        string   `json:"data"`
	Value       *big.Int `json:"value"`
}

type writer struct {
	cfg     Config
	conn    Connection
	log     log15.Logger
	stop    <-chan int
	sysErr  chan<- error // Reports fatal error to core
	metrics *metrics.ChainMetrics
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		cfg:     *cfg,
		conn:    conn,
		log:     log,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
	}
}

func (w *writer) start() error {
	w.log.Debug("Starting ton writer...")
	return nil
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
