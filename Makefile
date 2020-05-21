SHELL := /bin/bash

GO111MODULE = on

mock:
	mockgen -package=internal -destination=dht/internal/mocks/application_mock.go -source=dht/node.go
.PHONY: mock

proto:
	protoc --gogo_out=./pb/ --proto_path=./pb/ pastry.proto
.PHONY: proto

