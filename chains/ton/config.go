// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ByKeks/chainbridge-utils/core"
	"github.com/ByKeks/chainbridge-utils/msg"
)

// Config encapsulates all necessary parameters in ethereum compatible forms
type Config struct {
	name           string      // Human-readable chain name
	id             msg.ChainId // ChainID
	endpoint       string      // url for rpc endpoint
	from           string      // address of key to use
	keystorePath   string      // Location of keyfiles
	blockstorePath string
	freshStart     bool // Disables loading from blockstore at start
	gasLimit       *big.Int
	maxGasPrice    *big.Int
	http           bool // Config for type of connection
	startBlock     *big.Int
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func parseChainConfig(chainCfg *core.ChainConfig) (*Config, error) {

	config := &Config{
		name:           chainCfg.Name,
		id:             chainCfg.Id,
		endpoint:       chainCfg.Endpoint,
		from:           chainCfg.From,
		keystorePath:   chainCfg.KeystorePath,
		blockstorePath: chainCfg.BlockstorePath,
		freshStart:     chainCfg.FreshStart,
		http:           false,
		startBlock:     big.NewInt(0),
	}

	if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "true" {
		config.http = true
		delete(chainCfg.Opts, "http")
	} else if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "false" {
		config.http = false
		delete(chainCfg.Opts, "http")
	}

	if startBlock, ok := chainCfg.Opts["startBlock"]; ok && startBlock != "" {
		block := big.NewInt(0)
		_, pass := block.SetString(startBlock, 10)
		if pass {
			config.startBlock = block
			delete(chainCfg.Opts, "startBlock")
		} else {
			return nil, errors.New("unable to parse start block")
		}
	}

	if len(chainCfg.Opts) != 0 {
		return nil, fmt.Errorf("unknown Opts Encountered: %#v", chainCfg.Opts)
	}

	return config, nil
}
