// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package ante

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func feedRewardSchemaStrategyChecker(strategy string) error {
	if strategy != "" {
		_, ok := types.FeedRewardStrategyConvertor[strategy]
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "invalid feed reward strategy")
		}
	}

	return nil
}

// TODO: chainlink pubKey against observation signature :resp.GetAccount().GetChainlinkPublicKey() VS signature
// TODO: replace with validation logic here
func pubKeySignatureValidate(chainlinkPubKey, signature []byte) bool {
	return true
}

// TODO: observation VS observationSignature
func signaturePlainDataValidate(signature, data []byte) bool {
	return true
}
