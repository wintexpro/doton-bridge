# CLI Options

## Flags

### Global

```
--config value       JSON configuration file
--verbosity value    Supports levels crit (silent) to trce (trace) (default: "info")
--keystore value     Path to keystore directory (default: "./keys")
--blockstore value   Specify path for blockstore
--fresh              Disables loading from blockstore at start. Opts will still be used if specified. (default: false)
--latest             Overrides blockstore and start block, starts from latest block (default: false)
--metrics            Enables metric server (default: false)
--metricsPort value  Port to serve metrics on (default: 8001)
--testkey value      Applies a predetermined test keystore to the chains.
--help, -h           show help (default: false)
--version, -v        print the version (default: false)
```

### Deploy Management
#### `doton contracts send-grams`
send grams from the giver to the contracts
#### `doton contracts deploy`
deploy the set of ton contracts
#### `doton contracts deploy-relayer`
deploy the relayer ton contract
```
--from Address of relayer
```
#### `doton contracts deploy-wallet`
deploy the wallet contract
#### `doton contracts get-balance`
get of the wallet balance
#### `doton contracts setup`
setup contracts of brige

#### `doton contracts send-tokens`

transfer the balance from owner's account to another account

```
--amount  Amount of tokens
--to      Recipient
--nonce   Nonce
```

### Account Management

The commands can be used to manage keys in the local keystore. You can view available keys with `doton accounts list`.

#### `doton accounts generate`
```
--password value  Password used to encrypt the keystore. Used with --generate, --import, or --unlock
--sr25519         Specify account/key type as sr25519. (default: false)
--secp256k1       Specify account/key type as secp256k1. (default: false)
--ed25519         Specify account/key type as ed25519. (default: false)
--network value   Specify the network to use for the address encoding (substrate/polkadot/centrifuge) (default: substrate)
```

#### `doton accounts import`
```
--seedphrase value  Derive keys from a seed phrase and import them into a keystore.
--ton               Import an existing ton keystore. (default: false)
--privateKey value  Import a hex representation of a private key into a keystore.
--sr25519           Specify account/key type as sr25519. (default: false)
--secp256k1         Specify account/key type as secp256k1. (default: false)
--password value    Password used to encrypt the keystore. Used with --generate, --import, or --unlock
--network value     Specify the network to use for the address encoding (substrate/polkadot/centrifuge) (default: substrate)
```


## Environment Variables

- `KEYSTORE_PASSWORD`: The password to use when loading the keystore.
- `BLOCK_TIMEOUT`: The duration (seconds) until a chain is considered "unhealthy"