#!/usr/bin/env bash

./scripts/start.sh > "$(pwd)"/chainlinkd.log 2>&1 &
sleep 10

GRPC_INTEGRATION_TEST=true go test -v -run "^\QTestGRPCTestSuite\E$/^\QTestIntegration\E$" ./tests/grpc

pkill chainlinkd
echo "Chainlink GRPC tests has exited successfully."
exit 0
