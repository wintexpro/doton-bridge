// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"testing"

	"github.com/ByKeks/chainbridge-utils/core"
	"github.com/ByKeks/chainbridge-utils/keystore"
	"github.com/ByKeks/chainbridge-utils/msg"
)

func TestChain_ListenerShutdownOnFailure(t *testing.T) {
	cfg := &core.ChainConfig{
		Id:             msg.ChainId(0),
		Name:           "alice",
		Endpoint:       "http://localhost",
		From:           "0:841288ed3b55d9cdafa806807f02a0ae0c169aa5edfe88a789a6482429756a94",
		Insecure:       false,
		KeystorePath:   "./keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts:           map[string]string{},
	}
	t.Logf("1 %v", cfg)

	keyI, err := keystore.KeypairFromAddress(cfg.From, keystore.TonChain, cfg.KeystorePath, cfg.Insecure)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("keyI: %v", keyI)

	// logger := log.Root().New("test", cfg.Name)

	// t.Logf("ChainId: %v", cfg.Id)
	// t.Logf("Name: %v", cfg.Name)
	// t.Logf("Endpoint: %v", cfg.Endpoint)
	// t.Logf("From: %v", cfg.From)
	// t.Logf("keystore type: %v", keystore.TonChain)

	// sysErr := make(chan error)
	// _, err := InitializeChain(cfg, logger, sysErr, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
