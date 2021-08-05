// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"fmt"
	"testing"

	testnet "github.com/ChainSafe/chainlink-cosmos/testutil/network"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

var testTokens = sdk.NewIntWithDecimal(1000, 18)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// for generate test tx
	clientCtx client.Context

	cfg     testnet.Config
	testnet *testnet.Network

	account *testAccount
}

type testAccount struct {
	Priv   cryptotypes.PrivKey
	Addr   sdk.AccAddress
	Cosmos string
}

func (s *KeeperTestSuite) SetupTest() {
	s.account = generateKeyPair()
	s.T().Log("setting up test suite")

	cfg := testnet.DefaultConfig()
	genesisState := cfg.GenesisState
	cfg.NumValidators = 1

	var authData authtypes.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(genesisState[authtypes.ModuleName], &authData))

	genAccount, err := codectypes.NewAnyWithValue(&authtypes.BaseAccount{
		Address:       s.account.Addr.String(),
		AccountNumber: 1,
		Sequence:      0,
	})
	s.Require().NoError(err)
	authData.Accounts = append(authData.Accounts, genAccount)

	// configure genesis data for chainlink module
	var chainlinkData types.GenesisState
	chainlinkData.ModuleOwners = [](*types.MsgModuleOwner){&types.MsgModuleOwner{Address: s.account.Addr, PubKey: []byte(s.account.Cosmos), AssignerAddress: nil}}

	chainlinkDataBz, err := cfg.Codec.MarshalJSON(&chainlinkData)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = chainlinkDataBz

	s.cfg = cfg

	s.testnet = testnet.New(s.T(), cfg)
	_, err = s.testnet.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) BreakDownTest() {
	s.testnet.WaitForNextBlock()
	s.T().Log("Breaking down test suite.")
	s.testnet.Cleanup()
}

func generateKeyPair() *testAccount {
	priv, pub, addr := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pub)
	if err != nil {
		panic(err)
	}
	return &testAccount{priv, addr, cosmosPubKey}
}

func setupGRPCConn() *grpc.ClientConn {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	return grpcConn
}

func (s *KeeperTestSuite) TestGRPCQueryModuleOwner() {
	// var (
	// 	req *types.GetModuleOwnerRequest
	// 	res *types.GetModuleOwnerResponse
	// )

	val := s.testnet.Validators[0]
	baseURL := val.APIAddress
	// grpcConn := setupGRPCConn()
	// queryClient := types.NewQueryClient(grpcConn)

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		expErr   bool
		respType *types.GetModuleOwnerResponse
		expected *types.GetModuleOwnerResponse
	}{
		{
			"grpc query for GetAllModuleOwner",
			fmt.Sprintf("%s/chainlink.v1beta.Query/GetAllModuleOwner", baseURL),
			map[string]string{},
			false,
			&types.GetModuleOwnerResponse{},
			&types.GetModuleOwnerResponse{[](*types.MsgModuleOwner){&types.MsgModuleOwner{Address: s.account.Addr, PubKey: []byte(s.account.Cosmos), AssignerAddress: nil}}},
		},
	}
	//

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			resp, err := sdktestutil.GetRequestWithHeaders(tc.url, tc.headers)

			fmt.Println("tc.url: ", tc.url)
			fmt.Println("baseURL: ", baseURL)
			fmt.Println(fmt.Sprintf("RESP: %s", resp))
			/*
				RESP: {
					"code": 12,
					"message": "Not Implemented",
					"details": [
					]
				}
			*/
			s.Require().NoError(err)
			err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			}
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
