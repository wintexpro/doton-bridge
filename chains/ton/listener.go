// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"time"

	metrics "github.com/ByKeks/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/ChainBridge/chains"
	"github.com/ChainSafe/ChainBridge/connections/ton"
	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

type listener struct {
	cfg         Config
	conn        *ton.Connection
	latestBlock metrics.LatestBlock
	log         log15.Logger
	router      chains.Router
	stop        <-chan int
}

type Message struct {
	ID        int64  `json:"id"`
	Status    int64  `json:"status"`
	RawBody   string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

// Frequency of polling for a new block
var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5

func NewListener(conn *ton.Connection, cfg *Config, log log15.Logger, stop <-chan int) *listener {
	return &listener{
		cfg:         *cfg,
		conn:        conn,
		log:         log,
		latestBlock: metrics.LatestBlock{LastUpdated: time.Now()},
		stop:        stop,
	}
}

func (l *listener) setRouter(r chains.Router) {
	// l.router = r
}

// start creates the initial subscription for all events
func (l *listener) start() error {
	// // Check whether latest is less than starting block
	// header, err := l.conn.api.RPC.Chain.GetHeaderLatest()
	// if err != nil {
	// 	return err
	// }
	// if uint64(header.Number) < l.startBlock {
	// 	return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, header.Number)
	// }

	// // for _, sub := range Subscriptions {
	// // 	err := l.registerEventHandler(sub.name, sub.handler)
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // }

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

func loadAbi(name string) client.Abi {
	content, err := ioutil.ReadFile("./contracts/" + name + ".abi.json")
	if err != nil {
		panic(err)
	}
	abi := client.Abi{Type: client.ContractAbiType}
	if err = json.Unmarshal(content, &abi.Value); err != nil {
		panic(err)
	}

	return abi
}

func (l *listener) getBlock(blockNumber *big.Int) (*connection.BlockType, error) {
	params := client.ParamsOfQueryCollection{
		Collection: "blocks",
		Result:     "seq_no gen_utime",
		Limit:      null.Uint32From(1),
		Filter: json.RawMessage(`{
			"workchain_id":{"eq":-1},
			"status":{"eq": 2},
			"seq_no":{"eq": ` + blockNumber.String() + `}
		}`),
		Order: []client.OrderBy{{
			Path:      "seq_no",
			Direction: client.DescSortDirection,
		}},
	}

	res, err := l.conn.Client().NetQueryCollection(&params)

	if err != nil {
		return nil, err
	}

	if len(res.Result) <= 0 {
		return nil, errors.New("No blocks found")
	}

	currentBlock := &connection.BlockType{}

	s, err := json.Marshal(res.Result[0])

	if err != nil {
		return nil, err
	}

	json.Unmarshal(s, &currentBlock)

	return currentBlock, nil
}

func (l *listener) decodeMessageBody(message *Message, abi client.Abi) (*client.DecodedMessageBody, error) {
	params := client.ParamsOfDecodeMessageBody{
		Abi:  abi,
		Body: message.RawBody,
	}

	return l.conn.Client().AbiDecodeMessageBody(&params)
}

func (l *listener) getMessage(address string, prevBlock, currentBlock *connection.BlockType) (*Message, error) {
	params := client.ParamsOfQueryCollection{
		Collection: "messages",
		Result:     "id status created_at body",
		Filter: json.RawMessage(`{
			"status": { "eq": 5 },
			"src": { "eq": "` + address + `" },
			"created_at": {
				"ge": ` + strconv.FormatInt(prevBlock.CreatedAt, 10) + `,
				"lt": ` + strconv.FormatInt(currentBlock.CreatedAt, 10) + `
			}
		}`),
	}

	res, err := l.conn.Client().NetQueryCollection(&params)

	if err != nil {
		return nil, err
	}

	if len(res.Result) <= 0 {
		return nil, nil
	}

	message := &Message{}

	s, err := json.Marshal(res.Result[0])
	if err != nil {
		return nil, err
	}

	json.Unmarshal(s, &message)

	return message, nil
}

// registerEventHandler enables a handler for a given event. This cannot be used after Start is called.
// func (l *listener) registerEventHandler(name eventName, handler eventHandler) error {
// if l.subscriptions[name] != nil {
// 	return fmt.Errorf("event %s already registered", name)
// }
// l.subscriptions[name] = handler
// return nil
// }

// var ErrBlockNotReady = errors.New("required result to be 32 bytes, but got 0")

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")

	// AccessControllerABI := loadAbi("AccessController")
	// ReceiverABI := loadAbi("Receiver")
	SenderABI := loadAbi("Sender")

	currentBlockNumber := l.cfg.startBlock
	var prevBlockNumber *big.Int
	prevBlockNumber = new(big.Int)

	prevBlock, err := l.getBlock(prevBlockNumber.Sub(l.cfg.startBlock, big.NewInt(1)))

	if err != nil {
		return err
	}

	// var retry = BlockRetryLimit

	address := "0:dee8cdbf9937431376dd7ab7ee93367c14c62acc24d1d558cdd01186cf45704d"

	for {
		select {
		case <-l.stop:
			return errors.New("terminated")
		default:
			currentBlock, err := l.getBlock(currentBlockNumber)
			if err != nil {
				time.Sleep(BlockRetryInterval)
				continue
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Failed to query latest block", "block", currentBlock, "err", err)
				time.Sleep(BlockRetryInterval)
				continue
			}

			if currentBlock.Number > latestBlock.Number {
				l.log.Trace("Block not yet finalized", "target", currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			// message, err := l.getMessage(address, prevBlock, currentBlock)

			// if err != nil {
			// 	l.log.Error(err.Error(), "block", currentBlock, "err", err)
			// 	time.Sleep(BlockRetryInterval)
			// 	continue
			// }

			// body, err := l.decodeMessageBody(message, AccessControllerABI)
			// if err != nil {
			// 	l.log.Error("Failed to decode message", "block", currentBlock, "err", err)
			// 	time.Sleep(BlockRetryInterval)
			// 	continue
			// }

			message, err := l.getMessage(address, prevBlock, currentBlock)
			if err != nil {
				l.log.Error(err.Error(), "block", currentBlock, "err", err)
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Handle founded message
			if message != nil {
				body, err := l.decodeMessageBody(message, SenderABI)
				if err != nil {
					l.log.Error("Failed to decode message", "block", currentBlock, "err", err)
					time.Sleep(BlockRetryInterval)
					continue
				}

				fmt.Printf("Message: %v\n", body)
			} else {
				l.log.Error("No messages found", "block", currentBlock, "err", err)
			}

			prevBlock = currentBlock

			currentBlockNumber.Add(currentBlockNumber, big.NewInt(1))
			l.latestBlock.Height = big.NewInt(currentBlock.Number)
			l.latestBlock.LastUpdated = time.Now()

			time.Sleep(time.Second)
		}
	}
}

// processEvents fetches a block and parses out the events, calling Listener.handleEvents()
// func (l *listener) processEvents(hash types.Hash) error {
// l.log.Trace("Fetching block for events", "hash", hash.Hex())
// meta := l.conn.getMetadata()
// key, err := types.CreateStorageKey(&meta, "System", "Events", nil, nil)
// if err != nil {
// 	return err
// }

// var records types.EventRecordsRaw
// _, err = l.conn.api.RPC.State.GetStorage(key, &records, hash)
// if err != nil {
// 	return err
// }

// e := utils.Events{}
// err = records.DecodeEventRecords(&meta, &e)
// if err != nil {
// 	return err
// }

// l.handleEvents(e)
// l.log.Trace("Finished processing events", "block", hash.Hex())

// return nil
// }

// handleEvents calls the associated handler for all registered event types
// func (l *listener) handleEvents(evts utils.Events) {
// if l.subscriptions[FungibleTransfer] != nil {
// 	for _, evt := range evts.ChainBridge_FungibleTransfer {
// 		l.log.Trace("Handling FungibleTransfer event")
// 		l.submitMessage(l.subscriptions[FungibleTransfer](evt, l.log))
// 	}
// }
// if l.subscriptions[NonFungibleTransfer] != nil {
// 	for _, evt := range evts.ChainBridge_NonFungibleTransfer {
// 		l.log.Trace("Handling NonFungibleTransfer event")
// 		l.submitMessage(l.subscriptions[NonFungibleTransfer](evt, l.log))
// 	}
// }
// if l.subscriptions[GenericTransfer] != nil {
// 	for _, evt := range evts.ChainBridge_GenericTransfer {
// 		l.log.Trace("Handling GenericTransfer event")
// 		l.submitMessage(l.subscriptions[GenericTransfer](evt, l.log))
// 	}
// }

// if len(evts.System_CodeUpdated) > 0 {
// 	l.log.Trace("Received CodeUpdated event")
// 	err := l.conn.updateMetatdata()
// 	if err != nil {
// 		l.log.Error("Unable to update Metadata", "error", err)
// 	}
// }
// }

// submitMessage inserts the chainId into the msg and sends it to the router
// func (l *listener) submitMessage(m msg.Message, err error) {
// if err != nil {
// 	log15.Error("Critical error processing event", "err", err)
// 	return
// }
// m.Source = l.chainId
// err = l.router.Send(m)
// if err != nil {
// 	log15.Error("failed to process event", "err", err)
// }
// }
