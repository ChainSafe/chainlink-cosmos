// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package grpc

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ChainSafe/chainlink-cosmos/app"
	testnet "github.com/ChainSafe/chainlink-cosmos/testutil/network"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	signing2 "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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
	//alice = generateKeyPair("alice")
	//bob   = generateKeyPair("bob")
	//cerlo = generateKeyPair("cerlo")
	alice *testAccount
	bob   *testAccount
	cerlo *testAccount
)

//func generateKeyPair(name string) *testAccount {
//	priv, pub, addr := testdata.KeyTestPubAddr()
//	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pub)
//	if err != nil {
//		panic(err)
//	}
//	return &testAccount{name, priv, priv.PubKey(), addr, cosmosPubKey}
//}

//func formatKeyPair(info keyring.Info) *testAccount {
//	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, info.GetPubKey())
//	if err != nil {
//		panic(err)
//	}
//	return &testAccount{
//		Name:   info.GetName(),
//		Pub:    info.GetPubKey(),
//		Addr:   info.GetAddress(),
//		Cosmos: cosmosPubKey,
//	}
//}

func importKeyPair(t testing.TB, clientCtx client.Context, name string) *testAccount {
	info, err := clientCtx.Keyring.Key(name)
	require.NoError(t, err)
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, info.GetPubKey())
	require.NoError(t, err)
	return &testAccount{
		Name:   info.GetName(),
		Pub:    info.GetPubKey(),
		Addr:   info.GetAddress(),
		Cosmos: cosmosPubKey,
	}
}

func TestGRPCTestSuite(t *testing.T) {
	running := os.Getenv("GRPC_INTEGRATION_TEST")
	if running != "true" {
		t.SkipNow()
	}
	suite.Run(t, new(GRPCTestSuite))
}

type GRPCTestSuite struct {
	suite.Suite

	clientCtx client.Context
	rpcClient rpcclient.Client
	config    testnet.Config
	network   *testnet.Network
	grpcConn  *grpc.ClientConn
}

// SetupTest directly connected to daemon on port 9090
// fresh `scripts/start.sh` need to be run before each execution
func (s *GRPCTestSuite) SetupTest() {
	s.T().Log("setup test suite")

	userHomeDir, err := os.UserHomeDir()
	require.NoError(s.T(), err)
	home := filepath.Join(userHomeDir, ".chainlinkd")

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithHomeDir(home).
		WithViper("").
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

	clientCtx, err = config.ReadFromClientConfig(clientCtx)
	require.NoError(s.T(), err)

	backendKeyring, err := client.NewKeyringFromBackend(clientCtx, "test")
	require.NoError(s.T(), err)

	clientCtx = clientCtx.WithKeyring(backendKeyring)

	alice = importKeyPair(s.T(), clientCtx, "alice")
	bob = importKeyPair(s.T(), clientCtx, "bob")
	cerlo = importKeyPair(s.T(), clientCtx, "cerlo")

	s.clientCtx = clientCtx

	s.rpcClient, err = client.NewClientFromNode("tcp://127.0.0.1:26657")
	require.NoError(s.T(), err)

	s.grpcConn, err = grpc.Dial(
		"127.0.0.1:9090",
		grpc.WithInsecure(),
	)
	require.NoError(s.T(), err)
}

// SetupTest using testnet package
// TODO find a way to have the auth and bank genesis state not overwritten by testnet
// generateKeyPair and formatKeyPair can be used to generate or format existing key
func (s *GRPCTestSuite) SetupTestTMP() {
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
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.grpcConn, err = grpc.Dial(
		s.network.Validators[0].AppConfig.GRPC.Address,
		grpc.WithInsecure(),
	)
	require.NoError(s.T(), err)
}

// TODO can use BuildTX from TxFactory
func (s *GRPCTestSuite) BroadcastTx(ctx context.Context, submitter *testAccount, msgs ...sdk.Msg) *tx.BroadcastTxResponse {
	txClient := tx.NewServiceClient(s.grpcConn)

	encCfg := simapp.MakeTestEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msgs...)
	s.Require().NoError(err)

	txBuilder.SetGasLimit(testdata.NewTestGasLimit())
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("link", 3)))
	//txBuilder.SetMemo(...)
	//txBuilder.SetTimeoutHeight(...)

	s.Require().NoError(s.Sign(submitter, txBuilder))

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)

	res, err := txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_BLOCK,
		TxBytes: txBytes,
	})
	s.Require().NoError(err)

	return res
}

// TODO can use BuildTX from TxFactory
func (s *GRPCTestSuite) Sign(signer *testAccount, txBuilder client.TxBuilder) error {
	accNum, accSeq, err := s.clientCtx.AccountRetriever.GetAccountNumberSequence(s.clientCtx, signer.Addr)
	s.Require().NoError(err)

	//_, accSeq := uint64(1), uint64(0)

	encCfg := simapp.MakeTestEncodingConfig()

	signerData := signing2.SignerData{
		ChainID:       "testchain",
		AccountNumber: accNum,
		Sequence:      accSeq,
	}

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

	err = txBuilder.SetSignatures(sig)
	s.Require().NoError(err)

	bytesToSign, err := encCfg.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	s.Require().NoError(err)

	sigBytes, _, err := s.clientCtx.Keyring.Sign(signer.Name, bytesToSign) // FIXME
	s.Require().NoError(err)

	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
		//Signature: nil,
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

	_, _ = s.waitForBlock(1)

	s.T().Log("1 - Check initial module owner")

	getModuleOwnerResponse, err := queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	s.Require().NoError(err)
	moduleOwner := getModuleOwnerResponse.GetModuleOwner()
	s.Require().Equal(1, len(moduleOwner))
	s.Require().Equal(alice.Addr, moduleOwner[0].GetAddress())
	s.Require().Equal(alice.Cosmos, string(moduleOwner[0].GetPubKey()))

	s.T().Log("2 - Add new module owner by alice")

	addModuleOwnerTx := &types.MsgModuleOwner{
		Address:         bob.Addr,
		PubKey:          []byte(bob.Cosmos),
		AssignerAddress: alice.Addr,
	}
	s.Require().NoError(addModuleOwnerTx.ValidateBasic())
	addModuleOwnerResponse := s.BroadcastTx(ctx, alice, addModuleOwnerTx)
	s.Require().EqualValues(0, addModuleOwnerResponse.TxResponse.Code)

	getModuleOwnerResponse, err = queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	s.Require().NoError(err)
	moduleOwner = getModuleOwnerResponse.GetModuleOwner()
	s.Require().Equal(2, len(moduleOwner))

	s.T().Log("3 - Module ownership transfer by bob to alice")

	moduleOwnershipTransferTx := &types.MsgModuleOwnershipTransfer{
		NewModuleOwnerAddress: alice.Addr,
		NewModuleOwnerPubKey:  []byte(alice.Cosmos),
		AssignerAddress:       bob.Addr,
	}
	s.Require().NoError(addModuleOwnerTx.ValidateBasic())
	moduleOwnershipTransferResponse := s.BroadcastTx(ctx, bob, moduleOwnershipTransferTx)
	s.Require().EqualValues(0, moduleOwnershipTransferResponse.TxResponse.Code)

	s.T().Log("4 - Add new feed by alice")

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
		SubmissionCount:           10,
		HeartbeatTrigger:          2,
		DeviationThresholdTrigger: 3,
		FeedReward: &types.FeedRewardSchema{
			Amount:   100,
			Strategy: "",
		},
		ModuleOwnerAddress: alice.Addr,
	}
	s.Require().NoError(addFeedTx.ValidateBasic())
	addFeedResponse := s.BroadcastTx(ctx, alice, addFeedTx)
	s.Require().EqualValues(0, addFeedResponse.TxResponse.Code)

	getFeedByFeedIdResponse, err := queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed := getFeedByFeedIdResponse.GetFeed()

	s.Require().Equal(feedId, feed.GetFeedId())
	s.Require().EqualValues(cerlo.Addr, feed.GetFeedOwner())
	s.Require().EqualValues(1, len(feed.GetDataProviders()))
	s.Require().Contains(feed.GetDataProviders(), &types.DataProvider{Address: cerlo.Addr, PubKey: []byte(cerlo.Cosmos)})
	s.Require().EqualValues(10, feed.GetSubmissionCount())
	s.Require().EqualValues(2, feed.GetHeartbeatTrigger())
	s.Require().EqualValues(3, feed.GetDeviationThresholdTrigger())
	s.Require().EqualValues(uint32(0x64), feed.GetFeedReward().GetAmount())
	s.Require().EqualValues("", feed.GetFeedReward().GetStrategy())

	s.T().Log("5 - Add data provider by cerlo")

	addDataProviderTx := &types.MsgAddDataProvider{
		FeedId: feedId,
		DataProvider: &types.DataProvider{
			Address: bob.Addr,
			PubKey:  []byte(bob.Cosmos),
		},
		Signer: cerlo.Addr,
	}
	s.Require().NoError(addDataProviderTx.ValidateBasic())
	addDataProviderResponse := s.BroadcastTx(ctx, cerlo, addDataProviderTx)
	s.Require().EqualValues(0, addDataProviderResponse.TxResponse.Code)

	getFeedByFeedIdResponse, err = queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed = getFeedByFeedIdResponse.GetFeed()
	s.Require().EqualValues(2, len(feed.GetDataProviders()))
	s.Require().Contains(feed.GetDataProviders(), &types.DataProvider{Address: cerlo.Addr, PubKey: []byte(cerlo.Cosmos)})
	s.Require().Contains(feed.GetDataProviders(), &types.DataProvider{Address: bob.Addr, PubKey: []byte(bob.Cosmos)})

	s.T().Log("6 - Remove data provider by cerlo")

	removeDataProviderTx := &types.MsgRemoveDataProvider{
		FeedId:  feedId,
		Address: cerlo.Addr,
		Signer:  cerlo.Addr,
	}
	s.Require().NoError(removeDataProviderTx.ValidateBasic())
	removeDataProviderResponse := s.BroadcastTx(ctx, cerlo, removeDataProviderTx)
	s.Require().EqualValues(0, removeDataProviderResponse.TxResponse.Code)

	getFeedByFeedIdResponse, err = queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed = getFeedByFeedIdResponse.GetFeed()
	s.Require().EqualValues(1, len(feed.GetDataProviders()))
	s.Require().Contains(feed.GetDataProviders(), &types.DataProvider{Address: bob.Addr, PubKey: []byte(bob.Cosmos)})

	s.T().Log("7 - Feed ownership transfer to bob by cerlo")

	feedOwnershipTransferTx := &types.MsgFeedOwnershipTransfer{
		FeedId:              feedId,
		NewFeedOwnerAddress: bob.Addr,
		Signer:              cerlo.Addr,
	}
	s.Require().NoError(feedOwnershipTransferTx.ValidateBasic())
	feedOwnershipTransferResponse := s.BroadcastTx(ctx, cerlo, feedOwnershipTransferTx)
	s.Require().EqualValues(0, feedOwnershipTransferResponse.TxResponse.Code)

	getFeedByFeedIdResponse, err = queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed = getFeedByFeedIdResponse.GetFeed()
	s.Require().EqualValues(bob.Addr, feed.GetFeedOwner())

	s.T().Log("8 - Update submission count parameter")

	setSubmissionCountTx := &types.MsgSetSubmissionCount{
		FeedId:          feedId,
		SubmissionCount: 1,
		Signer:          bob.Addr,
	}
	setHeartbeatTriggerTx := &types.MsgSetHeartbeatTrigger{
		FeedId:           feedId,
		HeartbeatTrigger: 200,
		Signer:           bob.Addr,
	}
	setDeviationThresholdTriggerTx := &types.MsgSetDeviationThresholdTrigger{
		FeedId:                    feedId,
		DeviationThresholdTrigger: 300,
		Signer:                    bob.Addr,
	}
	setFeedRewardTx := &types.MsgSetFeedReward{
		FeedId: feedId,
		FeedReward: &types.FeedRewardSchema{
			Amount:   400,
			Strategy: "",
		},
		Signer: bob.Addr,
	}

	s.Require().NoError(setSubmissionCountTx.ValidateBasic())
	s.Require().NoError(setHeartbeatTriggerTx.ValidateBasic())
	s.Require().NoError(setDeviationThresholdTriggerTx.ValidateBasic())
	s.Require().NoError(setFeedRewardTx.ValidateBasic())
	setFeedParamsResponse := s.BroadcastTx(ctx, bob, setSubmissionCountTx, setHeartbeatTriggerTx, setDeviationThresholdTriggerTx, setFeedRewardTx)
	s.Require().EqualValues(0, setFeedParamsResponse.TxResponse.Code)

	getFeedByFeedIdResponse, err = queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	s.Require().NoError(err)
	feed = getFeedByFeedIdResponse.GetFeed()
	s.Require().EqualValues(1, feed.GetSubmissionCount())
	s.Require().EqualValues(200, feed.GetHeartbeatTrigger())
	s.Require().EqualValues(300, feed.GetDeviationThresholdTrigger())
	s.Require().EqualValues(400, feed.GetFeedReward().GetAmount())
	s.Require().EqualValues("", feed.GetFeedReward().GetStrategy())

	s.T().Log("9 - Submit feed data by bob")


	// add bob in account store first before submitting feed data
	addBobInAccountStoreTx := &types.MsgAccount{
		Submitter:           bob.Addr,
		ChainlinkPublicKey:  []byte("bobChainlinkPublicKey"),
		ChainlinkSigningKey: []byte("ChainlinkSigningKey"),
		PiggyAddress:        bob.Addr,
	}
	s.Require().NoError(addBobInAccountStoreTx.ValidateBasic())
	addBobInAccountStoreTxResponse := s.BroadcastTx(ctx, bob, addBobInAccountStoreTx)
	s.Require().EqualValues(0, addBobInAccountStoreTxResponse.TxResponse.Code)


	submitFeedDataTx := &types.MsgFeedData{
		FeedId:                        feedId,
		ObservationFeedData:           [][]byte{[]byte("data")},
		ObservationFeedDataSignatures: [][]byte{[]byte("signature_bob")},
		Submitter:                     bob.Addr,
		CosmosPubKeys:                 [][]byte{[]byte(bob.Cosmos)},
	}
	s.Require().NoError(submitFeedDataTx.ValidateBasic())
	submitFeedDataResponse := s.BroadcastTx(ctx, bob, submitFeedDataTx)
	s.Require().EqualValues(0, submitFeedDataResponse.TxResponse.Code)

	getRoundDataResponse, err := queryClient.GetRoundData(ctx, &types.GetRoundDataRequest{FeedId: feedId, RoundId: 1})
	s.Require().NoError(err)
	roundData := getRoundDataResponse.GetRoundData()
	s.Require().EqualValues(1, len(roundData))

	// TODO check round data when OCR ready
}

func (s *GRPCTestSuite) waitForBlock(h int64) (int64, error) {
	ticker := time.NewTicker(time.Second)
	timeout := time.After(10 * time.Second)

	var latestHeight int64

	for {
		select {
		case <-timeout:
			ticker.Stop()
			return latestHeight, errors.New("timeout exceeded waiting for block")
		case <-ticker.C:
			status, err := s.rpcClient.Status(context.Background())
			if err == nil && status != nil {
				latestHeight = status.SyncInfo.LatestBlockHeight
				if latestHeight >= h {
					return latestHeight, nil
				}
			}
		}
	}
}
