# doton-local-network

A docker development environment for a doton local network

## Requirements

 - [docker](https://docs.docker.com/engine/install/)
 - [docker-compose](https://docs.docker.com/compose/install/)

## Launching bridge components

If you got some old versions of doton components, make sure that you got latest docker images on your local environment. You can just check your latest tagged images digests with digests on [dockerhub](https://hub.docker.com/u/wintex). Or just remove all local images, because this guide's commands would download them if not find.

Clone this repo:
```sh
$ git clone git@github.com:wintexpro/doton-local-network.git
```

Inside `./doton-local-network` directory run:

```sh
$ make run-chains
```

After making sure TON OS SE and substrate services are running, you should run setup script

```sh
$ make run-setup
```

That executes to send initial value to contracts and deploy them to local TON node, also set up a substrate node for interacting with DOTON protocol. After setup, this command attaches you to a docker container with [tonos-cli](https://github.com/tonlabs/tonos-cli) and [halva-cli](https://github.com/halva-suite/halva) tools.

Make sure that containers which were launched during previous commands (doton-setup-bridge and doton-setup) stopped their work.

You need to run the relay nodes by running the command (please, make sure the setup scripts were done):

```sh
$ make run-alice
```

After making sure alice node is running, you must deploy another 2 relayers

At last, you can run

```sh
$ CONFIG_NAME=config2.json make deploy-relayer
$ CONFIG_NAME=config3.json make deploy-relayer
```

And run they

```sh
$ make run-bob
$ make run-charlie
```

## Run local polkadot UI

If you need to run local polkadotjs UI, you can run command:

```sh
$ docker run --rm -it --name polkadot-ui --network doton-local-network_default --link doton-sub-chain -p 8001:80 jacogr/polkadot-js-apps:0.79.1
```

## Setting up Polkadot JS Apps

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
  "LookupSource": "AccountId",
  "VrfResult": {
    "pk": "Vec<u8>",
    "val": "Vec<u8>",
    "proof": "Vec<u8>"
  }
}
```

## Setting substrate initial balance

For working example you should set some balance for account on substrate network. To do this you should perform some steps:
1. Open an https://polkadot.js.org/apps/
2. Make sure that local node was choosen
3. On https://polkadot.js.org/apps/ choose Accounts->Accounts->Add account. We reccomend you to create new account, because of existing accounts has got large amount of balance
4. Go throught Developers->sudo
5. On 'Sudo access' tab choose **balances** extrinsic with **setBalance** command
6. Choose created account for **who** option (Or choose any if you want or didn't create new)
7. Set **new_free** option with some amount (ex 10000) and click submit sudo.


## Substrate-TON transfer

1. On https://polkadot.js.org/apps/ choose developers->extrinsic
2. Set an *example* value for **submit the following extrinsic** option with *transferNative* function (to the right)
3. Set some **amount** . That value will be transfered to TON network
4. **recipient** field should be filled with TON recipient address. You can obtain it by launching this command from repository (doton-local-network) root
```
$ make get-balance
```
5. Set **dest_id** to *2* (chain id of TON network in `configs/config.json`)
6. Click **Submit transaction**
7. Unlock with password *123456* and click **Sign and Submit**

Wait for a while and check TON recipient balance with command from **4** . It should be equal with sended amount.

## TON-Substrate transfer

| Current version doesn't supports any user interface, excepts cli, so this steps can be complicated

For sending tokens to the substrate network you should use this command from repository root:

```
AMOUNT=<amount> TO=<substrate address> NONCE=<nonce> make ton-send-tokens
```

Parameters:

- \<amount\> - amount of tokens that will be sent
- \<substrate address\> - adress of substrate recepient
- \<nonce\> IMPORTANT! you should increment nonce by yourself. Start with *1*.

Example:
```
$ AMOUNT=10 TO=5CaNkosmQfEaLnYjdegNH5KSPAtPeMtkjkM4xAeQMv1gUGCp NONCE=1 make ton-send-tokens
```

Considering ton local node can't produce blocks without a messages, you should run several times previous command with --amount 0 and incremented nonce (one by one, start from previous + 1)

Example:
```
$ AMOUNT=0 TO=5CaNkosmQfEaLnYjdegNH5KSPAtPeMtkjkM4xAeQMv1gUGCp NONCE=2 make ton-send-tokens
$ AMOUNT=0 TO=5CaNkosmQfEaLnYjdegNH5KSPAtPeMtkjkM4xAeQMv1gUGCp NONCE=3 make ton-send-tokens
```

Wait for a while and check **Accounts** tab on https://polkadot.js.org/apps/. Balance is expected to increase. 