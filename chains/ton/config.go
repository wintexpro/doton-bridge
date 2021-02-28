// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/prometheus/common/log"
	"github.com/wintexpro/chainbridge-utils/core"
	"github.com/wintexpro/chainbridge-utils/msg"
)

// Chain specific options
var (
	SenderOpt              = "sender"
	ReceiverOpt            = "receiver"
	BurnedTokensHandlerOpt = "burnedTokensHandler"
	ContractsPath          = "contractsPath"
	WorkchainID            = "workchainID"
)

// Config encapsulates all necessary parameters in ethereum compatible forms
type Config struct {
	id             msg.ChainId // ChainID
	name           string      // Human-readable chain name
	endpoint       string      // url for rpc endpoint
	from           string      // address of key to use
	keystorePath   string      // Location of keyfiles
	blockstorePath string      // Location of blockstore
	contractsPath  string      // Location of abi files
	workchainID    string
	contracts      map[string]string
	freshStart     bool // Disables loading from blockstore at start
	http           bool // Config for type of connection
	startBlock     *big.Int
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func ParseChainConfig(chainCfg *core.ChainConfig) (*Config, error) {

	config := &Config{
		name:           chainCfg.Name,
		id:             chainCfg.Id,
		endpoint:       chainCfg.Endpoint,
		from:           chainCfg.From,
		keystorePath:   chainCfg.KeystorePath,
		blockstorePath: chainCfg.BlockstorePath,
		contracts:      make(map[string]string),
		freshStart:     chainCfg.FreshStart,
		http:           false,
		startBlock:     big.NewInt(0),
		workchainID:    "0",
	}

	if workchainID, ok := chainCfg.Opts[WorkchainID]; ok && workchainID != "" {
		config.workchainID = workchainID
		delete(chainCfg.Opts, WorkchainID)
	} else {
		log.Info("Workchain ID value set to the default: 0")
	}

	if contractsPath, ok := chainCfg.Opts[ContractsPath]; ok && contractsPath != "" {
		config.contractsPath = contractsPath
		delete(chainCfg.Opts, ContractsPath)
	} else {
		return nil, fmt.Errorf("must provide opts.contractsPath field for ton config")
	}

	if contract, ok := chainCfg.Opts[ReceiverOpt]; ok && contract != "" {
		config.contracts[ReceiverOpt] = contract
		delete(chainCfg.Opts, ReceiverOpt)
	} else {
		return nil, fmt.Errorf("must provide opts.receiver field for ton config")
	}

	if contract, ok := chainCfg.Opts[BurnedTokensHandlerOpt]; ok && contract != "" {
		config.contracts[BurnedTokensHandlerOpt] = contract
		delete(chainCfg.Opts, BurnedTokensHandlerOpt)
	} else {
		return nil, fmt.Errorf("must provide opts.burnedTokensHandler field for ton config")
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
