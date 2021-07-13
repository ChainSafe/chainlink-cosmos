package ante

import (
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
	ErrFeedDoesNotExist = "feed does not exist"
)

func NewAnteHandler(
	ak authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, chainLinkKeeper chainlinkkeeper.Keeper,
	sigGasConsumer authante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
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
		case *types.MsgAddDataProvider:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			if (types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetDataProvider().GetAddress()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data provider already registered")
			}
		case *types.MsgRemoveDataProvider:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
			if !(types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetAddress()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data provider not present")
			}
		case *types.MsgSetSubmissionCount:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
		case *types.MsgSetHeartbeatTrigger:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
		case *types.MsgSetDeviationThresholdTrigger:
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, ErrFeedDoesNotExist)
			}
		default:
			continue
		}
	}

	// TODO check feed owner #12
	// https://github.com/ChainSafe/chainlink-cosmos/issues/12

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
			feed := fd.chainLinkKeeper.GetFeed(ctx, t.GetFeedId())
			if feed.Feed.Empty() {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "feed not exist")
			}
			if !(types.DataProviders)(feed.GetFeed().GetDataProviders()).Contains(t.GetSubmitter()) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "invalid data provider")
			}
		default:
			continue
		}
	}

	return next(ctx, tx, simulate)
}
