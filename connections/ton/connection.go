// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ByKeks/chainbridge-utils/crypto/ed25519"
	"github.com/ChainSafe/log15"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

var BlockRetryInterval = time.Second * 5

type BlockType struct {
	Number int64 `json:"seq_no"`
}

type Connection struct {
	endpoint string
	http     bool
	kp       *ed25519.Keypair
	conn     *client.Client
	// opts     *bind.TransactOpts
	// callOpts *bind.CallOpts
	nonce uint64
	// optsLock sync.Mutex
	log  log15.Logger
	stop chan int // All routines should exit when this channel is closed
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(endpoint string, http bool, kp *ed25519.Keypair, log log15.Logger) *Connection {
	return &Connection{
		endpoint: endpoint,
		http:     http,
		kp:       kp,
		log:      log,
		stop:     make(chan int),
	}
}

func (c *Connection) Keypair() *ed25519.Keypair {
	return c.kp
}

func (c *Connection) Client() *client.Client {
	return c.conn
}

// func (c *Connection) Opts() *bind.TransactOpts {
// 	return c.opts
// }

// func (c *Connection) CallOpts() *bind.CallOpts {
// 	return c.callOpts
// }

// Connect starts the ton connection
func (c *Connection) Connect() error {
	c.log.Info("Connecting to ton chain...", "url", c.endpoint)

	conn, err := client.NewClient(client.Config{
		Network: &client.NetworkConfig{ServerAddress: c.endpoint},
	})
	if err != nil {
		return err
	}

	c.conn = conn

	// Construct tx opts, call opts, and nonce mechanism
	// opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	// if err != nil {
	// 	return err
	// }
	// c.opts = opts
	// c.nonce = 0
	// c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
// func (c *Connection) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) {
// privateKey := c.kp.PrivateKey()
// address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

// nonce, err := c.conn.PendingNonceAt(context.Background(), address)
// if err != nil {
// 	return nil, 0, err
// }

// auth := bind.NewKeyedTransactor(privateKey)
// auth.Nonce = big.NewInt(int64(nonce))
// auth.Value = value
// auth.GasLimit = uint64(gasLimit.Int64())
// auth.GasPrice = gasPrice
// auth.Context = context.Background()

// return auth, nonce, nil
// }

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
func (c *Connection) LockAndUpdateOpts() error {
	// c.optsLock.Lock()

	// nonce, err := c.conn.PendingNonceAt(context.Background(), c.opts.From)
	// if err != nil {
	// 	c.optsLock.Unlock()
	// 	return err
	// }
	// c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Connection) UnlockOpts() {
	// c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (*big.Int, error) {
	params := client.ParamsOfQueryCollection{
		Collection: "blocks",
		Result:     "seq_no",
		Limit:      null.Uint32From(1),
		Filter:     json.RawMessage(`{"workchain_id":{"eq":-1}, "status":{"eq": 2}}`),
		Order: []client.OrderBy{{
			Path:      "seq_no",
			Direction: client.DescSortDirection,
		}},
	}

	res, err := c.conn.NetQueryCollection(&params)
	if err != nil {
		return nil, err
	}

	if len(res.Result) <= 0 {
		return nil, errors.New("No blocks found")
	}

	latestBlock := &BlockType{}

	s, err := json.Marshal(res.Result[0])
	if err != nil {
		return nil, err
	}

	json.Unmarshal(s, &latestBlock)

	return big.NewInt(latestBlock.Number), nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Connection) EnsureHasBytecode(addr ethcommon.Address) error {
	// code, err := c.conn.CodeAt(context.Background(), addr, nil)
	// if err != nil {
	// 	return err
	// }

	// if len(code) == 0 {
	// 	return fmt.Errorf("no bytecode found at %s", addr.Hex())
	// }
	return nil
}

// WaitForBlock will poll for the block number until the current block is equal or greater than
func (c *Connection) WaitForBlock(block *big.Int) error {
	// for {
	// 	select {
	// 	case <-c.stop:
	// 		return errors.New("connection terminated")
	// 	default:
	// 		currBlock, err := c.LatestBlock()
	// 		if err != nil {
	// 			return err
	// 		}

	// 		// Equal or greater than target
	// 		if currBlock.Cmp(block) >= 0 {
	// 			return nil
	// 		}
	// 		c.log.Trace("Block not ready, waiting", "target", block, "current", currBlock)
	// 		time.Sleep(BlockRetryInterval)
	// 		continue
	// 	}
	// }
	return nil
}

// Close terminates the client connection and stops any running routines
func (c *Connection) Close() {

	if c.conn != nil {
		defer c.conn.Close()
	}
}
