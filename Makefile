PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

CHAINLINK_DAEMON_BINARY = chainlinkd

###############################################################################
###                                  Build&Run                                  ###
###############################################################################

all: update install start

update:
	${GO_MOD} go mod tidy
	${GO_MOD} go mod vendor

install:
	${GO_MOD} go install ./cmd/$(CHAINLINK_DAEMON_BINARY)

start:
	./scripts/start.sh

clean:
	@rm -rf ./vendor

###############################################################################
###                                   Lint                                  ###
###############################################################################

.PHONY: help lint test license
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: get-lint
	./bin/golangci-lint run ./... --timeout 5m0s

###############################################################################
###                                Check&Testing                            ###
###############################################################################

check:
	gosec ./...

test:
	go test --race ./...

test-addFeed:
	./scripts/addFeed.sh

###############################################################################
###                                   Protobuf                              ###
###############################################################################

proto-install:
	./scripts/proto-tools-installer.sh

protogen:
	./scripts/protocgen

###############################################################################
###                                   License                               ###
###############################################################################

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 ./x

## license-check: Checks for missing license headers
license-check:
	@echo "  >  \033[Checking for license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -check -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 ./x