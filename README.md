# DOTON Bridge

This is a fork of the [original ChainSafe repository](https://github.com/ChainSafe/ChainBridge) to which the functionality of interacting with the TON node and some additional business logic were added.

# Contents

- [DOTON Bridge](#doton-bridge)
- [Contents](#contents)
- [Getting Started](#getting-started)
- [Installation](#installation)
  - [Dependencies](#dependencies)
  - [Building](#building)
  - [Docker](#docker)
- [Configuration](#configuration)
    - [TON Options](#ton-options)
    - [Substrate Options](#substrate-options)
  - [Blockstore](#blockstore)
  - [Keystore](#keystore)
- [Chain Implementations](#chain-implementations)
- [Docs](#docs)

# Getting Started

- Check out our [documentation](https://wintexpro.github.io/doton-bridge/).
- Try [running DOTON Bridge locally](https://wintexpro.github.io/doton-bridge/local/).
- Chat with us on telegram [[EN](https://t.me/doton_bridge)] [[RU](https://t.me/doton_bridge_ru)]

# Installation

## Dependencies

- DOTON Bridge use [ton-client-go](https://github.com/radianceteam/ton-client-go) and require [Rust](https://www.rust-lang.org/tools/install) and [TON-SDK](https://github.com/tonlabs/TON-SDK) (see [radianceteam/ton-client-go](https://github.com/radianceteam/ton-client-go)) Only required if connecting to a ton chain.

- [Subkey](https://substrate.dev/docs/en/knowledgebase/integrate/subkey): Used for substrate key management. Only required if connecting to a substrate chain.

## Building

`make build`: Builds `doton` in `./build`.

**or**

`make install`: Uses `go install` to add `doton` to your GOBIN.

## Docker 
The official wintex Docker image can be found [here](https://hub.docker.com/repository/docker/wintex/doton-bridge).

To build the Docker image locally run:

```
docker build -t wintex/doton-bridge .
```

To start DotonBridge:

``` 
docker run -v ./config.json:/config.json wintex/doton-bridge
```

# Configuration

A chain configurations take this form:

```
{
    "name": "freeTON",                  // Human-readable name
    "type": "ton",                      // Chain type (eg. "ton" or "substrate")
    "id": "0",                          // Chain ID
    "endpoint": "ws://<host>:<port>",   // Node endpoint
    "from": "0:164d61e...",             // On-chain address of relayer
    "opts": {},                         // Chain-specific configuration options (see below)
}
```

See `config.json.example` for an example configuration. 

### TON Options

Ton chains support the following additional options:
```
{

    "contractsPath": "/contracts", // The path to contract files (ABI, TVC)
    "receiver": "0:e50f...92ee",   // The contract Reciver address (Deploy script in /scripts/Makefile target: ton-deploy-contracts)
    "startBlock": "1",             // The block to start processing events from (default: 0)
    "workchainID": "0"             // The workchain from which the events will be processing
}
```

### Substrate Options

Substrate supports the following additonal options:

```
{
    "startBlock": "1234" // The block to start processing events from (default: 0)
    "useExtendedCall": "true" // Extend extrinsic calls to substrate with ResourceID. Used for backward compatibility with example pallet.
}
```

## Blockstore

The blockstore is used to record the last block the relayer processed, so it can pick up where it left off. 

If a `startBlock` option is provided (see [Configuration](#configuration)), then the greater of `startBlock` and the latest block in the blockstore is used at startup.

To disable loading from the blockstore specify the `--fresh` flag. A custom path for the blockstore can be provided with `--blockstore <path>`. For development, the `--latest` flag can be used to start from the current block and override any other configuration.

## Keystore

DOTON Bridge requires keys to sign and submit transactions, and to identify each bridge node on chain.

To use secure keys, see `doton accounts --help`. The keystore password can be supplied with the `KEYSTORE_PASSWORD` environment variable.

To import external ton keys, such as those generated with tonos-cli, use `doton accounts import --ton /path/to/key`.

or

`doton accounts import --ton --seedphrase "action glow era all liquid critic achieve lawsuit era anger loud slight"`

To import private keys as keystores, use `doton account import --privateKey key`.

For testing purposes, doton provides 5 test keys for substrate. The can be used with `--testkey <name>`, where `name` is one of `Alice`, `Bob`, `Charlie`, `Dave`, or `Eve`.
# Chain Implementations

- TON: [doton-ton](chains/ton)

    The contracts for interact with DOTON protocol

- Substrate: [doton-substrate](chains/substrate)

    A substrate pallet that can be integrated into a chain, as well as an example pallet to demonstrate chain integration.

# Docs

MKdocs will generate static HTML files for DOTON markdown files located in `./docs/`

`make install-mkdocs`: Pull the docker image MkDocs

`make mkdocs`: Run MkDoc's docker image, building and hosting the html files on `localhost:8000`  
