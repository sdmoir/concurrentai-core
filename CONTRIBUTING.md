# Contributing

Thank you for your interest in contributing to Concurrent.ai! This document is very much a work in progress, so please feel free to make suggestions about how it can be improved.

## Running Locally

To run a local Concurrent.ai stack, follow the Getting Started guide in [concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra) for running a local stack via Minikube.

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

We are using [GitHub flow](https://guides.github.com/introduction/flow/) as our branching strategy (minimal structure with changes branched off of `main`).

All external changes should be created from a [personal fork](https://guides.github.com/activities/forking/) of concurrentai/concurrentai-core, and then submitted as a pull request.

When creating your PR:
- Make sure to target the `main` branch
- Include a descriptive summary of the change, a link to any issue associated with the change, etc.
- Make sure that all tests are passing

Happy coding!
