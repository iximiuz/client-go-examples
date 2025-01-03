# Kubernetes client-go examples

```diff
! Support development of this project > patreon.com/iximiuz
```

A collection of mini-programs demonstrating various [client-go](https://github.com/kubernetes/client-go) use cases augmented by a [preconfigured online development environment](https://labs.iximiuz.com/playgrounds/k8s-client-go/). Inspired by [client-go/examples](https://github.com/kubernetes/client-go/tree/master/examples).

The intention is to test a (more or less) fresh version of Go and `k8s.io` packages against the [currently maintained Kubernetes release branches](https://kubernetes.io/releases/).

What is tested at the moment:

- `go 1.22.10`
- `k8s.io/client-go 0.28.14 0.29.12 0.30.8 0.31.4` (maintained release branches)
- `Kubernetes 1.28.9 1.29.12 1.30.8 1.31.4` (best-effort match with versions supported by `kind`)

## Setup

Most examples expect at least two Kubernetes clusters - `shared1` and `shared2`.

```bash
curl -sLS https://get.arkade.dev | sudo sh
arkade get kind kubectl

kind create cluster --name shared1
kind create cluster --name shared2
```

## Run

Oversimplified (for now):

```bash
cd <program>
make test

# or from the root folder:
make test-all
```

## TODO

- Add more assertions to mini-programs
- Examples to be covered
  - setting API request timeout
  - configuring API request throttling
  - `delete`
  - `delete collection`
  - `list` filtration
  - `watch` filtration
  - `informer` filtration
  - `patch` with different strategies
  - `Server Side Apply` (SSA)
  - working with subresources
  - `ownerReference` (one and many)
  - optimistic locking
  - https://stackoverflow.com/questions/56115197/how-to-idiomatically-fill-empty-fields-with-default-values-for-kubernetes-api-ob


## Contribution

Contributions are always welcome! Want to participate but don't know where to start? The TODO list above could give you some ideas.
Before jumping to the code, please create an issue describing the addition/change first. This will allow me to coordinate the effort
and make sure multiple people don't work on the same task.
