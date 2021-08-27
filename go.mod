module github.com/ChainSafe/chainlink-cosmos

go 1.16

require (
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/ethereum/go-ethereum v1.10.6
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/smartcontractkit/libocr v0.0.0-20210803133922-ddddd3dce7e5
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/spm v0.0.0-20210524110815-6d7452d2dc4a
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.27.1 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
