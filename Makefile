SHELL := /bin/bash

GO111MODULE = on

mock:
	mockgen -package=internal -destination=dht/internal/mocks/application_mock.go -source=dht/node.go
.PHONY: mock