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

### Ton Options

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

ChainBridge requires keys to sign and submit transactions, and to identify each bridge node on chain.

To use secure keys, see `chainbridge accounts --help`. The keystore password can be supplied with the `KEYSTORE_PASSWORD` environment variable.

To import external ton keys, such as those generated with tonos-cli, use `chainbridge accounts import --ton /path/to/key.json`

or

`chainbridge accounts import --ton --seedphrase "action glow era all liquid critic achieve lawsuit era anger loud slight"`

To import private keys as keystores, use `chainbridge account import --privateKey key`.

For testing purposes, chainbridge provides 5 test keys. The can be used with `--testkey <name>`, where `name` is one of `Alice`, `Bob`, `Charlie`, `Dave`, or `Eve`. 
