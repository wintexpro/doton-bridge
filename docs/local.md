# Running Locally
# doton-local-network ([Github](https://github.com/wintexpro/doton-local-network))

A docker development environment for a doton local network

## Requirements

 - [docker](https://docs.docker.com/engine/install/)
 - [docker-compose](https://docs.docker.com/compose/install/)

## Getting Started

Clone this repo:
```sh
$ git clone git@github.com:wintexpro/doton-local-network.git
```

Inside `./doton-local-network` directory run:

```sh
$ make run-chains
```

After making sure TON OS SE and substrate services are running, you can run the environment setup script

```sh
$ make run-setup
```

That execute to send initial value to contracts and deploy them to local TON node, also set up a substrate node for interacting with DOTON protocol. After setup, this command attaches you to a docker container with [tonos-cli](https://github.com/tonlabs/tonos-cli) and [halva-cli](https://github.com/halva-suite/halva) tools.

At last, you can run the relay node by running the command (please, make sure the setup scripts were done):

```sh
$ make run-bridge
```

## Polkadot JS Apps

You can interact with a substrate local node by visiting https://polkadot.js.org/apps/ and choose "Local Node" network

You will need to add these definitions to the [developer settings](https://polkadot.js.org/apps/#/settings/developer):

```json
{
  "Message": "Text",
  "chainbridge::ChainId": "u8",
  "ExtAddress": "Text",
  "ChainId": "u8",
  "ResourceId": "[u8; 32]",
  "Nonce": "u64",
  "DepositNonce": "u64",
  "ProposalVotes": {
    "votes_for": "Vec<AccountId>",
    "votes_against": "Vec<AccountId>",
    "status": "enum"
  },
  "Erc721Token": {
    "id": "TokenId",
    "metadata": "Vec<u8>"
  },
  "TokenId": "U256",
  "Address": "AccountId",
  "LookupSource": "AccountId"
}
```

## FreeTON Local Node

TONOS Startup Edition (SE) is a pre-configured Docker image with a local blockchain that provides the same API as a Dapp Server.

You can interact with a FreeTON local node by visiting http://127.0.0.1/graphql

## Helpers

To send a message to Substrate through TON, you must follow a helper command:

```sh 
make ton-send-msg MSG="Hello substrate\!"
```

And to send a message to TON through Substrate:

```sh 
make sub-send-msg MSG="Hello ton\!"
```

## Configuration

This repository is fully complete to working with the bridge, but if you need to change any config files, you can do so without rebuilding docker images. Directories `./scripts`, `./configs`, `./contracts`, `./keys` will be mounted to docker containers

#### Keys:

The Repository contains prepared TON keys:

`./keys/0:df22eba0b48020b70efa7a6e9d6360ed1dc20877250947470cc1289b14c9cc1e.key` - The test relay node key derived from the seed phrase: `"action glow era all liquid critic achieve lawsuit era anger loud slight"`, also same key `0:df22eba0b48020b70efa7a6e9d6360ed1dc20877250947470cc1289b14c9cc1e.tonos.key` in tonos format.

And Substrate sr25519 key:
`5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY.key` derived from Uri `//Alice`

You can replace them with your own keys if you need.

#### Configuration files

##### configs/config.json
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

###### Ton Options:

Ton chains support the following additional options:

```
{

    "contractsPath": "/contracts", // The path to contract files (ABI, TVC)
    "receiver": "0:e50f...92ee",   // The contract Reciver address (Deploy script in /scripts/Makefile target: ton-deploy-contracts)
    "startBlock": "1",             // The block to start processing events from (default: 0)
    "workchainID": "0"             // The workchain from which the events will be processing
}
```

###### Substrate Options

Substrate supports the following additonal options:

```
{
    "startBlock": "1234" // The block to start processing events from (default: 0)
    "useExtendedCall": "true" // Extend extrinsic calls to substrate with ResourceID. Used for backward compatibility with example pallet.
}
```

###### configs/halva.js
```
const mnemonic = "bottom drive obey lake curtain smoke basket hold race lonely fit walk"; // don't used in this case

module.exports = {
  networks: {
    test: {
      mnemonic,
      ws: "ws://doton-sub-chain:9944", // substrate endpoint
    },
  },
  polkadotjs: {
    types: {
      "Message": "Text",
      "chainbridge::ChainId": "u8",
      "ExtAddress": "Text",
      "ChainId": "u8",
      "ResourceId": "[u8; 32]",
      "Nonce": "u64",
      "DepositNonce": "u64",
      "ProposalVotes": {
        "votes_for": "Vec<AccountId>",
        "votes_against": "Vec<AccountId>",
        "status": "enum"
      },
      "Erc721Token": {
        "id": "TokenId",
        "metadata": "Vec<u8>"
      },
      "TokenId": "U256",
      "Address": "AccountId",
      "LookupSource": "AccountId"
    }
  }
}
```