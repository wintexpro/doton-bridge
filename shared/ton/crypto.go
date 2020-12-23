// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"encoding/base64"
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

func ImportTonKeysFromFile(keystorepath, contractsPath, filename string, password []byte) (string, error) {
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

	kp.SetAddress(DeriveRelayerAddress(contractsPath, client.KeyPair{
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

func ImportTonPrivKey(keystorepath, contractsPath, key string, password []byte) (string, error) {
	if password == nil {
		password = keystore.GetPassword("Enter password to encrypt keystore file:")
	}

	kp, err := ed25519.NewKeypairFromSeed(key)
	if err != nil {
		return "", fmt.Errorf("could not generate ed25519 keypair from given string: %w", err)
	}

	kp.SetAddress(DeriveRelayerAddress(contractsPath, client.KeyPair{
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

func LoadAbi(path, name string) client.Abi {
	content, err := ioutil.ReadFile(path + "/" + name + ".abi.json")
	if err != nil {
		panic(err)
	}
	abi := client.Abi{Type: client.ContractAbiType}
	if err = json.Unmarshal(content, &abi.Value); err != nil {
		panic(err)
	}

	return abi
}

func LoadTvc(path, name string) string {
	content, err := ioutil.ReadFile(path + "/" + name + ".tvc")
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(content)
}

type RelayerInitialData struct {
	AccessControllerAddress string
	MyPublicKey             []byte
	MyInitState             []byte
	BridgeAddress           string
}

func DeriveRelayerAddress(contractsPath string, keys client.KeyPair) string {
	RelayerABI := LoadAbi(contractsPath, "Relayer")
	RelayerTVC := LoadTvc(contractsPath, "Relayer")

	conn, err := NewClient()
	if err != nil {
		panic(err)
	}

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: keys,
	}

	deploySet := client.DeploySet{
		Tvc:         RelayerTVC,
		WorkchainID: null.NewInt32(0, true),
	}

	callSet := client.CallSet{
		FunctionName: "constructor",
		Input: map[string]interface{}{
			// ignore this args
			"_accessControllerAddress": null.StringFrom("0:164d61e6cad0597545cb8ab98ecfdb2a29e0cc55d484daece02c63d8511e9a5f"),
			"_myPublicKey":             null.StringFrom("0x" + keys.Public),
			"_myInitState":             null.StringFrom(RelayerTVC),
			"_bridgeAddress":           null.StringFrom("0:164d61e6cad0597545cb8ab98ecfdb2a29e0cc55d484daece02c63d8511e9a5f"),
		},
	}

	params := client.ParamsOfEncodeMessage{
		Abi:       RelayerABI,
		DeploySet: &deploySet,
		CallSet:   &callSet,
		Signer:    signer,
	}

	res, err := conn.AbiEncodeMessage(&params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nAddress: %s \n\n", res.Address)

	return res.Address
}
