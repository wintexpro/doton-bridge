// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	utils "github.com/ChainSafe/ChainBridge/shared/ton"
	. "github.com/ChainSafe/ChainBridge/tonbindings"
	"github.com/ChainSafe/log15"
	log "github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/keystore"
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
		Endpoint:       "http://localhost",
		From:           "0:2089148264fb4b40dbb9ed7ba7a862403a715abf50a5730637da33d4b6453dd2",
		Insecure:       false,
		KeystorePath:   dir + "/../../keys",
		BlockstorePath: "",
		FreshStart:     true,
		Opts: map[string]string{
			"contractsPath":       dir + "/mocks/contracts",
			"receiver":            "0:1ba93200aa73341512bb7d406ccc3bae38b79628ef3fdcccd9ce2e0a133b1387",
			"burnedTokensHandler": "0:f5150024d6737c23f8a5057b391e9b39a93bc48fa047108d7613b53b6401141f",
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
		// t.Logf("\n\nEventType: %v", event.Type)
		// t.Logf("\n\nShardBlockID: %v", event.ShardBlockID)
		// t.Logf("\n\nMessageId: %v", event.MessageID)
		// t.Logf("\n\nevent: %#v", event)
	}

	workchainID := null.Int32From(0)

	var tip3HandlerAddress, relayerAddress, bridgeAddress, rootTokenContractAddress, tonTokenWalletAddress, burnedTokensHandlerAddress string

	// var (
	// 	accessControllerAddress, senderAddress, tip3HandlerAddress, tonTokenWalletAddress,
	// 	bridgeAddress, relayerAddress, receiverAddress, bridgeVoteControllerAddress, messageHandlerAddress,
	// 	burnedTokensHandlerAddress, rootTokenContractAddress string
	// )

	// // giver, err := NewGiver(chain.conn.Client(), signer, workchainID)
	// // if err != nil {
	// // 	t.Fatal(err)
	// // }

	ctx := ContractContext{
		Conn:        chain.conn.Client(),
		Signer:      &signer,
		WorkchainID: workchainID,
	}

	// // proposalContract := Proposal{Ctx: ctx}
	// senderContract := Sender{Ctx: ctx}
	// receiverContract := Receiver{Ctx: ctx}
	// accessControllerContract := AccessController{Ctx: ctx}
	// bridgeVoteControllerContract := BridgeVoteController{Ctx: ctx}
	bridgeContract := Bridge{Ctx: ctx}
	relayerContract := Relayer{Ctx: ctx}
	tip3HandlerContract := Tip3Handler{Ctx: ctx}
	rootTokenContract := RootTokenContract{Ctx: ctx}
	tonTokenWalletContract := TONTokenWallet{Ctx: ctx}
	// messageHandlerContract := MessageHandler{Ctx: ctx}
	burnedTokensHandlerContract := BurnedTokensHandler{Ctx: ctx}

	walletCode, err := tonTokenWalletContract.Code()
	if err != nil {
		t.Fatal(err)
	}

	// // proposalCode, err := proposalContract.Code()
	// // if err != nil {
	// // 	t.Fatal(err)
	// // }

	if burnedTokensHandlerAddress, err = burnedTokensHandlerContract.Address(); err != nil {
		t.Fatal(err)
	}
	// if senderAddress, err = senderContract.Address(); err != nil {
	// 	t.Fatal(err)
	// }
	// if receiverAddress, err = receiverContract.Address(); err != nil {
	// 	t.Fatal(err)
	// }
	// if accessControllerAddress, err = accessControllerContract.Address(); err != nil {
	// 	t.Fatal(err)
	// }
	// if bridgeVoteControllerAddress, err = bridgeVoteControllerContract.Address(); err != nil {
	// 	t.Fatal(err)
	// }
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
		Walletpublickey: "0x" + signer.Keys.Public,
		Owneraddress:    "0:0000000000000000000000000000000000000000000000000000000000000000",
	}

	if tonTokenWalletAddress, err = tonTokenWalletContract.Address(tonTokenWalletInitVars); err != nil {
		t.Fatal(err)
	}
	// if messageHandlerAddress, err = messageHandlerContract.Address(); err != nil {
	// 	t.Fatal(err)
	// }

	// fmt.Printf("sender: %s \n", senderAddress)
	// fmt.Printf("receiver: %s \n", receiverAddress)
	// fmt.Printf("accessController: %s \n", accessControllerAddress)
	// fmt.Printf("bridgeVoteController: %s \n", bridgeVoteControllerAddress)
	fmt.Printf("bridge: %s \n", bridgeAddress)
	fmt.Printf("relayer: %s \n", relayerAddress)
	fmt.Printf("rootTokenContract: %s \n", rootTokenContractAddress)
	fmt.Printf("tip3Handler: %s \n", tip3HandlerAddress)
	// fmt.Printf("messageHandler: %s \n", messageHandlerAddress)
	fmt.Printf("tonTokenWallet: %s \n", tonTokenWalletAddress)
	fmt.Printf("burnedTokensHandler: %s \n", burnedTokensHandlerAddress)

	// rootToken, err := rootTokenContract.New(rootTokenContractAddress, rootTokenContractInitVars)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ======== DecodeMessageBody Deposit

	// message, err := burnedTokensHandlerContract.DecodeMessageBody("te6ccgEBAQEAXwAAuloTCToBAAAAAAAAAAAAAAAAAAAAx26+SgK7w0eG2GCzVfWlzgAAAAAAAAAAAgAAAAAAAAAAAAAA6NSlEAC8VTHoeVnYNlUFd/t+bfnAVGaG+cEcOf4TVUkO2/hhcw==", true)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("message: %s \n", message.Value)

	// ======== New BurnedTokensHandlerContract

	// burnedTokensHandler, err := burnedTokensHandlerContract.New(burnedTokensHandlerAddress)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// bridgeVoteController, err := bridgeVoteControllerContract.New(bridgeVoteControllerAddress)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wallet, err := tonTokenWalletContract.New(tonTokenWalletAddress, tonTokenWalletInitVars)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// var wallet *TONTokenWalletContract
	// if wallet, err = tonTokenWalletContract.Deploy(tonTokenWalletInitVars, messageCallback); err != nil {
	// 	t.Fatal(err)
	// }

	// details, err := rootToken.GetDetails().Call()
	// fmt.Printf("\ndetails: %#v\n", details)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wdetails, err := wallet.GetDetails().Call()
	// fmt.Printf("\nwdetails: %#v\n", wdetails)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // // Set Simple Message Handler

	// _, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(SimpleMessageResourceID[:]), messageHandlerAddress).Send(messageCallback)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// time.Sleep(time.Second * 5)

	// // // Resolve message

	// m := msg.Message{
	// 	Source:       msg.ChainId(1),
	// 	Destination:  msg.ChainId(1),
	// 	Type:         SimpleMessageTransfer,
	// 	DepositNonce: msg.Nonce(1),
	// 	ResourceId:   msg.ResourceId(SimpleMessageResourceID),
	// 	Payload: []interface{}{
	// 		relayerAddress,
	// 		types.Text("hello ton"),
	// 	},
	// }

	// if chain.writer.ResolveMessage(m) {
	// 	t.Log("The message resolved")
	// } else {
	// 	t.Fatal("The message doesn't resolve")
	// }

	// // // Set TIP-3 Transfer Handler

	// fmt.Printf("\n\nResourceId: %s\n\n", "0x"+hex.EncodeToString(Tip3ResourceID[:]))

	// _, err = accessController.GrantSuperAdminRole(relayerAddress).Send(messageCallback)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// time.Sleep(time.Second * 5)

	// _, err = relayer.BridgeSetHandler("0x"+hex.EncodeToString(Tip3ResourceID[:]), tip3HandlerAddress).Send(messageCallback)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// time.Sleep(time.Second * 5)

	// m = msg.NewFungibleTransfer(
	// 	msg.ChainId(1), msg.ChainId(1), msg.Nonce(1), big.NewInt(1000),
	// 	Tip3ResourceID, []byte("0:3afffeb3a1beb9c13099552310a7e35958af24119c2ecd2923348204dbd5b624"),
	// )

	// if chain.writer.ResolveMessage(m) {
	// 	t.Log("The message resolved")
	// } else {
	// 	t.Fatal("The message doesn't resolve")
	// }

	// time.Sleep(time.Second * 5)

	// role, err := relayer.GetRole().Call()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// fmt.Printf("\n\n Role:: %s \n\n", role)

	// result, err := bridge.GetHandlerAddressByMessageType("0x" + hex.EncodeToString(Tip3ResourceID[:])).Call()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("\n\nGetHandlerAddressByMessageType: %#v\n\n", result)

	// details, err = rootToken.GetDetails().Call()
	// fmt.Printf("\ndetails: %#v\n", details)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// details, err := wallet.GetDetails().Call()
	// fmt.Printf("\nwdetails: %#v\n", details)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ========================

	// messageType := "0x" + hex.EncodeToString(FungibleTransfer3ResourceID[:])

	// input, err := json.Marshal(map[string]interface{}{
	// 	"destinationChainID": "1",
	// 	"resourceID":         messageType,
	// 	"depositNonce":       "1",
	// 	"amount":             "1000000000000",
	// 	"recipient":          "0xbc5531e87959d836550577fb7e6df9c0546686f9c11c39fe1355490edbf86173",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// paramsOfEncodeMessageBody := client.ParamsOfEncodeMessageBody{
	// 	Abi:        burnedTokensHandler.Abi,
	// 	Signer:     *burnedTokensHandler.Ctx.Signer,
	// 	IsInternal: true,
	// 	CallSet: client.CallSet{
	// 		FunctionName: "deposit",
	// 		Input:        input,
	// 	},
	// }

	// resultOfEncodeMessageBody, err := chain.conn.Client().AbiEncodeMessageBody(&paramsOfEncodeMessageBody)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = wallet.BurnByOwner("1000000000000", "100000000", burnedTokensHandlerAddress, resultOfEncodeMessageBody.Body).Send(messageCallback)

	// ========================

	// time.Sleep(time.Second * 5)

	// details, err = rootToken.GetDetails().Call()
	// fmt.Printf("\ndetails: %#v\n", details)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wdetails, err := wallet.GetDetails().Call()
	// fmt.Printf("\nwdetails: %#v\n", wdetails)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	chain.conn.Client().Close()
}
