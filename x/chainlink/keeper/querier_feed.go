package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func listFeed(ctx sdk.Context, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	msgs := keeper.GetAllFeed(ctx)

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, msgs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
