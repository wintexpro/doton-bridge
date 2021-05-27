// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

/*
The substrate package contains the logic for interacting with substrate chains.
The current supported transfer types are Fungible, Nonfungible, and generic.

There are 3 major components: the connection, the listener, and the writer.

Connection

The Connection handles connecting to the substrate client, and submitting transactions to the client.
It also handles state queries. The connection is shared by the writer and listener.

Listener

The substrate listener polls blocks and parses the associated events for the three transfer types. It then forwards these into the router.

Writer

As the writer receives messages from the router, it constructs proposals. If a proposal is still active, the writer will attempt to vote on it. Resource IDs are resolved to method name on-chain, which are then used in the proposals when constructing the resulting Call struct.

*/
package substrate

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/wintexpro/chainbridge-utils/blockstore"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/crypto/sr25519"
	"github.com/wintexpro/chainbridge-utils/keystore"
	metrics "github.com/wintexpro/chainbridge-utils/metrics/types"
	"github.com/wintexpro/chainbridge-utils/msg"
)

var _ core.Chain = &Chain{}

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection       // THe chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
	stop     chan<- int
	logger   log15.Logger
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than startBlock, then the latest block is returned, otherwise startBlock is.
func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}

func InitializeChain(cfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error, m *metrics.ChainMetrics) (*Chain, error) {
	kp, err := keystore.KeypairFromAddress(cfg.From, keystore.SubChain, cfg.KeystorePath, cfg.Insecure)
	if err != nil {
		return nil, err
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	// Attempt to load latest block
	bs, err := blockstore.NewBlockstore(cfg.BlockstorePath, cfg.Id, kp.Address())
	if err != nil {
		return nil, err
	}
	startBlock := parseStartBlock(cfg)
	if !cfg.FreshStart {
		startBlock, err = checkBlockstore(bs, startBlock)
		if err != nil {
			return nil, err
		}
	}

	stop := make(chan int)
	// Setup connection
	conn := NewConnection(cfg.Endpoint, cfg.Name, krp, logger, stop, sysErr)
	err = conn.Connect()
	if err != nil {
		return nil, err
	}

	err = conn.checkChainId(cfg.Id)
	if err != nil {
		return nil, err
	}

	if cfg.LatestBlock {
		curr, err := conn.api.RPC.Chain.GetHeaderLatest()
		if err != nil {
			return nil, err
		}
		startBlock = uint64(curr.Number)
	}

	ue := parseUseExtended(cfg)

	vrfkp, err := VrfGenerateKeypair()
	if err != nil {
		return nil, err
	}

	// Setup listener & writer
	l := NewListener(conn, cfg.Name, cfg.Id, startBlock, logger, bs, stop, sysErr, m)
	w := NewWriter(conn, logger, sysErr, m, ue, vrfkp)

	chain := &Chain{
		cfg:      cfg,
		conn:     conn,
		listener: l,
		writer:   w,
		stop:     stop,
		logger:   logger,
	}

	err = chain.sendVrfPublicKey(vrfkp)
	if err != nil {
		return nil, err
	}

	chain.checkEpoch(vrfkp)

	return chain, nil
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	c.conn.log.Debug("Successfully started chain", "chainId", c.cfg.Id)
	return nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
	c.listener.setRouter(r)
}

func (c *Chain) LatestBlock() metrics.LatestBlock {
	return c.listener.latestBlock
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) Stop() {
	close(c.stop)
}

func (c *Chain) getPublicRandomness(epochNumber int) (string, error) {
	epoch, err := types.EncodeToBytes(uint32(epochNumber))
	if err != nil {
		return "", err
	}

	data := c.conn.getMetadata()
	key, err := types.CreateStorageKey(&data, "DorrStorage", "EpochToRandomness", epoch, nil)
	if err != nil {
		return "", err
	}

	res, err := c.conn.api.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return "", err
	}
	if len(*res) == 0 {
		return "", errors.New("STORAGE IS EMPTY")
	}

	return hex.EncodeToString(*res), nil
}

func (c *Chain) sendVrfPublicKey(kp *Keypair) error {
	c.conn.log.Debug("Trying send VRF public key to network", "chainId", c.cfg.Id, "public", fmt.Sprintf("%x", kp.Public()))
	err := c.conn.SubmitTx("Dorr.set_pk", kp.Public())
	if err != nil {
		return err
	}

	c.conn.log.Debug("Successfully sended VRF public key to network", "chainId", c.cfg.Id)

	return nil
}

func (c *Chain) sendVrfResults(epoch int, vrfkp *Keypair) error {
	c.conn.log.Debug("Trying send VRF results to network", "chainId", c.cfg.Id)

	randomness, err := c.getPublicRandomness(epoch)
	if err != nil {
		return err
	}

	randomnessBytes, err := hex.DecodeString(randomness)
	if err != nil {
		return err
	}

	inout, proof, err := vrfkp.Sign(randomnessBytes)
	if err != nil {
		return err
	}

	err = c.conn.SubmitTx("Dorr.set_vrf_results", inout, proof)
	if err != nil {
		return err
	}

	c.conn.log.Debug("Successfully VRF results to network", "chainId", c.cfg.Id)

	return nil
}

func (c *Chain) checkEpoch(vrfkp *Keypair) {
	ticker := time.NewTicker(6 * time.Second)
	quit := make(chan struct{})
	accountID := types.NewAccountID(vrfkp.Public())

	go func() {
		for {
			select {
			case <-ticker.C:
				var epoch int
				err := c.conn.api.Client.Call(&epoch, "dorr_getCurrentEpoch")
				if err != nil {
					c.logger.Error("Failed call to dorr_getCurrentEpoch", "error", err)
					return
				}

				var epochByPk int
				err = c.conn.api.Client.Call(&epochByPk, "dorr_getEpochByPk", accountID)
				if err != nil {
					c.logger.Error("Failed call to dorr_getEpochByPk", "error", err)
					return
				}

				c.logger.Debug("Check current epoch", "epoch", epoch, "epochByPk", epochByPk)

				if epochByPk < epoch {
					c.sendVrfResults(epochByPk, vrfkp)
					ticker.Stop()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
