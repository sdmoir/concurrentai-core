# Contributing

Thank you for your interest in contributing to Concurrent.ai! This document is very much a work in progress, so please feel free to make suggestions about how it can be improved.

## Overview

Who this document is for...

## Running Locally

Follow the Getting Started guide in [concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra) to run a local Concurrent.ai stack via Minikube.

## Coding

The core Concurrent.ai components are written in [Go](https://golang.org/).

General coding guidelines:
- [Scannability and readability](https://www.geepawhill.org/2019/03/20/refactoring-pro-tip-i-optimize-scannability-then-readability-then-writability/) are top priorities
- Test-driven development is encouraged but not required

## Testing

### Unit/Micro Tests

```bash
# model-enricher
cd src/model-enricher && go test ./...

# model-executor
cd src/model-executor && go test ./...

# rendezvous-api
cd src/rendezvous-api && go test ./...

# rendezvous-collector
cd src/rendezvous-collector && go test ./...
```

### Integration Tests

_Coming Soon_


## Submitting a Pull Request
- When creating your PR, target the `master` branch
- Once tests are passing, request a review from a [Concurrent.ai team member](https://github.com/orgs/concurrentai/people)
