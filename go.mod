module github.com/ChainSafe/chainlink-cosmos

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.46.8
	github.com/ethereum/go-ethereum v1.10.17
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	github.com/tendermint/spm v0.0.0-20210524110815-6d7452d2dc4a
	github.com/tendermint/tendermint v0.34.24
	github.com/tendermint/tm-db v0.6.7
	google.golang.org/genproto v0.0.0-20221014213838-99cd37c6964a
	google.golang.org/grpc v1.50.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
