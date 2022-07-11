PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := '0.10.3'
COMMIT = $(shell git rev-parse --short=8 HEAD)

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Decimal \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=decd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=deccli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'
BUILD_TAGS := -tags cleveldb

all: install

install: go.sum
		go install $(BUILD_FLAGS) $(BUILD_TAGS) ./cmd/decd
		go install $(BUILD_FLAGS) ./cmd/deccli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified "
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

test:
	@go test -mod=readonly $(PACKAGES)
