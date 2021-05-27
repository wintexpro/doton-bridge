// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"os"
	"testing"
	"time"

	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/gtank/merlin"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/keystore"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var TestLogger = log15.New("chain", "test")

func TestSubChain(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// From:     "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
	// From:         "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty",

	var prepareChainConn func(string, string) *Chain

	prepareChainConn = func(name string, from string) *Chain {
		cfg := &core.ChainConfig{
			Id:           1,
			Name:         name,
			Endpoint:     "ws://localhost:9944",
			From:         from,
			KeystorePath: dir + "/../../keys",
			Opts: map[string]string{
				"useExtendedCall": "true",
				"startBlock":      "1",
			},
			LatestBlock: false,
			FreshStart:  true,
		}

		pswdStr := "123456"

		os.Setenv(keystore.EnvPassword, pswdStr)

		logger := log.Root().New("test", cfg.Name)

		sysErr := make(chan error)
		chain, err := InitializeChain(cfg, logger, sysErr, nil)
		if err != nil {
			t.Fatal(err)
		}

		TestLogger.SetHandler(log15.LvlFilterHandler(log15.LvlError, TestLogger.GetHandler()))

		router := core.NewRouter(TestLogger)
		chain.SetRouter(router)

		err = chain.Start()
		if err != nil {
			t.Fatal(err)
		}

		return chain
	}

	chainAlice := prepareChainConn("Alice", "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")
	chainBob := prepareChainConn("Bob", "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")

	time.Sleep(time.Second * 66)

	m := msg.Message{
		Source:       msg.ChainId(2),
		Destination:  msg.ChainId(1),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(1),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			chainAlice.cfg.From,
			types.Text("hello ton"),
		},
	}

	resolvedAlice := chainAlice.writer.ResolveMessage(m)
	resolvedBob := chainBob.writer.ResolveMessage(m)

	if !resolvedBob || !resolvedAlice {
		t.Fatal("The proposal is not resolved.", "Resolved by Alice:", resolvedAlice, "Resolved by Bob", resolvedBob)
	}
}

func NewSigningContext(context, msg []byte) *merlin.Transcript {
	t := merlin.NewTranscript("SigningContext")
	t.AppendMessage([]byte(""), context)
	t.AppendMessage([]byte("sign-bytes"), msg)
	return t
}
