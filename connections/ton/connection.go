// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

var BlockRetryInterval = time.Second * 5

type BlockType struct {
	Number    uint64 `json:"seq_no"`
	CreatedAt int64  `json:"gen_utime"`
}

type Connection struct {
	workchainID string
	endpoint    string
	http        bool
	conn        *client.Client
	log         log15.Logger
	stop        chan int // All routines should exit when this channel is closed
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(endpoint string, http bool, workchainID string, log log15.Logger) *Connection {
	return &Connection{
		workchainID: workchainID,
		endpoint:    endpoint,
		http:        http,
		log:         log,
		stop:        make(chan int),
	}
}

func (c *Connection) Client() *client.Client {
	return c.conn
}

// Connect starts the ton connection
func (c *Connection) Connect() error {
	c.log.Info("Connecting to ton chain...", "url", c.endpoint)

	conn, err := client.NewClient(client.Config{
		Network: &client.NetworkConfig{ServerAddress: null.StringFrom(c.endpoint)},
	})
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (*BlockType, error) {
	params := client.ParamsOfQueryCollection{
		Collection: "blocks",
		Result:     "seq_no gen_utime",
		Limit:      null.Uint32From(1),
		Filter: json.RawMessage(`{
			"workchain_id":{"eq":` + c.workchainID + `},
			"status":{"eq": 2}
		}`),
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

	return latestBlock, nil
}

// Close terminates the client connection and stops any running routines
func (c *Connection) Close() {

	if c.conn != nil {
		defer c.conn.Close()
	}
}
