// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ChainSafe/ChainBridge/config"
	"github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"

	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	. "github.com/ChainSafe/ChainBridge/tonbindings"
	log "github.com/ChainSafe/log15"
	"github.com/urfave/cli/v2"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	"github.com/wintexpro/chainbridge-utils/keystore"
)

var ZERO_ADDRESS = "0:0000000000000000000000000000000000000000000000000000000000000000"

// handleDeployCmd deploy the set of ton contracts

func handleSendGrams(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Send grams...")
	return execute(ctx, dHandler, sendGrams)
}

func handleDeployCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Starting deploy...")
	return execute(ctx, dHandler, deploy)
}

func handleSetupCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Setup...")
	return execute(ctx, dHandler, setup)
}

func handleDeployWalletCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Starting deploy Wallet...")
	return execute(ctx, dHandler, deployWallet)
}

func handleGetBalanceCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Getting balance...")
	return execute(ctx, dHandler, getBalance)
}

func execute(ctx *cli.Context, dHandler *dataHandler, fn func(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error) error {
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Check for test key flag
	var ks string
	var insecure bool
	if key := ctx.String(config.TestKeyFlag.Name); key != "" {
		ks = key
		insecure = true
	} else {
		ks = cfg.KeystorePath
	}

	// Used to signal core shutdown due to fatal error
	// sysErr := make(chan error)

	for _, chain := range cfg.Chains {
		logger := log.Root().New("chain", chain.Name)

		if chain.Type == "ton" {
			kpI, err := keystore.KeypairFromAddress(chain.From, keystore.TonChain, ks, insecure)
			if err != nil {
				return err
			}
			kp, _ := kpI.(*ed25519.Keypair)
			conn := connection.NewConnection(chain.Endpoint, false, chain.Opts["workchainID"], logger)

			workchainID, err := strconv.Atoi(chain.Opts["workchainID"])
			if err != nil {
				return err
			}

			workchainIDNull := null.Int32From(int32(workchainID))

			signer := client.Signer{
				Type: client.KeysSignerType,
				Keys: client.KeyPair{
					Public: kp.PublicKey(),
					Secret: kp.SecretKey(),
				},
			}

			if err := conn.Connect(); err != nil {
				return err
			}

			return fn(conn, workchainIDNull, &signer, logger)
		}
	}

	return nil
}

func deploy(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	var (
		accessControllerAddress, senderAddress, tip3HandlerAddress,
		bridgeAddress, relayerAddress, receiverAddress, bridgeVoteControllerAddress, messageHandlerAddress,
		burnedTokensHandlerAddress, rootTokenContractAddress string
	)

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
		}
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
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
	burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		return err
	}

	var proposalCode *client.ResultOfGetCodeFromTvc
	if proposalCode, err = proposalContract.Code(); err != nil {
		return err
	}
	if burnedTokensHandlerAddress, err = burnedTokensHandlerContract.Address(); err != nil {
		return err
	}
	if senderAddress, err = senderContract.Address(); err != nil {
		return err
	}
	if receiverAddress, err = receiverContract.Address(); err != nil {
		return err
	}
	if accessControllerAddress, err = accessControllerContract.Address(); err != nil {
		return err
	}
	if bridgeVoteControllerAddress, err = bridgeVoteControllerContract.Address(); err != nil {
		return err
	}
	if bridgeAddress, err = bridgeContract.Address(); err != nil {
		return err
	}
	if relayerAddress, err = relayerContract.Address(); err != nil {
		return err
	}
	if tip3HandlerAddress, err = tip3HandlerContract.Address(); err != nil {
		return err
	}

	rootTokenContractInitVars := &RootTokenContractInitVars{
		RandomNonce: "0",
		Name:        hex.EncodeToString([]byte("DOTON")),
		Symbol:      hex.EncodeToString([]byte("DTN")),
		Decimals:    "0xc",
		Walletcode:  walletCode.Code,
	}

	if rootTokenContractAddress, err = rootTokenContract.Address(rootTokenContractInitVars); err != nil {
		return err
	}

	if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
		return err
	}

	fmt.Printf("sender: %s \n", senderAddress)
	fmt.Printf("receiver: %s \n", receiverAddress)
	fmt.Printf("accessController: %s \n", accessControllerAddress)
	fmt.Printf("bridgeVoteController: %s \n", bridgeVoteControllerAddress)
	fmt.Printf("bridge: %s \n", bridgeAddress)
	fmt.Printf("relayer: %s \n", relayerAddress)
	fmt.Printf("rootTokenContract: %s \n", rootTokenContractAddress)
	fmt.Printf("tip3Handler: %s \n", tip3HandlerAddress)
	fmt.Printf("messageHandler: %s \n", messageHandlerAddress)
	fmt.Printf("burnedTokensHandler: %s \n", burnedTokensHandlerAddress)

	if _, err = accessControllerContract.Deploy(
		&AccessControllerDeployParams{
			AccessCardInitState: AccessCardTvc,
			InitialValue:        "2000000000",
		}, messageCallback("accessControllerContract.Deploy"),
	); err != nil {
		return err
	}
	if _, err = bridgeVoteControllerContract.Deploy(
		&BridgeVoteControllerDeployParams{
			ProposalCode:         proposalCode.Code,
			DeployInitialValue:   "2000000000",
			PublicKey:            "0x" + signer.Keys.Public,
			ProposalPublicKey:    "0x" + signer.Keys.Public,
			ProposalVotersAmount: "0x1",
			BridgeAddress:        bridgeAddress,
		}, messageCallback("bridgeVoteControllerContract.Deploy"),
	); err != nil {
		return err
	}
	if _, err = bridgeContract.Deploy(
		&BridgeDeployParams{
			RelayerInitState:        RelayerTvc,
			AccessControllerAddress: accessControllerAddress,
			VoteControllerAddress:   bridgeVoteControllerAddress,
		}, messageCallback("bridgeContract.Deploy"),
	); err != nil {
		return err
	}
	if _, err = relayerContract.Deploy(
		&RelayerDeployParams{
			AccessControllerAddress: accessControllerAddress,
			MyPublicKey:             "0x" + signer.Keys.Public,
			MyInitState:             RelayerTvc,
			BridgeAddress:           bridgeAddress,
		}, messageCallback("relayerContract.Deploy"),
	); err != nil {
		return err
	}
	if _, err := receiverContract.Deploy(messageCallback("receiverContract.Deploy")); err != nil {
		return err
	}
	if _, err = senderContract.Deploy(messageCallback("senderContract.Deploy")); err != nil {
		return err
	}
	if _, err = messageHandlerContract.Deploy(
		&MessageHandlerDeployParams{
			ProposalCode:                proposalCode.Code,
			BridgeVoteControllerAddress: bridgeVoteControllerAddress,
			BridgeVoteControllerPubKey:  "0x" + signer.Keys.Public,
		}, messageCallback("messageHandlerContract.Deploy"),
	); err != nil {
		return err
	}
	if _, err = tip3HandlerContract.Deploy(
		&Tip3HandlerDeployParams{
			ProposalCode:                proposalCode.Code,
			BridgeVoteControllerAddress: bridgeVoteControllerAddress,
			BridgeVoteControllerPubKey:  "0x" + signer.Keys.Public,
			Tip3RootAddress:             rootTokenContractAddress,
		}, messageCallback("tip3HandlerContract.Deploy"),
	); err != nil {
		return err
	}

	// var burnedTokensHandler *BurnedTokensHandlerContract
	if _, err = burnedTokensHandlerContract.Deploy(
		&BurnedTokensHandlerDeployParams{
			Tip3RootAddress: rootTokenContractAddress,
		}, messageCallback("burnedTokensHandlerContract.Deploy"),
	); err != nil {
		return err
	}
	rootTokenContractDeployParams := RootTokenContractDeployParams{
		Rootpublickey:    "0x0",
		Rootowneraddress: tip3HandlerAddress,
	}
	// var rootToken *RootTokenContractContract
	if _, err = rootTokenContract.Deploy(&rootTokenContractDeployParams, rootTokenContractInitVars, messageCallback("rootTokenContract.Deploy")); err != nil {
		return err
	}

	return nil
}

func setup(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	var tip3HandlerAddress, relayerAddress, messageHandlerAddress string
	var err error

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
		}
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
		WorkchainID: workchainID,
	}

	tip3HandlerContract := Tip3Handler{Ctx: ctx}
	accessControllerContract := AccessController{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}
	bridgeContract := Bridge{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}

	if relayerAddress, err = relayerContract.Address(); err != nil {
		return err
	}
	if tip3HandlerAddress, err = tip3HandlerContract.Address(); err != nil {
		return err
	}
	if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
		return err
	}

	var accessController *AccessControllerContract
	if accessController, err = accessControllerContract.New(""); err != nil {
		return err
	}
	var relayer *RelayerContract
	if relayer, err = relayerContract.New(""); err != nil {
		return err
	}

	var bridge *BridgeContract
	if bridge, err = bridgeContract.New(""); err != nil {
		return err
	}

	fmt.Printf("accessController: %s \n", accessController.Address)
	fmt.Printf("bridge: %s \n", bridge.Address)
	fmt.Printf("relayer: %s \n", relayer.Address)

	// Set Handlers
	var SimpleMessageResourceID = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 105, 109, 112, 108, 101, 77, 101, 115, 115, 97, 103, 101, 82, 101, 115, 111, 117, 114, 99, 101,
	}

	if _, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(SimpleMessageResourceID[:]), messageHandlerAddress).Send(messageCallback("BridgeSetHandler(SimpleMessageResourceID)")); err != nil {
		return err
	}

	if _, err = relayer.BridgeSetHandler("0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00", tip3HandlerAddress).Send(messageCallback("BridgeSetHandler(Tip3ResourceID)")); err != nil {
		return err
	}

	_, err = accessController.GrantSuperAdminRole(relayerAddress).Send(messageCallback("GrantSuperAdminRole"))
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 30)

	result, err := bridge.GetHandlerAddressByMessageType("0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00").Call()
	if err != nil {
		return err
	}

	fmt.Printf("\n\n %s Handler address: %#v\n\n", "0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00", result)

	result, err = bridge.GetHandlerAddressByMessageType("0x" + hex.EncodeToString(SimpleMessageResourceID[:])).Call()
	if err != nil {
		return err
	}

	fmt.Printf("\n\n %s Handler address: %#v\n\n", "0x"+hex.EncodeToString(SimpleMessageResourceID[:]), result)

	return nil
}

func sendGrams(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	var (
		accessControllerAddress, senderAddress, tip3HandlerAddress, tonTokenWalletAddress,
		bridgeAddress, relayerAddress, receiverAddress, bridgeVoteControllerAddress, messageHandlerAddress,
		burnedTokensHandlerAddress, rootTokenContractAddress string
	)

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
		}
	}

	giver, err := NewGiver(conn.Client(), *signer, workchainID)
	if err != nil {
		return err
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
		WorkchainID: workchainID,
	}

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
	burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}

	if burnedTokensHandlerAddress, err = burnedTokensHandlerContract.Address(); err != nil {
		return err
	}
	if senderAddress, err = senderContract.Address(); err != nil {
		return err
	}
	if receiverAddress, err = receiverContract.Address(); err != nil {
		return err
	}
	if accessControllerAddress, err = accessControllerContract.Address(); err != nil {
		return err
	}
	if bridgeVoteControllerAddress, err = bridgeVoteControllerContract.Address(); err != nil {
		return err
	}
	if bridgeAddress, err = bridgeContract.Address(); err != nil {
		return err
	}
	if relayerAddress, err = relayerContract.Address(); err != nil {
		return err
	}
	if tip3HandlerAddress, err = tip3HandlerContract.Address(); err != nil {
		return err
	}

	if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
		return err
	}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		return err
	}

	rootTokenContractInitVars := &RootTokenContractInitVars{
		RandomNonce: "0",
		Name:        hex.EncodeToString([]byte("DOTON")),
		Symbol:      hex.EncodeToString([]byte("DTN")),
		Decimals:    "0xc",
		Walletcode:  walletCode.Code,
	}

	if rootTokenContractAddress, err = rootTokenContract.Address(rootTokenContractInitVars); err != nil {
		return err
	}

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	if tonTokenWalletAddress, err = tonTokenWalletContract.Address(tonTokenWalletInitVars); err != nil {
		return err
	}

	if _, err = giver.SendGrams(tonTokenWalletAddress, big.NewInt(200000000000), messageCallback("SendGrams(tonTokenWalletAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(burnedTokensHandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(burnedTokensHandlerAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(accessControllerAddress, big.NewInt(200000000000), messageCallback("SendGrams(accessControllerAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(bridgeVoteControllerAddress, big.NewInt(200000000000), messageCallback("SendGrams(bridgeVoteControllerAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(bridgeAddress, big.NewInt(200000000000), messageCallback("SendGrams(bridgeAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(relayerAddress, big.NewInt(200000000000), messageCallback("SendGrams(relayerAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(senderAddress, big.NewInt(200000000000), messageCallback("SendGrams(senderAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(receiverAddress, big.NewInt(200000000000), messageCallback("SendGrams(receiverAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(tip3HandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(tip3HandlerAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(rootTokenContractAddress, big.NewInt(200000000000), messageCallback("SendGrams(rootTokenContractAddress)")); err != nil {
		return err
	}
	if _, err = giver.SendGrams(messageHandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(messageHandlerAddress)")); err != nil {
		return err
	}

	return nil
}

type Result = map[string]interface{}
type WalletDetails = map[string]interface{}

func getBalance(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
		WorkchainID: workchainID,
	}

	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		return err
	}

	rootTokenContract := RootTokenContract{Ctx: ctx}

	rootTokenContractInitVars := &RootTokenContractInitVars{
		RandomNonce: "0",
		Name:        hex.EncodeToString([]byte("DOTON")),
		Symbol:      hex.EncodeToString([]byte("DTN")),
		Decimals:    "0xc",
		Walletcode:  walletCode.Code,
	}

	rootTokenContractAddress, err := rootTokenContract.Address(rootTokenContractInitVars)
	if err != nil {
		return err
	}

	rootToken, err := rootTokenContract.New(rootTokenContractAddress, rootTokenContractInitVars)
	if err != nil {
		return err
	}

	addressResult, err := rootToken.GetWalletAddress("0x"+signer.Keys.Public, ZERO_ADDRESS).Call()
	if err != nil {
		return err
	}
	address := addressResult.(Result)["value0"].(string)

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	wallet, err := tonTokenWalletContract.New(address, tonTokenWalletInitVars)
	if err != nil {
		return err
	}

	result, err := wallet.GetDetails().Call()

	fmt.Printf(" \n Balance of %s: %s \n", address, result.(Result)["value0"].(WalletDetails)["balance"])

	return nil
}

func deployWallet(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
		}
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
		WorkchainID: workchainID,
	}

	rootTokenContract := RootTokenContract{Ctx: ctx}
	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		return err
	}

	rootTokenContractInitVars := &RootTokenContractInitVars{
		RandomNonce: "0",
		Name:        hex.EncodeToString([]byte("DOTON")),
		Symbol:      hex.EncodeToString([]byte("DTN")),
		Decimals:    "0xc",
		Walletcode:  walletCode.Code,
	}

	rootTokenContractAddress, err := rootTokenContract.Address(rootTokenContractInitVars)
	if err != nil {
		return err
	}

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.Keys.Public,
		Owneraddress:    ZERO_ADDRESS,
	}

	rootToken, err := rootTokenContract.New(rootTokenContractAddress, rootTokenContractInitVars)
	if err != nil {
		return err
	}

	addressResult, err := rootToken.GetWalletAddress("0x"+signer.Keys.Public, ZERO_ADDRESS).Call()
	if err != nil {
		return err
	}

	if _, err = tonTokenWalletContract.Deploy(tonTokenWalletInitVars, messageCallback("tonTokenWalletContract.Deploy")); err != nil {
		return err
	}

	logger.Info("Wallet address", "Address", addressResult.(Result)["value0"].(string))

	return nil
}
