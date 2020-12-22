// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/wintexpro/chainbridge-utils/msg"
)

type eventName string
type eventHandler func(map[string]interface{}, interface{}, log15.Logger) (msg.Message, error)

type SimpleMessagePayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

const SimpleMessage eventName = "DataReceived"

var SimpleMessageTransfer msg.TransferType = "SimpleMessageTransfer"

var SimpleMessageResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 105, 109, 112, 108, 101, 77, 101, 115, 115, 97, 103, 101, 82, 101, 115, 111, 117, 114, 99, 101,
}

var Subscriptions = []struct {
	name        eventName
	handler     eventHandler
	abiName     string
	contractKey string
}{
	{SimpleMessage, SimpleMessageTransferHandler, "Receiver", "receiver"},
}

func SimpleMessageTransferHandler(message map[string]interface{}, body interface{}, log log15.Logger) (msg.Message, error) {
	chainIDAsBytes, err := hex.DecodeString(
		((body.(*client.DecodedMessageBody).Value.(map[string]interface{})["destinationChainId"]).(string))[2:],
	)
	if err != nil {
		panic(err)
	}

	chainID := big.NewInt(0).SetBytes(chainIDAsBytes).Uint64()

	nonceAsBytes, err := hex.DecodeString(
		((body.(*client.DecodedMessageBody).Value.(map[string]interface{})["nonce"]).(string))[2:],
	)
	if err != nil {
		panic(err)
	}

	nonce := big.NewInt(0).SetBytes(nonceAsBytes).Uint64()

	data, err := hex.DecodeString(
		(body.(*client.DecodedMessageBody).Value.(map[string]interface{})["data"].(string))[2:],
	)
	if err != nil {
		panic(err)
	}

	payload := SimpleMessagePayload{
		From:    message["src"].(string),
		Message: string(data),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	log.Info("Got simple message transfer event!", "destination", chainID, "nonce", nonce, "resourceId", hex.EncodeToString(SimpleMessageResourceID[:]))

	m := msg.Message{
		Source:       msg.ChainId(0),
		Destination:  msg.ChainId(chainID),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(nonce),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			[]byte(payloadJSON),
		},
	}

	return m, nil
}
