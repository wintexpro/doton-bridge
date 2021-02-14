// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	utils "github.com/ChainSafe/ChainBridge/shared/ton"
	. "github.com/ChainSafe/ChainBridge/tonbindings"
	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
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
		From:           "0:e8bc02f9e8e56c04d9248743cfed5b8c3a0d6b5171f7fcf083a0dd206f847891",
		Insecure:       false,
		KeystorePath:   dir + "/../../keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts: map[string]string{
			"contractsPath": dir + "/mocks/contracts",
			"receiver":      "0:21235dbf90a09dc61d5389545530613c1fa02f9cc7cf596988cf4a215f8e9220",
			"startBlock":    "3",
			"workchainID":   "0",
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

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: client.KeyPair{
			Public: chain.kp.PublicKey(),
			Secret: chain.kp.SecretKey(),
		},
	}

	TestLogger.SetHandler(log15.LvlFilterHandler(log15.LvlError, TestLogger.GetHandler()))

	r := core.NewRouter(TestLogger)
	chain.SetRouter(r)

	err = chain.Start()
	if err != nil {
		t.Fatal(err)
	}

	messageCallback := func(event *client.ProcessingEvent) {
		t.Logf("\n\nEventType: %v", event.Type)
		// t.Logf("\n\nShardBlockID: %v", event.ShardBlockID)
		t.Logf("\n\nMessageId: %v", event.MessageID)
		// t.Logf("\n\nevent: %#v", event)
	}

	workchainID := null.Int32From(0)

	var (
		accessControllerAddress, senderAddress, tip3HandlerAddress, tonTokenWalletAddress,
		bridgeAddress, relayerAddress, receiverAddress, bridgeVoteControllerAddress, messageHandlerAddress,
		rootTokenContractAddress string
	)

	giver, err := NewGiver(chain.conn.Client(), signer, workchainID)
	if err != nil {
		t.Fatal(err)
	}

	ctx := ContractContext{
		Conn:        chain.conn.Client(),
		Signer:      &signer,
		WorkchainID: workchainID,
	}

	proposalContract := Proposal{Ctx: ctx}
	senderContract := Sender{Ctx: ctx}
	receiverContract := Receiver{Ctx: ctx}
	accessControllerContract := AccessController{Ctx: ctx}
	bridgeVoteControllerContract := BridgeVoteController{Ctx: ctx}
	bridgeContract := Bridge{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}
	tip3HandlerContract := Tip3Handler{Ctx: ctx}
	rootTokenContract := RootTokenContract{Ctx: ctx}
	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		t.Fatal(err)
	}

	proposalCode, err := proposalContract.Code()
	if err != nil {
		t.Fatal(err)
	}

	if senderAddress, err = senderContract.Address(); err != nil {
		t.Fatal(err)
	}
	if receiverAddress, err = receiverContract.Address(); err != nil {
		t.Fatal(err)
	}
	if accessControllerAddress, err = accessControllerContract.Address(); err != nil {
		t.Fatal(err)
	}
	if bridgeVoteControllerAddress, err = bridgeVoteControllerContract.Address(); err != nil {
		t.Fatal(err)
	}
	if bridgeAddress, err = bridgeContract.Address(); err != nil {
		t.Fatal(err)
	}
	if relayerAddress, err = relayerContract.Address(); err != nil {
		t.Fatal(err)
	}
	if tip3HandlerAddress, err = tip3HandlerContract.Address(); err != nil {
		t.Fatal(err)
	}

	rootTokenContractInitVars := &RootTokenContractInitVars{
		RandomNonce:      "0",
		Name:             hex.EncodeToString([]byte("DOTON")),
		Symbol:           hex.EncodeToString([]byte("DTN")),
		Decimals:         "0x0",
		Walletcode:       walletCode.Code,
		Rootpublickey:    "0x0",
		Rootowneraddress: tip3HandlerAddress,
	}

	if rootTokenContractAddress, err = rootTokenContract.Address(rootTokenContractInitVars); err != nil {
		t.Fatal(err)
	}

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	if tonTokenWalletAddress, err = tonTokenWalletContract.Address(tonTokenWalletInitVars); err != nil {
		t.Fatal(err)
	}
	if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("senderAddress: %s \n", senderAddress)
	fmt.Printf("receiverAddress: %s \n", receiverAddress)
	fmt.Printf("accessControllerAddress: %s \n", accessControllerAddress)
	fmt.Printf("bridgeVoteControllerAddress: %s \n", bridgeVoteControllerAddress)
	fmt.Printf("bridgeAddress: %s \n", bridgeAddress)
	fmt.Printf("relayerAddress: %s \n", relayerAddress)
	fmt.Printf("rootTokenContractAddress: %s \n", rootTokenContractAddress)
	fmt.Printf("tip3HandlerAddress: %s \n", tip3HandlerAddress)
	fmt.Printf("messageHandlerAddress: %s \n", messageHandlerAddress)
	fmt.Printf("tonTokenWalletAddress: %s \n", tonTokenWalletAddress)

	if _, err = giver.SendGrams(accessControllerAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(bridgeVoteControllerAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(bridgeAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(relayerAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(senderAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(receiverAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(tip3HandlerAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(rootTokenContractAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(tonTokenWalletAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = giver.SendGrams(messageHandlerAddress, big.NewInt(500000000000), messageCallback); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	var accessController *AccessControllerContract
	if accessController, err = accessControllerContract.Deploy(
		&AccessControllerDeployParams{
			AccessCardInitState: AccessCardTvc,
			InitialValue:        "2000000000",
		}, messageCallback,
	); err != nil {
		t.Fatal(err)
	}
	if _, err = bridgeVoteControllerContract.Deploy(
		&BridgeVoteControllerDeployParams{
			ProposalCode:         proposalCode.Code,
			DeployInitialValue:   "2000000000",
			PublicKey:            "0x" + signer.Keys.Public,
			ProposalPublicKey:    "0x" + signer.Keys.Public,
			ProposalVotersAmount: "0x1",
			BridgeAddress:        bridgeAddress,
		}, messageCallback,
	); err != nil {
		t.Fatal(err)
	}
	var bridge *BridgeContract
	if bridge, err = bridgeContract.Deploy(
		&BridgeDeployParams{
			RelayerInitState:        RelayerTvc,
			AccessControllerAddress: accessControllerAddress,
			VoteControllerAddress:   bridgeVoteControllerAddress,
		},
		messageCallback,
	); err != nil {
		t.Fatal(err)
	}
	var relayer *RelayerContract
	if relayer, err = relayerContract.Deploy(
		&RelayerDeployParams{
			AccessControllerAddress: accessControllerAddress,
			MyPublicKey:             "0x" + signer.Keys.Public,
			MyInitState:             RelayerTvc,
			BridgeAddress:           bridgeAddress,
		}, messageCallback,
	); err != nil {
		t.Fatal(err)
	}
	if _, err := receiverContract.Deploy(messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = senderContract.Deploy(messageCallback); err != nil {
		t.Fatal(err)
	}
	if _, err = messageHandlerContract.Deploy(
		&MessageHandlerDeployParams{
			ProposalCode:                proposalCode.Code,
			BridgeVoteControllerAddress: bridgeVoteControllerAddress,
			BridgeVoteControllerPubKey:  "0x" + signer.Keys.Public,
		}, messageCallback,
	); err != nil {
		t.Fatal(err)
	}
	if _, err = tip3HandlerContract.Deploy(
		&Tip3HandlerDeployParams{
			ProposalCode:                proposalCode.Code,
			BridgeVoteControllerAddress: bridgeVoteControllerAddress,
			BridgeVoteControllerPubKey:  "0x" + signer.Keys.Public,
			Tip3RootAddress:             rootTokenContractAddress,
		},
		messageCallback,
	); err != nil {
		t.Fatal(err)
	}

	var rootToken *RootTokenContractContract
	if rootToken, err = rootTokenContract.Deploy(rootTokenContractInitVars, messageCallback); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 3)

	var wallet *TONTokenWalletContract
	if wallet, err = tonTokenWalletContract.Deploy(tonTokenWalletInitVars, messageCallback); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	totalSupply, err := rootToken.GetTotalSupply().Call()
	fmt.Printf("\ntotalSupply: %s\n", totalSupply)
	if err != nil {
		t.Fatal(err)
	}

	totalGranted, err := rootToken.GetTotalGranted().Call()
	fmt.Printf("\ntotalGranted: %s\n", totalGranted)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := wallet.GetBalance().Call()
	fmt.Printf("\n balance: %s\n\n", balance)
	if err != nil {
		t.Fatal(err)
	}

	// // Set Simple Message Handler

	_, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(SimpleMessageResourceID[:]), messageHandlerAddress).Send(messageCallback)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	// // Resolve message

	m := msg.Message{
		Source:       msg.ChainId(1),
		Destination:  msg.ChainId(1),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(1),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			relayerAddress,
			types.Text("hello ton"),
		},
	}

	if chain.writer.ResolveMessage(m) {
		t.Log("The message resolved")
	} else {
		t.Fatal("The message doesn't resolve")
	}

	// // Set TIP-3 Transfer Handler

	fmt.Printf("\n\nResourceId: %s\n\n", "0x"+hex.EncodeToString(Tip3ResourceID[:]))

	_, err = accessController.GrantSuperAdminRole(relayerAddress).Send(messageCallback)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	_, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(Tip3ResourceID[:]), tip3HandlerAddress).Send(messageCallback)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	m = msg.NewFungibleTransfer(
		msg.ChainId(1), msg.ChainId(1), msg.Nonce(1), big.NewInt(1000),
		Tip3ResourceID, []byte("0:5ce502787cc360dabd569bf9e26c66c88a4264149501b52419564d05d83b4aff"),
	)

	if chain.writer.ResolveMessage(m) {
		t.Log("The message resolved")
	} else {
		t.Fatal("The message doesn't resolve")
	}

	time.Sleep(time.Second * 5)

	role, err := relayer.GetRole().Call()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("\n\n Role:: %s \n\n", role)

	result, err := bridge.GetHandlerAddressByMessageType("0x" + hex.EncodeToString(Tip3ResourceID[:])).Call()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\n\nGetHandlerAddressByMessageType: %#v\n\n", result)

	balance, err = wallet.GetBalance().Call()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\n\n balance:: %s \n\n", balance)

	chain.conn.Client().Close()
}
