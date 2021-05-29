# Copyright 2020 Wintex
# SPDX-License-Identifier: LGPL-3.0-only

# Build TON SDK

FROM phusion/baseimage:0.11 as sdk-builder
ARG PROFILE=release
WORKDIR /ton

ENV RUST_TOOLCHAIN=nightly-2020-11-19

RUN apt-get update && \
	apt-get dist-upgrade -y -o Dpkg::Options::="--force-confold" && \
	apt-get install -y cmake pkg-config libssl-dev git clang
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y && \
	export PATH="$PATH:$HOME/.cargo/bin" && \
	rustup default stable
RUN git clone --depth 1 --branch 1.14.1 https://github.com/tonlabs/TON-SDK.git && \
  export PATH="$PATH:$HOME/.cargo/bin" && \
  cd TON-SDK && \
  rustup toolchain install $RUST_TOOLCHAIN && \
  rustup default $RUST_TOOLCHAIN && \
  cargo build "--$PROFILE"

# Build bridge

FROM golang:1.15.6-buster AS builder
ADD . /src
WORKDIR /src

COPY --from=sdk-builder /ton/TON-SDK /TON-SDK

ENV CGO_LDFLAGS="-L//TON-SDK/target/release/deps/ -lton_client"
ENV LD_LIBRARY_PATH="/TON-SDK/target/release/deps/"

RUN echo "deb http://security.ubuntu.com/ubuntu bionic-security main" >> /etc/apt/sources.list && \
  apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 3B4FE6ACC0B21F32 && \
  apt-get update && apt-cache policy libssl1.0-dev && \
  apt-get install -y libssl1.0-dev

RUN go mod download
RUN cd cmd/doton && go build -o /bridge .

# Make small image

# FROM debian:stretch-slim
FROM phusion/baseimage:0.11 AS final

COPY --from=sdk-builder /ton/TON-SDK /TON-SDK

ENV CGO_LDFLAGS="-L//TON-SDK/target/release/deps/ -lton_client"
ENV LD_LIBRARY_PATH="/TON-SDK/target/release/deps/"

RUN apt-get update && \
	apt-get dist-upgrade -y -o Dpkg::Options::="--force-confold" && \
	apt-get install -y libssl-dev ca-certificates wget

RUN wget -P /usr/local/bin/ https://chainbridge.ams3.digitaloceanspaces.com/subkey-rc6 \
  && mv /usr/local/bin/subkey-rc6 /usr/local/bin/subkey \
  && chmod +x /usr/local/bin/subkey
RUN subkey --version

COPY --from=builder /bridge ./
RUN chmod +x ./bridge

ENTRYPOINT ["./bridge"]
