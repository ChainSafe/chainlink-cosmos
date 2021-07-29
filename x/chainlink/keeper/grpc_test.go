// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

	//_ = os.Setenv("KEYRING_BACKEND", "test")

	clientCtx := client.Context{}.
		WithHomeDir(home).
		WithViper("")

	//_ = clientCtx.Viper.BindEnv("KEYRING_BACKEND")

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

	encCfg := simapp.MakeTestEncodingConfig()

	signerData := xauthsigning.SignerData{
		ChainID:       "testchain",
		AccountNumber: 1,
		Sequence:      1,
	}

	signMode := encCfg.TxConfig.SignModeHandler().DefaultMode()
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: 1,
	}

	var prevSignatures []signing.SignatureV2

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
		Sequence: 1,
	}

	prevSignatures = append(prevSignatures, sig)
	return txBuilder.SetSignatures(prevSignatures...)
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

func TestEmptyFeedRound(t *testing.T) {
	ctx := context.Background()
	grpcConn := setupGRPCConn(t)
	queryClient := types.NewQueryClient(grpcConn)

	res, err := queryClient.GetRoundData(ctx, &types.GetRoundDataRequest{
		FeedId:  "non-existing-feed",
		RoundId: 1,
	})
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, 0, len(res.GetRoundData()))
}

func TestFeedRound(t *testing.T) {
	ctx := initClientContext(t)
	grpcConn := setupGRPCConn(t)

	submitter := "alice"

	clientCtx := ctx.Value(client.ClientContextKey).(*client.Context)
	key, err := clientCtx.Keyring.Key(submitter)
	require.NoError(t, err)

	addFeedTx := &types.MsgFeed{
		FeedId:    "testfeed1",
		FeedOwner: acc2.Addr,
		DataProviders: []*types.DataProvider{
			{
				Address: acc3.Addr,
				PubKey:  []byte(acc3.Cosmos),
			},
		},
		SubmissionCount:           1,
		HeartbeatTrigger:          2,
		DeviationThresholdTrigger: 3,
		ModuleOwnerAddress:        key.GetAddress(),
		FeedReward:                4,
	}
	require.NoError(t, addFeedTx.ValidateBasic())

	res := BroadcastTx(t, ctx, grpcConn, submitter, addFeedTx)

	fmt.Printf("%+v\n", res.TxResponse)

	require.Equal(t, 0, res.TxResponse.Code)

	//res, err := queryClient.GetRoundData(ctx, &types.GetRoundDataRequest{
	//	FeedId:  "feedid1",
	//	RoundId: 1,
	//})
	//if err != nil {
	//	t.Error(err)
	//}

	//require.Equal(t, 0, len(res.GetRoundData()))
}
