// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"

	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/wintexpro/chainbridge-utils/msg"
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (msg.Message, error)

const GenericTransfer eventName = "DataReceived"

var genericResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 244, 75, 230, 77, 45, 232, 149, 69, 76, 52, 103, 2, 25, 40, 229, 94, 1,
}

var receiverABI = LoadAbi("Receiver")

var Subscriptions = []struct {
	name        eventName
	handler     eventHandler
	abi         client.Abi
	contractKey string
}{
	{GenericTransfer, genericTransferHandler, receiverABI, "receiver"},
}

func genericTransferHandler(body interface{}, log log15.Logger) (msg.Message, error) {
	data := body.(*client.DecodedMessageBody).Value.(map[string]interface{})["data"]

	dataAsBytes, err := hex.DecodeString((data.(string))[2:])
	if err != nil {
		panic(err)
	}

	log.Info("Got generic transfer event!", "destination", msg.ChainId(1), "resourceId", hex.EncodeToString(genericResourceID[:]))

	return msg.NewGenericTransfer(
		msg.ChainId(2), // Unset
		msg.ChainId(1), // TODO: get from message body
		msg.Nonce(1),   // TODO: get from message body
		msg.ResourceId(genericResourceID),
		dataAsBytes,
	), nil
}
