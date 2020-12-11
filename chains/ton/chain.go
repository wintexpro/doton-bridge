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
	"encoding/json"
	"errors"
	"math/big"

	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/blockstore"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	"github.com/wintexpro/chainbridge-utils/keystore"
	metrics "github.com/wintexpro/chainbridge-utils/metrics/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

// var _ core.Chain = &Chain{}

var _ Connection = &connection.Connection{}

type Connection interface {
	Connect() error
	Client() *client.Client
	LatestBlock() (*connection.BlockType, error)
	Close()
}

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     Connection        // The chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
	kp       *ed25519.Keypair
	stop     chan<- int
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
func setupBlockstore(cfg *Config, kp *ed25519.Keypair) (*blockstore.Blockstore, error) {
	bs, err := blockstore.NewBlockstore(cfg.blockstorePath, cfg.id, kp.Address())
	if err != nil {
		return nil, err
	}

	if !cfg.freshStart {
		latestBlock, err := bs.TryLoadLatestBlock()
		if err != nil {
			return nil, err
		}

		if latestBlock.Cmp(cfg.startBlock) == 1 {
			cfg.startBlock = latestBlock
		}
	}

	return bs, nil
}

func InitializeChain(chainCfg *core.ChainConfig, log log15.Logger, sysErr chan<- error, m *metrics.ChainMetrics) (*Chain, error) {
	cfg, err := parseChainConfig(chainCfg)
	if err != nil {
		return nil, err
	}

	kpI, err := keystore.KeypairFromAddress(cfg.from, keystore.TonChain, cfg.keystorePath, chainCfg.Insecure)
	if err != nil {
		return nil, err
	}
	kp, _ := kpI.(*ed25519.Keypair)

	bs, err := setupBlockstore(cfg, kp)
	if err != nil {
		return nil, err
	}

	stop := make(chan int)
	conn := connection.NewConnection(cfg.endpoint, cfg.http, log)
	err = conn.Connect()
	if err != nil {
		return nil, err
	}

	if chainCfg.LatestBlock {
		curr, err := conn.LatestBlock()
		if err != nil {
			return nil, err
		}
		cfg.startBlock = big.NewInt(int64(curr.Number))
	}

	c := &Chain{
		cfg:      chainCfg,
		conn:     conn,
		kp:       kp,
		writer:   NewWriter(conn, cfg, log, stop, sysErr, m),
		listener: NewListener(conn, bs, cfg, log, stop),
		stop:     stop,
	}

	err = c.checkAndDeploySender(cfg)

	return c, err
}

func (c *Chain) checkAndDeploySender(cfg *Config) error {
	keys := client.KeyPair{
		Public: c.kp.PublicKey(),
		Secret: c.kp.SecretKey(),
	}

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: keys,
	}

	SenderABI := LoadAbi(cfg.contractsPath, "Sender")
	SenderTVC := LoadTvc(cfg.contractsPath, "Sender")

	deploySet := client.DeploySet{
		Tvc:         SenderTVC,
		WorkchainID: null.NewInt32(0, true),
	}

	callSet := client.CallSet{
		FunctionName: "constructor",
	}

	paramsOfEncodeMsg := client.ParamsOfEncodeMessage{
		Abi:       SenderABI,
		DeploySet: &deploySet,
		CallSet:   &callSet,
		Signer:    signer,
	}

	encodedMessage, err := c.conn.Client().AbiEncodeMessage(&paramsOfEncodeMsg)
	if err != nil {
		return err
	}

	paramsOfQuery := client.ParamsOfQueryCollection{
		Collection: "accounts",
		Filter:     json.RawMessage(`{ "id": {"eq": "` + encodedMessage.Address + `"} }`),
		Result:     "id balance",
	}

	res, err := c.conn.Client().NetQueryCollection(&paramsOfQuery)
	if err != nil {
		return err
	}

	if len(res.Result) == 0 {
		return errors.New("Your sender address: " + encodedMessage.Address + " does not have balance for deploy the contract code")
	}

	paramsOfSendMessage := client.ParamsOfSendMessage{
		Message: encodedMessage.Message,
		Abi:     &SenderABI,
	}

	_, err = c.conn.Client().ProcessingSendMessage(&paramsOfSendMessage, MessageCallback)

	return nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
	c.listener.setRouter(r)
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	err = c.writer.start()
	if err != nil {
		return err
	}

	c.writer.log.Debug("Successfully started chain")
	return nil
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) LatestBlock() metrics.LatestBlock {
	return c.listener.latestBlock
}

// Stop signals to any running routines to exit
func (c *Chain) Stop() {
	close(c.stop)
	if c.conn != nil {
		c.conn.Close()
	}
}
