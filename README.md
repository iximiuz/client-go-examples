# Kubernetes client-go examples

A collection of mini-programs covering various [client-go](https://github.com/kubernetes/client-go) use cases inspired by [client-go/examples](https://github.com/kubernetes/client-go/tree/master/examples).
The intention (at least so far) is to test (more or less) fresh version of Go and packages against a few latest
Kubernetes versions.

What tested at the moment:

- `go 1.17`
- `k8s.io/client-go v0.23.1`
- `Kubernetes v1.22.3`

## Setup

All examples expect `minikube` with at least two Kubernetes clusters - `shared1` and `shared2`.

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
- Test different Kubernetes versions
- Setup GitHub action(s)
- Examples to be covered
  - `delete` + owner references
  - `delete collection`
  - `list` filtration
  - `watch` filtration
  - `informer` filtration
  - `patch` with different strategies
  - `workqueue` - controllers' fundamentals
  - `retry.RetryOnConflict`