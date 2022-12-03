# Kubernetes client-go examples

```diff
! Support development of this project > patreon.com/iximiuz
```

A collection of mini-programs covering various [client-go](https://github.com/kubernetes/client-go) use cases inspired by [client-go/examples](https://github.com/kubernetes/client-go/tree/master/examples).
The intention (at least so far) is to test (more or less) fresh version of Go and packages against a few latest
Kubernetes versions.

What tested at the moment:

- `go 1.18`
- `k8s.io/client-go v0.23.1 v0.24.7 v0.25.3`
- `Kubernetes v1.22.15 1.23.12 v1.24.6 v1.25.2`

## Setup

Most examples expect `minikube` with at least two Kubernetes clusters - `shared1` and `shared2`.

```bash
curl -sLS https://get.arkade.dev | sudo sh
arkade get minikube kubectl

minikube start --profile shared1
minikube start --profile shared2
```

## Run

Oversimplified (for now):

```bash
cd <program>
go run main.go
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
