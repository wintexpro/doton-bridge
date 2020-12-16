// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prometheus/common/log"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	"github.com/wintexpro/chainbridge-utils/keystore"
)

type EncryptedTonKeystore struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

func ImportTonKeysFromFile(keystorepath, filename string, password []byte) (string, error) {
	if password == nil {
		password = keystore.GetPassword("Enter password to encrypt keystore file:")
	}

	importdata, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return "", fmt.Errorf("could not read import file: %w", err)
	}

	ksjson := new(EncryptedTonKeystore)
	err = json.Unmarshal(importdata, ksjson)

	kp := ed25519.Keypair{}
	kp.Decode([]byte(ksjson.Secret))

	kp.SetAddress(DeriveSenderAddress(client.KeyPair{
		Public: kp.PublicKey(),
		Secret: kp.SecretKey(),
	}))

	fp, err := filepath.Abs(keystorepath + "/" + kp.Address() + ".key")
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

	err = keystore.EncryptAndWriteToFile(file, &kp, password)
	if err != nil {
		return "", fmt.Errorf("could not write key to file: %w", err)
	}

	log.Info("private key imported", "public key", kp.PublicKey(), "file", fp)
	return fp, nil
}

func ImportTonPrivKey(keystorepath, key string, password []byte) (string, error) {
	if password == nil {
		password = keystore.GetPassword("Enter password to encrypt keystore file:")
	}

	kp, err := ed25519.NewKeypairFromSeed(key)
	if err != nil {
		return "", fmt.Errorf("could not generate ed25519 keypair from given string: %w", err)
	}

	kp.SetAddress(DeriveSenderAddress(client.KeyPair{
		Public: kp.PublicKey(),
		Secret: kp.SecretKey(),
	}))

	fp, err := filepath.Abs(keystorepath + "/" + kp.Address() + ".key")
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

func DeriveSenderAddress(keys client.KeyPair) string {
	conn, err := NewClient()
	if err != nil {
		panic(err)
	}

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: keys,
	}

	deploySet := client.DeploySet{
		Tvc:         SenderTVC,
		WorkchainID: null.NewInt32(0, true),
	}

	callSet := client.CallSet{
		FunctionName: "constructor",
	}

	abi := client.Abi{Type: client.ContractAbiType}
	if err = json.Unmarshal([]byte(SenderABI), &abi.Value); err != nil {
		panic(err)
	}

	params := client.ParamsOfEncodeMessage{
		Abi:       abi,
		DeploySet: &deploySet,
		CallSet:   &callSet,
		Signer:    signer,
	}

	res, err := conn.AbiEncodeMessage(&params)

	return res.Address
}
