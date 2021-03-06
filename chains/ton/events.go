// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

type eventName string
type eventHandler func(Message, json.RawMessage, log15.Logger) (*msg.Message, error)

type SimpleMessagePayload struct {
	From    string     `json:"from"`
	Message types.Text `json:"message"`
}

type SimpleMessage struct {
	Data               string `json:"data"`
	DestinationChainID string `json:"destinationChainId"`
	Nonce              string `json:"nonce"`
}

type FungibleMessage struct {
	ResourceID         string `json:"resourceID"`
	DestinationChainID string `json:"destinationChainID"`
	Nonce              string `json:"depositNonce"`
	Amount             string `json:"amount"`
	Recipient          string `json:"recipient"`
}

const DepositEventName eventName = "Deposit"
const DepositTransfer msg.TransferType = "DepositTransfer"

const SimpleMessageEventName eventName = "DataReceived"
const SimpleMessageTransfer msg.TransferType = "SimpleMessageTransfer"

var SimpleMessageResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 105, 109, 112, 108, 101, 77, 101, 115, 115, 97, 103, 101, 82, 101, 115, 111, 117, 114, 99, 101,
}

var Subscriptions = []struct {
	name        eventName
	handler     eventHandler
	abiName     string
	contractKey string
}{
	{SimpleMessageEventName, SimpleMessageTransferHandler, "Receiver", "receiver"},
	{DepositEventName, DepositTransferHandler, "BurnedTokensHandler", "burnedTokensHandler"},
}

func DepositTransferHandler(message Message, body json.RawMessage, log log15.Logger) (*msg.Message, error) {
	var resourceId [32]byte
	fungibleMessage := FungibleMessage{}

	err := json.Unmarshal(body, &fungibleMessage)
	if err != nil {
		return nil, err
	}

	chainIDInt, err := strconv.Atoi(fungibleMessage.DestinationChainID)
	if err != nil {
		return nil, err
	}

	nonceInt, err := strconv.Atoi(fungibleMessage.Nonce)
	if err != nil {
		return nil, err
	}

	amountInt, err := strconv.Atoi(fungibleMessage.Amount)
	if err != nil {
		return nil, err
	}

	amount := big.NewInt(int64(amountInt))
	resourceIDBytes, err := hex.DecodeString(fungibleMessage.ResourceID[2:])
	if err != nil {
		return nil, err
	}

	recipient, err := hex.DecodeString(fungibleMessage.Recipient[2:])
	if err != nil {
		return nil, err
	}

	for i, b := range resourceIDBytes {
		resourceId[i] = b
	}

	newFungibleMessage := msg.NewFungibleTransfer(
		0, // Unset
		msg.ChainId(int64(chainIDInt)),
		msg.Nonce(int64(nonceInt)),
		amount,
		resourceId,
		recipient,
	)

	return &newFungibleMessage, nil
}

func SimpleMessageTransferHandler(message Message, body json.RawMessage, log log15.Logger) (*msg.Message, error) {
	simpleMesage := SimpleMessage{}

	err := json.Unmarshal(body, &simpleMesage)
	if err != nil {
		return nil, err
	}

	chainIDAsBytes, err := hex.DecodeString(simpleMesage.DestinationChainID[2:])
	if err != nil {
		return nil, err
	}

	chainID := big.NewInt(0).SetBytes(chainIDAsBytes).Uint64()

	nonceAsBytes, err := hex.DecodeString(simpleMesage.Nonce[2:])
	if err != nil {
		return nil, err
	}

	nonce := big.NewInt(0).SetBytes(nonceAsBytes).Uint64()

	data, err := hex.DecodeString(simpleMesage.Data[2:])
	if err != nil {
		return nil, err
	}

	log.Info("Got simple message transfer event!", "destination", chainID, "nonce", nonce, "resourceId", hex.EncodeToString(SimpleMessageResourceID[:]))

	m := msg.Message{
		Source:       msg.ChainId(0),
		Destination:  msg.ChainId(chainID),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(nonce),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			message.Src,
			types.Text(data),
		},
	}

	return &m, nil
}
