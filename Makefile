SHELL := /bin/bash

GO111MODULE = on

mock:
	mockgen -package=internal -destination=dht/internal/mocks/dht_mocks.go -source=dht/dht.go
.PHONY: mock

proto:
	protoc --gogo_out=./pb/ --proto_path=./pb/ pastry.proto
.PHONY: proto

