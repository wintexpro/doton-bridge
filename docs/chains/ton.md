# TON Implementation Specification

Ton implementation of bridge should consist of some set of contracts: Bridge, BridgeVoteController, Proposal, Handler. As a source chain of transfer flow, ton implementation has Sender and Receiver contracts.

## Transfer Flow
1. Some users calls `sendData` function of `Sender` contract.
2. `Sender` contract calls `Receiver` by given address.
3. `Receiver` increment nonce by given chainId and emits `DataReceived` event
4. Relayers parse the `DataReceived` event and retrieve the associated record from the handler to construct a message.

### As Source Chain

### As Destination Chain
1. A Relayer calls `voteThrougthBridge` method of `Relayer` contract. This method calls `Bridge` contract method, named `relayerVoteForProposal`
2. `relayerVoteForProposal` method of `Bridge` contract calls `voteByBridge` method of `BridgeVoteController` contracts with some proposal data.
3. If a `Proposal` contract corresponding with the parameters passed does not created (deployed), it is created and the Relayer's vote is recorded into. If the proposal already exists, the Relayer's vote is simply recorded.
4. Upon passing the proposal threshold, `Proposal` cals `executeProposal` method of `Handler`, which in turn execute a ProposalExecuted event (outbound external message)

## Relayer Contract
Relayers will interact via `Relayer` contract, which is essentially private relayer interaction facade. 
```
function voteThroughBridge(uint8 choice, uint8 chainId, bytes32 messageType, uint64 nonce, bytes32 data, uint256 proposalPublicKey) onlyOwner external
```

## Bridge Contract
`Bridge` contract is a holder of bridge component addresses and also this components executor. Bridge validate, that relayer method was called by `Relayer` with valid code.
```
function relayerVoteForProposal(uint8 choice, uint8 chainId, bytes32 messageType, uint64 nonce, bytes32 data, uint256 relayerPubKey, uint256 proposalPublicKey) isValidRelayer(relayerPubKey) external
```

## BridgeVoteController
`BridgeVoteController` is a bridge component responsible for `Proposal` creation and execution (votes saving). Creation of proposal is a deploying a `Proposal` contract.

```
function voteByBridge(address voter, uint8 choice, uint8 chainId, bytes32 messageType, address handlerAddress, uint64 nonce, bytes32 data, uint256 proposalPublicKey) external

function createProposal(uint8 chainId, uint64 nonce, bytes32 data, uint8 initializerChoice, address initializerAddress, uint256 proposalPublicKey, address handlerAddress, bytes32 messageType) public returns (address proposalAddress)
```

## Proposal
`Proposal` contract stores relayer votes and execute handler

```
function voteByController(address voter, uint8 choice, bytes32 messageType, address handlerAddress) external
```

## Handler
`Handler` is a final component of bridge. This contracts execute event with proposal results data

```
event ProposalExecuted(uint8 chainId, uint64 nonce, bytes32 messageType, bytes32 data);
```

## Sender
`Sender` is a contract for sending some data for transfering it via bridge

```
function sendData(IReceiver destination, bool bounce, uint128 value, bytes32 data, uint256 destinationChainId) external onlyOwner
```

## Receiver
`Receiver` is a message emitter for relayers. another task for `Receiver` is a chain id nonces storing

```
function receiveData(bytes32 data, uint256 destinationChainId) external
```