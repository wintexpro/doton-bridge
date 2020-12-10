// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"errors"
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
	abi         map[string]client.Abi
}

// Frequency of polling for a new block
var BlockRetryInterval = time.Second * 5

// var BlockRetryLimit = 5

func NewListener(conn *ton.Connection, blockstore blockstore.Blockstorer, cfg *Config, log log15.Logger, stop <-chan int) *listener {
	abi := make(map[string]client.Abi)

	for _, subscription := range Subscriptions {
		abi[subscription.abiName] = LoadAbi(cfg.abiPath, subscription.abiName)
	}

	return &listener{
		cfg:         *cfg,
		conn:        conn,
		blockstore:  blockstore,
		log:         log,
		latestBlock: metrics.LatestBlock{LastUpdated: time.Now()},
		stop:        stop,
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

	prevBlock, err := GetBlock(l.conn.Client(), prevBlockNumber.Sub(l.cfg.startBlock, big.NewInt(1)))

	if err != nil {
		return err
	}

	// var retry = BlockRetryLimit

	for {
		select {
		case <-l.stop:
			return errors.New("terminated")
		default:
			currentBlock, err := GetBlock(l.conn.Client(), currentBlockNumber)
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

			err = l.processEvents(prevBlock, currentBlock)
			if err != nil {
				l.log.Error(err.Error(), "block", currentBlock, "err", err)
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(big.NewInt(0).SetUint64(currentBlock.Number))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}

			prevBlock = currentBlock

			currentBlockNumber.Add(currentBlockNumber, big.NewInt(1))
			l.latestBlock.Height = big.NewInt(int64(currentBlock.Number))
			l.latestBlock.LastUpdated = time.Now()

			time.Sleep(time.Second)
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
		for _, message := range *messages {
			body, err := DecodeMessageBody(l.conn.Client(), &message, l.abi[subscription.abiName])
			if err != nil {
				return err
			}

			l.submitMessage(subscription.handler(body, l.log))
		}
	}

	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m msg.Message, err error) {
	if err != nil {
		log15.Error("Critical error processing event", "err", err)
		return
	}

	m.Source = l.cfg.id
	err = l.router.Send(m)
	if err != nil {
		log15.Error("failed to process event", "err", err)
	}
}
