# client/transport

## Overview

Client-side transports and RPC wiring to communicate with Phantom servers. Builds authenticated gRPC clients, TLS dialers, and reconnect logic. Core logic addresses mTLS within the transport package.

## Go Files

- `mtls.go` – Configures the mutual TLS gRPC transport, including dial options and credential refresh logic.
