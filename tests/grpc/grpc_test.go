// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package grpc

import (
	"context"
	"testing"
	"time"

	testnet "github.com/ChainSafe/chainlink-cosmos/testutil/network"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type testAccount struct {
	Name   string
	Priv   cryptotypes.PrivKey
	Pub    cryptotypes.PubKey
	Addr   sdk.AccAddress
	Cosmos string
}

var (
	alice = generateKeyPair("alice")
	bob   = generateKeyPair("bob")
	cerlo = generateKeyPair("cerlo")
	//alice *testAccount
	//bob   *testAccount
	//cerlo *testAccount
)

func generateKeyPair(name string) *testAccount {
	priv, pub, addr := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pub)
	if err != nil {
		panic(err)
	}
	return &testAccount{name, priv, priv.PubKey(), addr, cosmosPubKey}
}

func formatKeyPair(info keyring.Info) *testAccount {
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, info.GetPubKey())
	if err != nil {
		panic(err)
	}
	return &testAccount{
		Name:   info.GetName(),
		Pub:    info.GetPubKey(),
		Addr:   info.GetAddress(),
		Cosmos: cosmosPubKey,
	}
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}

type GRPCTestSuite struct {
	suite.Suite

	ctx       sdk.Context
	clientCtx client.Context

	config   testnet.Config
	network  *testnet.Network
	grpcConn *grpc.ClientConn
}

func (s *GRPCTestSuite) SetupTest() {
	s.T().Log("setup test suite")

	s.config = testnet.DefaultConfig()
	s.config.NumValidators = 1

	//configure genesis data for auth module
	var authGenState authtypes.GenesisState
	s.config.Codec.MustUnmarshalJSON(s.config.GenesisState[authtypes.ModuleName], &authGenState)
	genAccounts, err := authtypes.UnpackAccounts(authGenState.Accounts)
	s.Require().NoError(err)
	genAccounts = append(genAccounts, &authtypes.BaseAccount{
		Address:       alice.Addr.String(),
		AccountNumber: 1,
		Sequence:      0,
	})
	genAccounts = append(genAccounts, &authtypes.BaseAccount{
		Address:       bob.Addr.String(),
		AccountNumber: 2,
		Sequence:      0,
	})
	genAccounts = append(genAccounts, &authtypes.BaseAccount{
		Address:       cerlo.Addr.String(),
		AccountNumber: 3,
		Sequence:      0,
	})
	accounts, err := authtypes.PackAccounts(genAccounts)
	s.Require().NoError(err)
	authGenState.Accounts = accounts
	s.config.GenesisState[authtypes.ModuleName] = s.config.Codec.MustMarshalJSON(&authGenState)

	// configure genesis data for bank module
	balances := sdk.NewCoins(
		sdk.NewCoin("link", s.config.AccountTokens),
		sdk.NewCoin(s.config.BondDenom, s.config.StakingTokens),
	)
	var bankGenState banktypes.GenesisState
	s.config.Codec.MustUnmarshalJSON(s.config.GenesisState[banktypes.ModuleName], &bankGenState)
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: alice.Addr.String(),
		Coins:   balances.Sort(),
	})
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: bob.Addr.String(),
		Coins:   balances.Sort(),
	})
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: cerlo.Addr.String(),
		Coins:   balances.Sort(),
	})
	s.config.GenesisState[banktypes.ModuleName] = s.config.Codec.MustMarshalJSON(&bankGenState)

	// configure genesis data for chainlink module
	var chainlinkGenState types.GenesisState
	chainlinkGenState.ModuleOwners = []*types.MsgModuleOwner{{Address: alice.Addr, PubKey: []byte(alice.Cosmos), AssignerAddress: nil}}
	s.config.GenesisState[types.ModuleName] = s.config.Codec.MustMarshalJSON(&chainlinkGenState)

	s.network = testnet.New(s.T(), s.config)
	//_, err := s.network.WaitForHeight(1)
	//s.Require().NoError(err)
	//
	//kb := s.network.Validators[0].ClientCtx.Keyring
	//
	//_, err = kb.SavePubKey(alice.Name, alice.Pub, hd.Secp256k1.Name())
	//s.Require().NoError(err)
	//
	//_, err = kb.SavePubKey(bob.Name, bob.Pub, hd.Secp256k1.Name())
	//s.Require().NoError(err)
	//
	//_, err = kb.SavePubKey(cerlo.Name, cerlo.Pub, hd.Secp256k1.Name())
	//s.Require().NoError(err)

	//aliceAcc, _, err := kb.NewMnemonic("alice", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	//s.Require().NoError(err)
	//alice = formatKeyPair(aliceAcc)
	//
	//bobAcc, _, err := kb.NewMnemonic("bob", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	//s.Require().NoError(err)
	//alice = formatKeyPair(bobAcc)
	//
	//cerloAcc, _, err := kb.NewMnemonic("cerlo", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	//s.Require().NoError(err)
	//alice = formatKeyPair(cerloAcc)

	s.grpcConn, err = grpc.Dial(
		s.network.Validators[0].AppConfig.GRPC.Address,
		grpc.WithInsecure(),
	)
	require.NoError(s.T(), err)
}

func (s *GRPCTestSuite) BroadcastTx(ctx context.Context, submitter *testAccount, msgs ...sdk.Msg) *tx.BroadcastTxResponse {
	txClient := tx.NewServiceClient(s.grpcConn)

	encCfg := simapp.MakeTestEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msgs...)
	s.Require().NoError(err)

	txBuilder.SetGasLimit(testdata.NewTestGasLimit())
	//txBuilder.SetFeeAmount(...)
	//txBuilder.SetMemo(...)
	//txBuilder.SetTimeoutHeight(...)

	s.Require().NoError(s.Sign(submitter, txBuilder))

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)

	res, err := txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	s.Require().NoError(err)

	return res
}

func (s *GRPCTestSuite) Sign(signer *testAccount, txBuilder client.TxBuilder) error {
	//accNum, accSeq, err := s.clientCtx.AccountRetriever.GetAccountNumberSequence(s.clientCtx, signer.Addr)
	//s.Require().NoError(err)

	_, accSeq := uint64(1), uint64(0)

	encCfg := simapp.MakeTestEncodingConfig()

	//signerData := xauthsigning.SignerData{
	//	ChainID:       "testchain",
	//	AccountNumber: accNum,
	//	Sequence:      accSeq,
	//}

	signMode := encCfg.TxConfig.SignModeHandler().DefaultMode()
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   signer.Pub,
		Data:     &sigData,
		Sequence: accSeq,
	}

	err := txBuilder.SetSignatures(sig)
	s.Require().NoError(err)

	//bytesToSign, err := encCfg.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	//s.Require().NoError(err)

	//sigBytes, _, err := s.clientCtx.Keyring.Sign("alice", bytesToSign) // FIXME
	//s.Require().NoError(err)

	sigData = signing.SingleSignatureData{
		SignMode: signMode,
		//Signature: sigBytes,
		Signature: nil,
	}
	sig = signing.SignatureV2{
		PubKey:   signer.Pub,
		Data:     &sigData,
		Sequence: accSeq,
	}

	return txBuilder.SetSignatures(sig)
}

func (s *GRPCTestSuite) TestIntegration() {
	ctx := context.Background()
	queryClient := types.NewQueryClient(s.grpcConn)

	// 1 - Check initial module owner
	getModuleOwnerResponse, err := queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	s.Require().NoError(err)
	moduleOwner := getModuleOwnerResponse.GetModuleOwner()
	s.Require().Equal(1, len(moduleOwner))
	s.Require().Equal(alice.Addr, moduleOwner[0].GetAddress())
	s.Require().Equal(alice.Cosmos, string(moduleOwner[0].GetPubKey()))

	// 2 - Add new module owner by alice
	addModuleOwnerTx := &types.MsgModuleOwner{
		Address:         bob.Addr,
		PubKey:          []byte(bob.Cosmos),
		AssignerAddress: alice.Addr,
	}
	s.Require().NoError(addModuleOwnerTx.ValidateBasic())
	addModuleOwnerTxResponse := s.BroadcastTx(ctx, alice, addModuleOwnerTx)
	s.Require().EqualValues(0, addModuleOwnerTxResponse.TxResponse.Code)

	time.Sleep(5 * time.Second)

	getModuleOwnerResponse, err = queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	s.Require().NoError(err)
	moduleOwner = getModuleOwnerResponse.GetModuleOwner()
	s.Require().Equal(2, len(moduleOwner))

	// 3 - Module ownership transfer by bob to alice
	moduleOwnershipTransferTx := &types.MsgModuleOwnershipTransfer{
		NewModuleOwnerAddress: alice.Addr,
		NewModuleOwnerPubKey:  []byte(alice.Cosmos),
		AssignerAddress:       bob.Addr,
	}
	s.Require().NoError(addModuleOwnerTx.ValidateBasic())
	moduleOwnershipTransferTxResponse := s.BroadcastTx(ctx, bob, moduleOwnershipTransferTx)
	s.Require().EqualValues(0, moduleOwnershipTransferTxResponse.TxResponse.Code)

	// 4 - Add new feed by alice
	feedId := "testfeed1"
	addFeedTx := &types.MsgFeed{
		FeedId:    feedId,
		FeedOwner: cerlo.Addr,
		DataProviders: []*types.DataProvider{
			{
				Address: cerlo.Addr,
				PubKey:  []byte(cerlo.Cosmos),
			},
		},
		SubmissionCount:           1,
		HeartbeatTrigger:          2,
		DeviationThresholdTrigger: 3,
		FeedReward:                4,
		ModuleOwnerAddress:        alice.Addr,
	}
	s.Require().NoError(addFeedTx.ValidateBasic())
	addFeedTxResponse := s.BroadcastTx(ctx, alice, addFeedTx)
	s.Require().EqualValues(0, addFeedTxResponse.TxResponse.Code)

	time.Sleep(5 * time.Second)

	getFeedByFeedIdResponse, err := queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed := getFeedByFeedIdResponse.GetFeed()

	s.Require().Equal(feedId, feed.GetFeedId())
	s.Require().EqualValues(1, feed.GetSubmissionCount())
	s.Require().EqualValues(2, feed.GetHeartbeatTrigger())
	s.Require().EqualValues(3, feed.GetDeviationThresholdTrigger())
	s.Require().EqualValues(4, feed.GetFeedReward())
}
