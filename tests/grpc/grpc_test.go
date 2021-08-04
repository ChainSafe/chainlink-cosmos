// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package grpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ChainSafe/chainlink-cosmos/app"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	xauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type testAccount struct {
	Priv   cryptotypes.PrivKey
	Addr   sdk.AccAddress
	Cosmos string
}

var (
	acc1 = generateKeyPair()
	acc2 = generateKeyPair()
	acc3 = generateKeyPair()
)

func initClientContext(t testing.TB) context.Context {
	userHomeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	home := filepath.Join(userHomeDir, ".chainlinkd")

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithHomeDir(home).
		WithViper("").
		WithAccountRetriever(xauthtypes.AccountRetriever{}).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

	clientCtx, err = config.ReadFromClientConfig(clientCtx)
	require.NoError(t, err)

	keyring, err := client.NewKeyringFromBackend(clientCtx, "test")
	require.NoError(t, err)

	clientCtx = clientCtx.WithKeyring(keyring)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	return ctx
}

func generateKeyPair() *testAccount {
	priv, pub, addr := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pub)
	if err != nil {
		panic(err)
	}
	return &testAccount{priv, addr, cosmosPubKey}
}

func setupGRPCConn(t testing.TB) *grpc.ClientConn {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Error(err)
	}

	return grpcConn
}

func Sign(t testing.TB, ctx context.Context, name string, txBuilder client.TxBuilder) error {
	clientCtx := ctx.Value(client.ClientContextKey).(*client.Context)

	key, err := clientCtx.Keyring.Key(name)
	require.NoError(t, err)
	pubKey := key.GetPubKey()

	accNum, accSeq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(*clientCtx, key.GetAddress())
	require.NoError(t, err)

	encCfg := simapp.MakeTestEncodingConfig()

	signerData := xauthsigning.SignerData{
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
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: accSeq,
	}

	err = txBuilder.SetSignatures(sig)
	require.NoError(t, err)

	bytesToSign, err := encCfg.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	require.NoError(t, err)

	sigBytes, _, err := clientCtx.Keyring.Sign(name, bytesToSign)
	require.NoError(t, err)

	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: accSeq,
	}

	return txBuilder.SetSignatures(sig)
}

func BroadcastTx(t testing.TB, ctx context.Context, grpcConn *grpc.ClientConn, submitter string, msgs ...sdk.Msg) *tx.BroadcastTxResponse {
	//clientCtx := ctx.Value(client.ClientContextKey).(*client.Context)
	txClient := tx.NewServiceClient(grpcConn)

	encCfg := simapp.MakeTestEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msgs...)
	require.NoError(t, err)

	txBuilder.SetGasLimit(testdata.NewTestGasLimit())
	//txBuilder.SetFeeAmount(...)
	//txBuilder.SetMemo(...)
	//txBuilder.SetTimeoutHeight(...)

	require.NoError(t, Sign(t, ctx, submitter, txBuilder))

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err)

	res, err := txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	require.NoError(t, err)

	return res
}

func TestGRPC(t *testing.T) {
	ctx := initClientContext(t)
	grpcConn := setupGRPCConn(t)
	queryClient := types.NewQueryClient(grpcConn)

	clientCtx := ctx.Value(client.ClientContextKey).(*client.Context)

	alice, err := clientCtx.Keyring.Key("alice")
	require.NoError(t, err)
	aliceCosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, alice.GetPubKey())
	require.NoError(t, err)

	bob, err := clientCtx.Keyring.Key("bob")
	require.NoError(t, err)
	bobCosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, bob.GetPubKey())
	require.NoError(t, err)

	cerlo, err := clientCtx.Keyring.Key("cerlo")
	require.NoError(t, err)
	cerloCosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, cerlo.GetPubKey())
	require.NoError(t, err)

	// 1 - Check initial module owner
	getModuleOwnerResponse, err := queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	require.NoError(t, err)
	moduleOwner := getModuleOwnerResponse.GetModuleOwner()
	require.Equal(t, 1, len(moduleOwner))
	require.Equal(t, alice.GetAddress(), moduleOwner[0].GetAddress())
	require.Equal(t, aliceCosmosPubKey, string(moduleOwner[0].GetPubKey()))

	// 2 - Add new module owner by alice
	addModuleOwnerTx := &types.MsgModuleOwner{
		Address:         bob.GetAddress(),
		PubKey:          []byte(bobCosmosPubKey),
		AssignerAddress: alice.GetAddress(),
	}
	require.NoError(t, addModuleOwnerTx.ValidateBasic())
	addModuleOwnerTxResponse := BroadcastTx(t, ctx, grpcConn, alice.GetName(), addModuleOwnerTx)
	require.EqualValues(t, 0, addModuleOwnerTxResponse.TxResponse.Code)

	time.Sleep(5 * time.Second)

	getModuleOwnerResponse, err = queryClient.GetAllModuleOwner(ctx, &types.GetModuleOwnerRequest{})
	require.NoError(t, err)
	moduleOwner = getModuleOwnerResponse.GetModuleOwner()
	require.Equal(t, 2, len(moduleOwner))

	// 3 - Module ownership transfer by bob to alice
	moduleOwnershipTransferTx := &types.MsgModuleOwnershipTransfer{
		NewModuleOwnerAddress: alice.GetAddress(),
		NewModuleOwnerPubKey:  []byte(aliceCosmosPubKey),
		AssignerAddress:       bob.GetAddress(),
	}
	require.NoError(t, addModuleOwnerTx.ValidateBasic())
	moduleOwnershipTransferTxResponse := BroadcastTx(t, ctx, grpcConn, bob.GetName(), moduleOwnershipTransferTx)
	require.EqualValues(t, 0, moduleOwnershipTransferTxResponse.TxResponse.Code)

	// 4 - Add new feed by alice
	feedId := "testfeed1"
	addFeedTx := &types.MsgFeed{
		FeedId:    feedId,
		FeedOwner: cerlo.GetAddress(),
		DataProviders: []*types.DataProvider{
			{
				Address: cerlo.GetAddress(),
				PubKey:  []byte(cerloCosmosPubKey),
			},
		},
		SubmissionCount:           1,
		HeartbeatTrigger:          2,
		DeviationThresholdTrigger: 3,
		FeedReward:                4,
		ModuleOwnerAddress:        alice.GetAddress(),
	}
	require.NoError(t, addFeedTx.ValidateBasic())
	addFeedTxResponse := BroadcastTx(t, ctx, grpcConn, alice.GetName(), addFeedTx)
	require.EqualValues(t, 0, addFeedTxResponse.TxResponse.Code)

	time.Sleep(5 * time.Second)

	getFeedByFeedIdResponse, err := queryClient.GetFeedByFeedId(ctx, &types.GetFeedByIdRequest{FeedId: feedId})
	require.NoError(t, err)
	feed := getFeedByFeedIdResponse.GetFeed()

	require.Equal(t, feedId, feed.GetFeedId())
	require.EqualValues(t, 1, feed.GetSubmissionCount())
	require.EqualValues(t, 2, feed.GetHeartbeatTrigger())
	require.EqualValues(t, 3, feed.GetDeviationThresholdTrigger())
	require.EqualValues(t, 4, feed.GetFeedReward())
}
