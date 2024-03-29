// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	utils "github.com/ChainSafe/ChainBridge/shared/ton"
	. "github.com/ChainSafe/ChainBridge/tonbindings"
	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/keystore"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var TestLogger = log15.New("chain", "test")

var Tip3ResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 84, 105, 112, 51, 82, 101, 115, 111, 117, 114, 99, 101,
}

var FungibleTransfer3ResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 110, 190, 74, 2, 187, 195, 71, 134, 216, 96, 179, 85, 245, 165, 206, 0,
}

func TestTonChain(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &core.ChainConfig{
		Id:             2,
		Name:           "alice",
		Endpoint:       "http://192.168.0.109",
		From:           "0:b8d83bb3d617ba74e0ea44542a510ef43b5709eee93e2c8f2e6254bc5e59237f",
		Insecure:       false,
		KeystorePath:   dir + "/../../keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts: map[string]string{
			"contractsPath":       dir + "/mocks/contracts",
			"receiver":            "0:f7fdc0170f9c7e0184962aea78b1f208fe857681537854104684d62479e76e5d",
			"burnedTokensHandler": "0:dd510027840f11ce3b7b5ef0d177ccdad55f7f0fb104d8591c8c6f69babc9cc8",
			"epochVoteController": "0:35895ba7f51c612cda5f9ae7ab96a51e29eb6ad2de4c948568b8c04409912f57",
			"startBlock":          "3",
			"workchainID":         "0",
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
		EnumTypeValue: client.KeysSigner{
			Keys: client.KeyPair{
				Public: chain.kp.PublicKey(),
				Secret: chain.kp.SecretKey(),
			},
		},
	}

	TestLogger.SetHandler(log15.LvlFilterHandler(log15.LvlError, TestLogger.GetHandler()))

	router := core.NewRouter(TestLogger)
	chain.SetRouter(router)

	err = chain.Start()
	if err != nil {
		t.Fatal(err)
	}

	// messageCallback := func(event *client.ProcessingEvent) {
	// 	t.Logf("\n\nevent: %#v", event)
	// }

	workchainID := null.Int32From(0)

	var (
		tip3HandlerAddress, relayerAddress,
		bridgeAddress, rootTokenContractAddress,
		tonTokenWalletAddress, burnedTokensHandlerAddress,
		senderAddress, receiverAddress, accessControllerAddress, epochVoteControllerAddress,
		messageHandlerAddress string
	)

	ctx := ContractContext{
		Conn:        chain.conn.Client(),
		Signer:      &signer,
		WorkchainID: workchainID,
	}

	// proposalContract := Proposal{Ctx: ctx}
	// epochContract := Epoch{Ctx: ctx}
	senderContract := Sender{Ctx: ctx}
	receiverContract := Receiver{Ctx: ctx}
	accessControllerContract := AccessController{Ctx: ctx}
	epochVoteControllerContract := EpochController{Ctx: ctx}
	bridgeContract := Bridge{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}
	tip3HandlerContract := Tip3Handler{Ctx: ctx}
	rootTokenContract := RootTokenContract{Ctx: ctx}
	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}
	messageHandlerContract := MessageHandler{Ctx: ctx}
	burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		t.Fatal(err)
	}

	if burnedTokensHandlerAddress, err = burnedTokensHandlerContract.Address(); err != nil {
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
	if epochVoteControllerAddress, err = epochVoteControllerContract.Address(); err != nil {
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
		RandomNonce: "0",
		Name:        hex.EncodeToString([]byte("DOTON")),
		Symbol:      hex.EncodeToString([]byte("DTN")),
		Decimals:    "0x0",
		Walletcode:  walletCode.Code,
	}

	if rootTokenContractAddress, err = rootTokenContract.Address(rootTokenContractInitVars); err != nil {
		t.Fatal(err)
	}

	tonTokenWalletInitVars := &TONTokenWalletInitVars{
		Rootaddress:     rootTokenContractAddress,
		Code:            walletCode.Code,
		Walletpublickey: "0x" + signer.EnumTypeValue.(client.KeysSigner).Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	if tonTokenWalletAddress, err = tonTokenWalletContract.Address(tonTokenWalletInitVars); err != nil {
		t.Fatal(err)
	}
	if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
		t.Fatal(err)
	}

	fmt.Print("\n")
	fmt.Printf("%s :sender \n", senderAddress)
	fmt.Printf("%s :receiver \n", receiverAddress)
	fmt.Printf("%s :accessController \n", accessControllerAddress)
	fmt.Printf("%s :epochVoteController \n", epochVoteControllerAddress)
	fmt.Printf("%s :bridge \n", bridgeAddress)
	fmt.Printf("%s :relayer \n", relayerAddress)
	fmt.Printf("%s :rootTokenContract \n", rootTokenContractAddress)
	fmt.Printf("%s :tip3Handler \n", tip3HandlerAddress)
	fmt.Printf("%s :messageHandler \n", messageHandlerAddress)
	fmt.Printf("%s :tonTokenWallet \n", tonTokenWalletAddress)
	fmt.Printf("%s :burnedTokensHandler \n", burnedTokensHandlerAddress)
	fmt.Print("\n")

	epochVoteController, err := epochVoteControllerContract.New(epochVoteControllerAddress)
	if err != nil {
		t.Fatal(err)
	}

	currentEpochNumberMap, err := epochVoteController.CurrentEpochNumber().Call()
	if err != nil {
		t.Fatal(err)
	}
	currentEpochNumber := currentEpochNumberMap.(map[string]interface{})["currentEpochNumber"].(string)

	epochAddressMap, err := epochVoteController.GetEpochAddress(currentEpochNumber).Call()
	if err != nil {
		t.Fatal(err)
	}
	epochAddress := epochAddressMap.(map[string]interface{})["epoch"].(string)

	t.Logf("epochAddress: %s", epochAddress)

	time.Sleep(time.Second * 12)

	m := msg.Message{
		Source:       msg.ChainId(1),
		Destination:  msg.ChainId(2),
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

	chain.conn.Client().Close()
}
