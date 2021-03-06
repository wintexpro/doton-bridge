# Installation

## Relay Node

### Dependencies

- [go 1.15](https://golang.org/dl/)

- [Subkey](https://github.com/paritytech/substrate): 
Used for substrate key management. Only required if connecting to a substrate chain.

- [TON-SDK](https://github.com/tonlabs/TON-SDK):  `export CGO_LDFLAGS="-L/.../TON-SDK/target/release/deps/ -lton_client"`

### Building from Source

To build `doton` in `./build`.
```
$ make build
```

**or**

Use `go install` to add `doton` to your GOBIN.

```
$ make install
```

## Substrate Doton Node

### Dependencies

- [Rust Developer Environment](https://substrate.dev/docs/en/knowledgebase/getting-started/)

- Rust Nightly Toolchain

```
$ rustup install nightly-2020-10-01
$ rustup default nightly-2020-10-01
$ rustup target add wasm32-unknown-unknown --toolchain nightly-2020-10-01
```

### Building from Source

```
$ git clone https://github.com/wintexpro/doton-substrate-chain.git
$ cargo build --release
```

