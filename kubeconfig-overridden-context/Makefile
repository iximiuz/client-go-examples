CUR_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


.PHONY: test
test: go-mod-tidy
	kind create cluster --name cluster2
	go run ${CUR_DIR}/main.go kind-cluster2

.PHONY: go-mod-tidy
go-mod-tidy:
	cd ${CUR_DIR} && go mod tidy