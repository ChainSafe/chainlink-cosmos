// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package ante

import (
	"bytes"
	chainlinkkeeper "github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ErrFeedDoesNotExist         = "feed does not exist"
	ErrSignerIsNotFeedOwner     = "account %s (%s) is not a feed owner"
	ErrAccountAlreadyExists     = "there is already a chainlink account associated with this cosmos address"
	ErrUnregisteredDataProvider = "linked account not found in account store"
	ErrDoesNotExist             = "no chainlink account associated with this cosmos address"
	ErrSubmitterDoesNotMatch    = "submitter address does not match"
)

func NewAnteHandler(
	ak authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, chainLinkKeeper chainlinkkeeper.Keeper,
	sigGasConsumer authante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler, externalTxDataValidationFunc func(sdk.Msg) bool,
) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		anteHandler := sdk.ChainAnteDecorators(
			authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
			authante.NewRejectExtensionOptionsDecorator(),
			authante.NewMempoolFeeDecorator(),
			authante.NewValidateBasicDecorator(),
			authante.TxTimeoutHeightDecorator{},
			authante.NewValidateMemoDecorator(ak),
			authante.NewConsumeGasForTxSizeDecorator(ak),
			authante.NewRejectFeeGranterDecorator(),
			authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
			authante.NewValidateSigCountDecorator(ak),
			authante.NewDeductFeeDecorator(ak, bankKeeper),
			authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
			authante.NewSigVerificationDecorator(ak, signModeHandler),
			authante.NewIncrementSequenceDecorator(ak),
			// all customized anteHandler below
			NewModuleOwnerDecorator(chainLinkKeeper),
			NewFeedDecorator(chainLinkKeeper),
			NewFeedDataDecorator(chainLinkKeeper),
			NewValidationDecorator(externalTxDataValidationFunc),
			NewAccountDecorator(chainLinkKeeper),
		)

		return anteHandler(ctx, tx, sim)
	}
}

type ModuleOwnerDecorator struct {
	chainLinkKeeper chainlinkkeeper.Keeper
}

func NewModuleOwnerDecorator(chainLinkKeeper chainlinkkeeper.Keeper) ModuleOwnerDecorator {
	return ModuleOwnerDecorator{
		chainLinkKeeper: chainLinkKeeper,
	}
}

func (mod ModuleOwnerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) == 0 {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg: %T", tx)
	}

	existingModuleOwnerList, err := mod.chainLinkKeeper.GetAllModuleOwner(sdk.WrapSDKContext(ctx), nil)
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrLogic, "module owner check failed at anteHandler[ModuleOwnerDecorator]")
	}

	// no checking if module owner list is empty
	if len(existingModuleOwnerList.GetModuleOwner()) == 0 {
		return next(ctx, tx, simulate)
	}

	signers := make([]sdk.AccAddress, 0)

	// get the signers of module owner Msg types
	for _, msg := range tx.GetMsgs() {
		switch t := msg.(type) {
		case *types.MsgModuleOwner:
			if len(t.GetSigners()) == 0 {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Tx: empty signer: %T", t)
			}
			signers = append(signers, t.GetSigners()[0])
		case *types.MsgModuleOwnershipTransfer:
			if len(t.GetSigners()) == 0 {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Tx: empty signer: %T", t)
			}
			signers = append(signers, t.GetSigners()[0])
		case *types.MsgFeed:
			if len(t.GetSigners()) == 0 {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Tx: empty signer: %T", t)
			}
			signers = append(signers, t.GetSigners()[0])
		default:
			continue
		}
	}

	for _, signer := range signers {
		if !(types.MsgModuleOwners)(existingModuleOwnerList.GetModuleOwner()).Contains(signer) {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "account %s (%s) is not a module owner", common.BytesToAddress(signer.Bytes()), signer)
		}
	}

	return next(ctx, tx, simulate)
}

type FeedDecorator struct {
	chainLinkKeeper chainlinkkeeper.Keeper
}

func NewFeedDecorator(chainLinkKeeper chainlinkkeeper.Keeper) FeedDecorator {
	return FeedDecorator{
		chainLinkKeeper: chainLinkKeeper,
	}
}

func (fd FeedDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) == 0 {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg: %T", tx)
	}

	for _, msg := range tx.GetMsgs() {
		switch t := msg.(type) {
		case *types.MsgFeed:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if !feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "feed already exists")
			}
			// check reward schema strategy
			err := feedRewardSchemaStrategyChecker(t.GetFeedReward().GetStrategy())
			if err != nil {
				return ctx, err
			}
		case *types.MsgAddDataProvider:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			if (types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetDataProvider().GetAddress()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data provider already registered")
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgRemoveDataProvider:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			if !(types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetAddress()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data provider not present")
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgSetSubmissionCount:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgSetHeartbeatTrigger:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgSetDeviationThresholdTrigger:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgSetFeedReward:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
			// check reward schema strategy
			err := feedRewardSchemaStrategyChecker(t.GetFeedReward().GetStrategy())
			if err != nil {
				return ctx, err
			}
		case *types.MsgFeedOwnershipTransfer:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		case *types.MsgRequestNewRound:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			signer := t.GetSigners()[0]
			if !feed.GetFeed().GetFeedOwner().Equals(signer) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, ErrSignerIsNotFeedOwner, common.BytesToAddress(signer.Bytes()), signer)
			}
		default:
			continue
		}
	}

	return next(ctx, tx, simulate)
}

type FeedDataDecorator struct {
	chainLinkKeeper chainlinkkeeper.Keeper
}

func NewFeedDataDecorator(chainLinkKeeper chainlinkkeeper.Keeper) FeedDataDecorator {
	return FeedDataDecorator{
		chainLinkKeeper: chainLinkKeeper,
	}
}

func (fd FeedDataDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) == 0 {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg: %T", tx)
	}

	for _, msg := range tx.GetMsgs() {
		switch t := msg.(type) {
		case *types.MsgFeedData:
			// get tx fee
			feeTx, ok := tx.(sdk.FeeTx)
			if !ok {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
			}
			txFee := feeTx.GetFee()
			if len(txFee) == 0 {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "empty tx fee coin slices")
			}
			// get the first coin from txFee as the tx fee charged for MsgFeedData tx
			t.TxFee = &types.Coin{
				Denom:  txFee[0].Denom,
				Amount: txFee[0].Amount.Uint64(),
			}

			// get feed by feedId
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "feed not exist")
			}

			// basic checking
			if !(types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetSubmitter()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "submitter is not a valid data provider")
			}
			if uint32(len(t.GetObservationFeedDataSignatures())) < feed.GetFeed().GetSubmissionCount() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "not enough signatures")
			}

			// observation signatures VS original observation data
			if len(t.GetObservationFeedData()) != len(t.GetObservationFeedDataSignatures()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "number of observation signatures and observation data does not match")
			}
			feedDataValidationFlag := len(t.GetObservationFeedData())
			for _, observation := range t.GetObservationFeedData() {
				for _, observationSignature := range t.GetObservationFeedDataSignatures() {
					if signaturePlainDataValidate(observationSignature, observation) {
						feedDataValidationFlag--
						break
					}
				}
			}
			if feedDataValidationFlag != 0 {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "observation signatures validation against observation data failed")
			}

			for _, pubKey := range t.GetCosmosPubKeys() {
				cosmosAddr, err := types.DeriveCosmosAddrFromPubKey(string(pubKey))
				if err != nil {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "invalid data provider: invalid cosmos pubkey")
				}

				dataProviderAddr, err := sdk.AccAddressFromBech32(cosmosAddr.String())
				if err != nil {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "invalid data provider: failed to derive cosmos address, invalid cosmos pubkey")
				}

				// valid data provider checking
				if !(types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(dataProviderAddr) {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "invalid data provider: data provider not in the list")
				}

				resp := fd.chainLinkKeeper.GetAccount(ctx, &types.GetAccountRequest{AccountAddress: dataProviderAddr})
				if resp.GetAccount().GetSubmitter().String() == "" {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ErrUnregisteredDataProvider)
				}

				// chainlink pubKey VS observation signature validation
				chainlinkPubKeyValidationFlag := false
				for _, signature := range t.GetObservationFeedDataSignatures() {
					if pubKeySignatureValidate(resp.GetAccount().GetChainlinkPublicKey(), signature) {
						chainlinkPubKeyValidationFlag = true
						break
					}
				}
				if !chainlinkPubKeyValidationFlag {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "observation signature validation against chainlink pubKey failed")
				}
			}

		default:
			continue
		}
	}

	return next(ctx, tx, simulate)
}

type ValidationDecorator struct {
	validationFn func(sdk.Msg) bool
}

func NewValidationDecorator(validationFunc func(sdk.Msg) bool) ValidationDecorator {
	return ValidationDecorator{
		validationFn: validationFunc,
	}
}

func (fd ValidationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) == 0 {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg: %T", tx)
	}

	for _, msg := range tx.GetMsgs() {
		switch t := msg.(type) {
		case *types.MsgFeedData:
			t.IsFeedDataValid = t.Validate(fd.validationFn)
		default:
			continue
		}
	}

	return next(ctx, tx, simulate)
}

type AccountDecorator struct {
	chainLinkKeeper chainlinkkeeper.Keeper
}

func NewAccountDecorator(chainLinkKeeper chainlinkkeeper.Keeper) AccountDecorator {
	return AccountDecorator{
		chainLinkKeeper: chainLinkKeeper,
	}
}

func (fd AccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) == 0 {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg: %T", tx)
	}

	for _, msg := range tx.GetMsgs() {
		switch t := msg.(type) {
		case *types.MsgAccount:
			// case to add a new chainlink account to the Account Store
			req := &types.GetAccountRequest{AccountAddress: t.Submitter}
			resp := fd.chainLinkKeeper.GetAccount(ctx, req)
			if resp.Account.Submitter.String() != "" {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ErrAccountAlreadyExists)
			}
		// case to edit an existing chainlink account in the Account Store
		case *types.MsgEditAccount:
			req := &types.GetAccountRequest{AccountAddress: t.Submitter}
			// submitters must match
			resp := fd.chainLinkKeeper.GetAccount(ctx, req)
			if resp.Account.Submitter.String() == "" {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ErrDoesNotExist)
			}
			if !bytes.Equal(t.Submitter.Bytes(), resp.Account.Submitter.Bytes()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ErrSubmitterDoesNotMatch)
			}
		default:
			continue
		}
	}

	return next(ctx, tx, simulate)
}
