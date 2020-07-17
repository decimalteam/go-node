PACKAGES=$(shell go list ./... | grep -v '/simulation')

# VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION := '0.9.8'
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Decimal \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=decd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=deccli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/decd
		go install $(BUILD_FLAGS) ./cmd/deccli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify