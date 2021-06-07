PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

CHAINLINK_DAEMON_BINARY = chainlinkd

###############################################################################
###                                  Build                                  ###
###############################################################################

install:
	${GO_MOD} go install ./cmd/$(CHAINLINK_DAEMON_BINARY)

clean:
	@rm -rf ./vendor

.PHONY: help lint test
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
###                                Testing                                  ###
###############################################################################

test:
	go test ./...
