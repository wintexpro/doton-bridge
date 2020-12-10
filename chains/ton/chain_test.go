// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
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
		Endpoint:       "http://net.ton.dev",
		From:           "ebc77aae202a4f12237e10892f4fe0e44f8fb3dfc07008dcc12b37f8f70c1149",
		Insecure:       false,
		KeystorePath:   dir + "/../../keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts: map[string]string{
			"abiPath":    dir + "/mocks/contracts",
			"startBlock": "1751455",
			"sender":     "0:dee8cdbf9937431376dd7ab7ee93367c14c62acc24d1d558cdd01186cf45704d",
			"receiver":   "0:c0c4627877c03b66d81d4d037dc696a322d63b6e14bea1e6fd39955734af6f5b",
		},
		LatestBlock: false,
	}

	pswdStr := "123456"

	os.Setenv(keystore.EnvPassword, pswdStr)
	importTonPrivKey(cfg.KeystorePath, "action glow era all liquid critic achieve lawsuit era anger loud slight", []byte(pswdStr))

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

	time.Sleep(time.Second * 20)

	chain.conn.Client().Close()
}

func importTonPrivKey(keystorepath, key string, password []byte) (string, error) {
	if password == nil {
		password = keystore.GetPassword("Enter password to encrypt keystore file:")
	}

	kp, err := ed25519.NewKeypairFromSeed(key)
	if err != nil {
		return "", fmt.Errorf("could not generate ed25519 keypair from given string: %w", err)
	}

	fp, err := filepath.Abs(keystorepath + "/" + kp.PublicKey() + ".key")
	if err != nil {
		return "", fmt.Errorf("invalid filepath: %w", err)
	}

	file, err := os.OpenFile(filepath.Clean(fp), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("Unable to Open File: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Error("import private key: could not close keystore file")
		}
	}()

	err = keystore.EncryptAndWriteToFile(file, kp, password)
	if err != nil {
		return "", fmt.Errorf("could not write key to file: %w", err)
	}

	log.Info("private key imported", "public key", kp.PublicKey(), "file", fp)
	return fp, nil
}
