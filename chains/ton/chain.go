// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only
/*
The ton package contains the logic for interacting with ton chains.

There are 3 major components: the connection, the listener, and the writer.
The currently supported transfer types are Fungible (ERC20), Non-Fungible (ERC721), and generic.

Connection

The connection contains the ton RPC client and can be accessed by both the writer and listener.

Listener

The listener polls for each new block and looks for deposit events in the bridge contract. If a deposit occurs, the listener will fetch additional information from the handler before constructing a message and forwarding it to the router.

Writer

The writer recieves the message and creates a proposals on-chain. Once a proposal is made, the writer then watches for a finalization event and will attempt to execute the proposal if a matching event occurs. The writer skips over any proposals it has already seen.
*/
package ton

import (
	"github.com/ByKeks/chainbridge-utils/core"
	metrics "github.com/ByKeks/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/log15"
)

// var _ core.Chain = &Chain{}

// var _ Connection = &connection.Connection{}

// type Connection interface {
// 	Connect() error
// 	Keypair() *secp256k1.Keypair
// 	Opts() *bind.TransactOpts
// 	CallOpts() *bind.CallOpts
// 	LockAndUpdateOpts() error
// 	UnlockOpts()
// 	Client() *ethclient.Client
// 	EnsureHasBytecode(address common.Address) error
// 	LatestBlock() (*big.Int, error)
// 	WaitForBlock(block *big.Int) error
// 	Close()
// }

type Chain struct {
	cfg *core.ChainConfig // The config of the chain
	// conn Connection        // The chains connection
	// listener *listener         // The listener of this chain
	// writer   *writer           // The writer of the chain
	// stop chan<- int
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
// func setupBlockstore(cfg *Config, kp *secp256k1.Keypair) (*blockstore.Blockstore, error) {
// 	bs, err := blockstore.NewBlockstore(cfg.blockstorePath, cfg.id, kp.Address())
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !cfg.freshStart {
// 		latestBlock, err := bs.TryLoadLatestBlock()
// 		if err != nil {
// 			return nil, err
// 		}

// 		if latestBlock.Cmp(cfg.startBlock) == 1 {
// 			cfg.startBlock = latestBlock
// 		}
// 	}

// 	return bs, nil
// }

func InitializeChain(chainCfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error, m *metrics.ChainMetrics) (*Chain, error) {
	// cfg, err := parseChainConfig(chainCfg)
	// if err != nil {
	// 	return nil, err
	// }

	// keystore.KeypairFromAddress(cfg.from, keystore.TonChain, cfg.keystorePath, chainCfg.Insecure)
	// if err != nil {
	// 	return nil, err
	// }
	// kp, _ := kpI.(*secp256k1.Keypair)

	// bs, err := setupBlockstore(cfg, kp)
	// if err != nil {
	// 	return nil, err
	// }

	// stop := make(chan int)
	// conn := connection.NewConnection(cfg.endpoint, cfg.http, kp, logger, cfg.gasLimit, cfg.maxGasPrice)
	// err = conn.Connect()
	// if err != nil {
	// 	return nil, err
	// }

	// if chainCfg.LatestBlock {
	// 	curr, err := conn.LatestBlock()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	cfg.startBlock = curr
	// }

	// listener := NewListener(conn, cfg, logger, bs, stop, sysErr, m)
	// listener.setContracts(bridgeContract, erc20HandlerContract, erc721HandlerContract, genericHandlerContract)

	// writer := NewWriter(conn, cfg, logger, stop, sysErr, m)
	// writer.setContract(bridgeContract)

	return &Chain{
		cfg: chainCfg,
		// conn: conn,
		// writer:   writer,
		// listener: listener,
		// stop:     stop,
	}, nil
}

// func (c *Chain) SetRouter(r *core.Router) {
// 	r.Listen(c.cfg.Id, c.writer)
// 	c.listener.setRouter(r)
// }

// func (c *Chain) Start() error {
// 	err := c.listener.start()
// 	if err != nil {
// 		return err
// 	}

// 	err = c.writer.start()
// 	if err != nil {
// 		return err
// 	}

// 	c.writer.log.Debug("Successfully started chain")
// 	return nil
// }

// func (c *Chain) Id() msg.ChainId {
// 	return c.cfg.Id
// }

// func (c *Chain) Name() string {
// 	return c.cfg.Name
// }

// func (c *Chain) LatestBlock() metrics.LatestBlock {
// 	return c.listener.latestBlock
// }

// // Stop signals to any running routines to exit
// func (c *Chain) Stop() {
// 	close(c.stop)
// 	if c.conn != nil {
// 		c.conn.Close()
// 	}
// }
