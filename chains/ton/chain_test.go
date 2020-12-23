// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"os"
	"testing"
	"time"

	utils "github.com/ChainSafe/ChainBridge/shared/ton"
	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/keystore"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var TestLogger = log15.New("chain", "test")

func TestTonChain(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &core.ChainConfig{
		Id:             msg.ChainId(1),
		Name:           "alice",
		Endpoint:       "http://localhost",
		From:           "0:df22eba0b48020b70efa7a6e9d6360ed1dc20877250947470cc1289b14c9cc1e",
		Insecure:       false,
		KeystorePath:   dir + "/../../keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts: map[string]string{
			"contractsPath": dir + "/mocks/contracts",
			"receiver":      "0:5a1921e4c0ec016f6be00917e06adbf06af2a26acc860af7d989f9a754f2ae89",
			"startBlock":    "3",
			"workchainID":   "-1",
		},
		LatestBlock: false,
	}

	pswdStr := "123456"

	os.Setenv(keystore.EnvPassword, pswdStr)
	utils.ImportTonPrivKey(cfg.KeystorePath, cfg.Opts["contractsPath"], "action glow era all liquid critic achieve lawsuit era anger loud slight", []byte(pswdStr))

	logger := log.Root().New("test", cfg.Name)

	sysErr := make(chan error)
	chain, err := InitializeChain(cfg, logger, sysErr, nil)
	if err != nil {
		t.Fatal(err)
	}

	TestLogger.SetHandler(log15.LvlFilterHandler(log15.LvlError, TestLogger.GetHandler()))

	r := core.NewRouter(TestLogger)
	chain.SetRouter(r)

	err = chain.Start()
	if err != nil {
		t.Fatal(err)
	}

	m := msg.Message{
		Source:       msg.ChainId(1),
		Destination:  msg.ChainId(2),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(1),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			types.Text("Hello ton!"),
		},
	}

	chain.writer.ResolveMessage(m)

	time.Sleep(time.Second * 20)

	chain.conn.Client().Close()
}
