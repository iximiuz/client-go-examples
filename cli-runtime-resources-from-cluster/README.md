# Simple tool to look up (and print out) Kubernetes objects by resource(s) and name(s)

```bash
go run main.go --help
go run main.go po
go run main.go pod
go run main.go pods
go run main.go services,deployments
go run main.go --namespace=default service/kubernetes
```
