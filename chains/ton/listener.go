// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/ChainBridge/chains"
	"github.com/ChainSafe/ChainBridge/connections/ton"
	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/wintexpro/chainbridge-utils/blockstore"
	metrics "github.com/wintexpro/chainbridge-utils/metrics/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

type listener struct {
	cfg         Config
	conn        *ton.Connection
	blockstore  blockstore.Blockstorer
	latestBlock metrics.LatestBlock
	log         log15.Logger
	router      chains.Router
	stop        <-chan int
	sysErr      chan<- error
	abi         map[string]client.Abi
}

type Message struct {
	ID        string  `json:"id"`
	Status    int8    `json:"status"`
	CreatedAt big.Int `json:"created_at"`
	Body      string  `json:"body"`
	Src       string  `json:"src"`
}

// Frequency of polling for a new block
var BlockRetryInterval = time.Second * 5

var BlockRetryLimit = 10

func NewListener(conn *ton.Connection, blockstore blockstore.Blockstorer, cfg *Config, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	abi := make(map[string]client.Abi)

	for _, subscription := range Subscriptions {
		abi[subscription.abiName] = LoadAbi(cfg.contractsPath, subscription.abiName)
	}

	return &listener{
		cfg:         *cfg,
		conn:        conn,
		blockstore:  blockstore,
		log:         log,
		latestBlock: metrics.LatestBlock{LastUpdated: time.Now()},
		stop:        stop,
		sysErr:      sysErr,
		abi:         abi,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// start creates the initial subscription for all events
func (l *listener) start() error {
	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")

	currentBlockNumber := l.cfg.startBlock
	var prevBlockNumber *big.Int
	prevBlockNumber = new(big.Int)

	prevBlock, err := GetBlock(l.conn.Client(), prevBlockNumber.Sub(l.cfg.startBlock, big.NewInt(1)), l.cfg.workchainID)

	if err != nil {
		return err
	}

	var retry = BlockRetryLimit

	for {
		select {
		case <-l.stop:
			return errors.New("terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.sysErr <- fmt.Errorf("event polling retries exceeded (chain=%d, name=%s)", l.cfg.id, l.cfg.name)
				return nil
			}

			currentBlock, err := GetBlock(l.conn.Client(), currentBlockNumber, l.cfg.workchainID)
			if err != nil {
				l.log.Debug("Failed to query current block", "block", currentBlock, "err", err)
				time.Sleep(BlockRetryInterval)
				continue
			}

			l.log.Debug("Block: " + currentBlockNumber.String() + " is being processed")

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Failed to query latest block", "block", latestBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			if currentBlock.Number > latestBlock.Number {
				l.log.Trace("Block not yet finalized", "target", currentBlock, "latest", latestBlock)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processEvents(prevBlock, currentBlock)
			if err != nil {
				l.log.Error(err.Error(), "block", currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(big.NewInt(0).SetUint64(currentBlock.Number))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}

			prevBlock = currentBlock

			l.log.Debug("Block: " + currentBlockNumber.String() + " is done")

			currentBlockNumber.Add(currentBlockNumber, big.NewInt(1))
			l.latestBlock.Height = big.NewInt(int64(currentBlock.Number))
			l.latestBlock.LastUpdated = time.Now()
			retry = BlockRetryLimit
		}
	}
}

// processEvents fetches a block and parses out the events, calling Listener.handleEvents()
func (l *listener) processEvents(prevBlock, currentBlock *ton.BlockType) error {
	for _, subscription := range Subscriptions {
		messages, err := GetMessage(l.conn.Client(), l.cfg.contracts[subscription.contractKey], prevBlock, currentBlock)
		if err != nil {
			return err
		}

		// Handle founded messages
		for _, rawMessage := range *messages {
			message := Message{}

			err := json.Unmarshal(rawMessage, &message)
			if err != nil {
				return err
			}

			rawBody := json.RawMessage(message.Body)

			l.log.Debug(fmt.Sprintf("Attemping decode message: %#v", message))

			bodyRaw, err := DecodeMessageBody(l.conn.Client(), &rawBody, l.abi[subscription.abiName])
			if err != nil {
				return err
			}

			msg, err := subscription.handler(message, bodyRaw.Value, l.log)

			if err != nil {
				log15.Error("Critical error processing event", "err", err)
				return err
			}

			l.submitMessage(msg)
		}
	}

	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *msg.Message) {
	m.Source = l.cfg.id

	l.log.Info("Trying send message", "Source", m.Source, "Destination", m.Destination, "ResourceId", fmt.Sprintf("%x", m.ResourceId), "DepositNonce", m.DepositNonce)
	err := l.router.Send(*m)

	if err != nil {
		log15.Error("failed to process event", "err", err)
	}
}
