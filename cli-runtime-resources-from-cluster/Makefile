CUR_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


.PHONY: test
test: go-mod-tidy
	go run ${CUR_DIR}/main.go -n kube-system svc

.PHONY: go-mod-tidy
go-mod-tidy:
	cd ${CUR_DIR} && go mod tidy