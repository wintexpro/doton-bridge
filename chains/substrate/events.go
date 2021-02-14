// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"encoding/hex"
	"fmt"
	"math/big"

	utils "github.com/ChainSafe/ChainBridge/shared/substrate"
	events "github.com/ChainSafe/chainbridge-substrate-events"
	"github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (msg.Message, error)

const FungibleTransfer eventName = "FungibleTransfer"
const NonFungibleTransfer eventName = "NonFungibleTransfer"
const GenericTransfer eventName = "GenericTransfer"
const MessageReceived eventName = "MessageReceived"

type EventMessageReceived struct {
	Phase       types.Phase
	From        types.AccountID
	Message     types.Text
	Destination types.U8
	Nonce       types.U64
	Topics      []types.Hash
}

var SimpleMessageResourceID = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 105, 109, 112, 108, 101, 77, 101, 115, 115, 97, 103, 101, 82, 101, 115, 111, 117, 114, 99, 101,
}

var Subscriptions = []struct {
	name    eventName
	handler eventHandler
}{
	{FungibleTransfer, fungibleTransferHandler},
	{NonFungibleTransfer, nonFungibleTransferHandler},
	{GenericTransfer, genericTransferHandler},
	{MessageReceived, simpleMessageHandler},
}

func fungibleTransferHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(events.EventFungibleTransfer)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventFungibleTransfer type")
	}

	resourceId := msg.ResourceId(evt.ResourceId)
	log.Info("Got fungible transfer event!", "destination", evt.Destination, "resourceId", resourceId.Hex(), "amount", evt.Amount)

	return msg.NewFungibleTransfer(
		0, // Unset
		msg.ChainId(evt.Destination),
		msg.Nonce(evt.DepositNonce),
		evt.Amount.Int,
		resourceId,
		evt.Recipient,
	), nil
}

func nonFungibleTransferHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(events.EventNonFungibleTransfer)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventNonFungibleTransfer type")
	}

	log.Info("Got non-fungible transfer event!", "destination", evt.Destination, "resourceId", hex.EncodeToString(evt.ResourceId[:]))

	return msg.NewNonFungibleTransfer(
		0, // Unset
		msg.ChainId(evt.Destination),
		msg.Nonce(evt.DepositNonce),
		msg.ResourceId(evt.ResourceId),
		big.NewInt(0).SetBytes(evt.TokenId[:]),
		evt.Recipient,
		evt.Metadata,
	), nil
}

func genericTransferHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(events.EventGenericTransfer)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventGenericTransfer type")
	}

	log.Info("Got generic transfer event!", "destination", evt.Destination, "resourceId", hex.EncodeToString(evt.ResourceId[:]))

	return msg.NewGenericTransfer(
		0, // Unset
		msg.ChainId(evt.Destination),
		msg.Nonce(evt.DepositNonce),
		msg.ResourceId(evt.ResourceId),
		evt.Metadata,
	), nil
}

func simpleMessageHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(utils.EventMessageReceived)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventMessageReceived type")
	}

	log.Info("Got simple message recived event!", "destination", evt.Destination, "resourceId", hex.EncodeToString(SimpleMessageResourceID[:]))

	m := msg.Message{
		Source:       msg.ChainId(0),
		Destination:  msg.ChainId(evt.Destination),
		Type:         SimpleMessageTransfer,
		DepositNonce: msg.Nonce(evt.Nonce),
		ResourceId:   msg.ResourceId(SimpleMessageResourceID),
		Payload: []interface{}{
			evt.From,
			evt.Message,
		},
	}

	return m, nil
}
