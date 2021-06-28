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

func NewAnteHandler(
	ak authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, chainLinkKeeper chainlinkkeeper.Keeper,
	sigGasConsumer authante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		anteHandler = sdk.ChainAnteDecorators(
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
			NewModuleOwnerDecorator(chainLinkKeeper),
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
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Msg: empty Msg", tx)
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
	for _, msg := range tx.GetMsgs() {
		t, ok := msg.(*types.ModuleOwner)
		if !ok {
			continue
		}

		if len(t.GetSigners()) == 0 {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Tx: empty signer", t)
		}
		txSigner := t.GetSigners()[0]

		signers = append(signers, txSigner)
	}

	for _, signer := range signers {
		if !(types.ModuleOwners)(existingModuleOwnerList.GetModuleOwner()).Contains(signer) {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "account %s (%s) is not a module owner", common.BytesToAddress(signer.Bytes()), signer)
		}
	}

	return next(ctx, tx, simulate)
}
