// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ChainSafe/ChainBridge/config"
	"github.com/btcsuite/btcutil/base58"
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

func handleSendTokensCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Sending tokens...")
	return execute(ctx, dHandler, sendTokens(ctx.String(config.AmountFlag.Name), ctx.String(config.ToFlag.Name), ctx.String(config.NonceFlag.Name)))
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
				EnumTypeValue: client.KeysSigner{
					Keys: client.KeyPair{
						Public: kp.PublicKey(),
						Secret: kp.SecretKey(),
					},
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
		bridgeAddress, relayerAddress, receiverAddress, messageHandlerAddress,
		burnedTokensHandlerAddress, rootTokenContractAddress, epochControllerAddress string
	)

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			// logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
		}
	}

	ctx := ContractContext{
		Conn:        conn.Client(),
		Signer:      signer,
		WorkchainID: workchainID,
	}

	proposalContract := Proposal{Ctx: ctx}
	epochContract := Epoch{Ctx: ctx}
	epochControllerContract := EpochController{Ctx: ctx}
	senderContract := Sender{Ctx: ctx}
	receiverContract := Receiver{Ctx: ctx}
	accessControllerContract := AccessController{Ctx: ctx}
	bridgeContract := Bridge{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}
	tip3HandlerContract := Tip3Handler{Ctx: ctx}
	rootTokenContract := RootTokenContract{Ctx: ctx}
	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}
	burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}

	var err error

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
	if epochControllerAddress, err = epochControllerContract.Address(); err != nil {
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

	var proposalCode *client.ResultOfGetCodeFromTvc
	if proposalCode, err = proposalContract.Code(); err != nil {
		return err
	}

	var epochCode *client.ResultOfGetCodeFromTvc
	if epochCode, err = epochContract.Code(); err != nil {
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
	fmt.Printf("epochControllerAddress: %s \n", epochControllerAddress)
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
	time.Sleep(time.Second)
	// epochControllerDeployParams *EpochControllerDeployParams, messageCallback
	if _, err = epochControllerContract.Deploy(&EpochControllerDeployParams{
		EpochCode:            epochCode.Code,
		ProposalCode:         proposalCode.Code,
		DeployInitialValue:   "2000000000",
		PublicKey:            "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		ProposalVotersAmount: "0x1",
		BridgeAddress:        bridgeAddress,
		FirstEraDuration:     "0xa",
		SecondEraDuration:    "0xa",
	}, messageCallback("epochControllerContract.Deploy")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err = bridgeContract.Deploy(
		&BridgeDeployParams{
			RelayerInitState:        RelayerTvc,
			AccessControllerAddress: accessControllerAddress,
			VoteControllerAddress:   epochControllerAddress,
		}, messageCallback("bridgeContract.Deploy"),
	); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err = relayerContract.Deploy(
		&RelayerDeployParams{
			AccessControllerAddress: accessControllerAddress,
			MyPublicKey:             "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
			MyInitState:             RelayerTvc,
			BridgeAddress:           bridgeAddress,
		}, messageCallback("relayerContract.Deploy"),
	); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err := receiverContract.Deploy(messageCallback("receiverContract.Deploy")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err = senderContract.Deploy(messageCallback("senderContract.Deploy")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err = messageHandlerContract.Deploy(
		&MessageHandlerDeployParams{
			ProposalCode:          proposalCode.Code,
			EpochControllerPubKey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		}, messageCallback("messageHandlerContract.Deploy"),
	); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err = tip3HandlerContract.Deploy(
		&Tip3HandlerDeployParams{
			ProposalCode:          proposalCode.Code,
			EpochControllerPubKey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
			Tip3RootAddress:       rootTokenContractAddress,
		}, messageCallback("tip3HandlerContract.Deploy"),
	); err != nil {
		return err
	}
	time.Sleep(time.Second)
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
	time.Sleep(time.Second)

	return nil
}

func setup(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	var tip3HandlerAddress, relayerAddress, messageHandlerAddress string
	var err error

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			// logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
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

	_, err = accessController.GrantSuperAdminRole(relayerAddress).Send(messageCallback("GrantSuperAdminRole"))
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	if _, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(SimpleMessageResourceID[:]), messageHandlerAddress).Send(messageCallback("BridgeSetHandler(SimpleMessageResourceID)")); err != nil {
		return err
	}
	time.Sleep(time.Second)

	if _, err = relayer.BridgeSetHandler("0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00", tip3HandlerAddress).Send(messageCallback("BridgeSetHandler(Tip3ResourceID)")); err != nil {
		return err
	}
	time.Sleep(time.Second)

	result, err := bridge.GetHandlerAddressByMessageType("0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00").Call()
	if err != nil {
		return err
	}

	fmt.Printf("\n %s Handler address: %#v", "0x000000000000000000000000000000c76ebe4a02bbc34786d860b355f5a5ce00", result)

	result, err = bridge.GetHandlerAddressByMessageType("0x" + hex.EncodeToString(SimpleMessageResourceID[:])).Call()
	if err != nil {
		return err
	}

	fmt.Printf("\n %s Handler address: %#v", "0x"+hex.EncodeToString(SimpleMessageResourceID[:]), result)

	return nil
}

func sendGrams(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	var (
		accessControllerAddress, senderAddress, tip3HandlerAddress, tonTokenWalletAddress,
		bridgeAddress, relayerAddress, receiverAddress, messageHandlerAddress,
		burnedTokensHandlerAddress, rootTokenContractAddress, epochControllerAddress string
	)

	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			logger.Info("Message status:", "label", eventLabel, "event", event)
			// logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
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

	epochControllerContract := EpochController{Ctx: ctx}
	senderContract := Sender{Ctx: ctx}
	receiverContract := Receiver{Ctx: ctx}
	accessControllerContract := AccessController{Ctx: ctx}
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
	if epochControllerAddress, err = epochControllerContract.Address(); err != nil {
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
		Walletpublickey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	if tonTokenWalletAddress, err = tonTokenWalletContract.Address(tonTokenWalletInitVars); err != nil {
		return err
	}

	logger.Info("SendGrams(tonTokenWalletAddress)")
	if _, err = giver.SendGrams(tonTokenWalletAddress, big.NewInt(200000000000), messageCallback("SendGrams(tonTokenWalletAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(burnedTokensHandlerAddress)")
	if _, err = giver.SendGrams(burnedTokensHandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(burnedTokensHandlerAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(accessControllerAddress)")
	if _, err = giver.SendGrams(accessControllerAddress, big.NewInt(200000000000), messageCallback("SendGrams(accessControllerAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(epochControllerAddress)")
	if _, err = giver.SendGrams(epochControllerAddress, big.NewInt(200000000000), messageCallback("SendGrams(epochControllerAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(bridgeAddress)")
	if _, err = giver.SendGrams(bridgeAddress, big.NewInt(200000000000), messageCallback("SendGrams(bridgeAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(relayerAddress)")
	if _, err = giver.SendGrams(relayerAddress, big.NewInt(200000000000), messageCallback("SendGrams(relayerAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(senderAddress)")
	if _, err = giver.SendGrams(senderAddress, big.NewInt(200000000000), messageCallback("SendGrams(senderAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(receiverAddress)")
	if _, err = giver.SendGrams(receiverAddress, big.NewInt(200000000000), messageCallback("SendGrams(receiverAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(tip3HandlerAddress)")
	if _, err = giver.SendGrams(tip3HandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(tip3HandlerAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(rootTokenContractAddress)")
	if _, err = giver.SendGrams(rootTokenContractAddress, big.NewInt(200000000000), messageCallback("SendGrams(rootTokenContractAddress)")); err != nil {
		return err
	}
	time.Sleep(time.Second)
	logger.Info("SendGrams(messageHandlerAddress)")
	if _, err = giver.SendGrams(messageHandlerAddress, big.NewInt(200000000000), messageCallback("SendGrams(messageHandlerAddress)")); err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	logger.Info("Exit: 0")

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

	addressResult, err := rootToken.GetWalletAddress("0x"+signer.EnumTypeValue.(client.KeysSigner).Keys.Public, ZERO_ADDRESS).Call()
	if err != nil {
		return err
	}
	address := addressResult.(Result)["value0"].(string)

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	wallet, err := tonTokenWalletContract.New(address, tonTokenWalletInitVars)
	if err != nil {
		return err
	}

	result, err := wallet.GetDetails().Call()
	if err != nil {
		return err
	}

	amount := new(big.Int)
	amount.SetString(result.(Result)["value0"].(WalletDetails)["balance"].(string), 10)
	amount.Div(amount, big.NewInt(1000000000000))

	fmt.Printf(" \n Balance of %s: %s DTN \n", address, amount)

	return nil
}

func decodeSS58AddressToPublicKey(address string) string {
	return fmt.Sprintf("0x%x", base58.Decode(string(address))[1:33])
}

func sendTokens(amount string, to string, nonce string) func(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	return func(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
		ctx := ContractContext{
			Conn:        conn.Client(),
			Signer:      signer,
			WorkchainID: workchainID,
		}

		burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}
		burnedTokensAddress, err := burnedTokensHandlerContract.Address()
		if err != nil {
			return err
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

		addressResult, err := rootToken.GetWalletAddress("0x"+signer.EnumTypeValue.(client.KeysSigner).Keys.Public, ZERO_ADDRESS).Call()
		if err != nil {
			return err
		}
		address := addressResult.(Result)["value0"].(string)

		tonTokenWalletInitVars := &TONTokenWalletInitVars{
			Rootaddress:     rootTokenContractAddress,
			Code:            walletCode.Code,
			Walletpublickey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
			Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
		}

		wallet, err := tonTokenWalletContract.New(address, tonTokenWalletInitVars)
		if err != nil {
			return err
		}

		messageType := "0x" + hex.EncodeToString(
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 110, 190, 74, 2, 187, 195, 71, 134, 216, 96, 179, 85, 245, 165, 206, 0},
		)

		decimals := big.NewInt(12)
		mul := big.NewInt(10)
		bamount := new(big.Int)
		_, ok := bamount.SetString(amount, 10)

		if !ok {
			panic(errors.New("error"))
		}

		mul.Exp(mul, decimals, nil)

		bamount.Mul(bamount, mul)

		input, err := json.Marshal(map[string]interface{}{
			"destinationChainID": "1",
			"resourceID":         messageType,
			"depositNonce":       nonce,
			"amount":             bamount.String(),
			"recipient":          decodeSS58AddressToPublicKey(to),
		})
		if err != nil {
			return err
		}

		abi, err := burnedTokensHandlerContract.Abi()
		if err != nil {
			return err
		}

		paramsOfEncodeMessageBody := client.ParamsOfEncodeMessageBody{
			Abi:        *abi,
			Signer:     *signer,
			IsInternal: true,
			CallSet: client.CallSet{
				FunctionName: "deposit",
				Input:        input,
			},
		}

		resultOfEncodeMessageBody, err := conn.Client().AbiEncodeMessageBody(&paramsOfEncodeMessageBody)
		if err != nil {
			return err
		}

		messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
			return func(event *client.ProcessingEvent) {
				// logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
			}
		}

		wallet.BurnByOwner(bamount.String(), "100000000", burnedTokensAddress, resultOfEncodeMessageBody.Body).Send(messageCallback("wallet.BurnByOwner"))

		return nil
	}
}

func deployWallet(conn *connection.Connection, workchainID null.Int32, signer *client.Signer, logger log.Logger) error {
	messageCallback := func(eventLabel string) func(event *client.ProcessingEvent) {
		return func(event *client.ProcessingEvent) {
			// logger.Info("Message status:", "label", eventLabel, "Status", event.Type, "MessageID", event.MessageID, "ShardBlockID", event.ShardBlockID)
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
		Walletpublickey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		Owneraddress:    ZERO_ADDRESS,
	}

	rootToken, err := rootTokenContract.New(rootTokenContractAddress, rootTokenContractInitVars)
	if err != nil {
		return err
	}

	addressResult, err := rootToken.GetWalletAddress("0x"+signer.EnumTypeValue.(client.KeysSigner).Keys.Public, ZERO_ADDRESS).Call()
	if err != nil {
		return err
	}

	if _, err = tonTokenWalletContract.Deploy(tonTokenWalletInitVars, messageCallback("tonTokenWalletContract.Deploy")); err != nil {
		return err
	}

	logger.Info("Wallet address", "Address", addressResult.(Result)["value0"].(string))

	return nil
}
