// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package ocr

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr/signature"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeVerifyOCR(t *testing.T) {

	////////////////////// ENCODE //////////////////////

	var configDigest [16]byte
	_, err := rand.Read(configDigest[:])
	require.NoError(t, err)

	reportContext := &types.ReportContext{
		ConfigDigest: &types.ConfigDigest{Value: configDigest[:]},
		Epoch:        uint32(time.Now().Unix()),
		Round:        1,
	}

	fmt.Printf("ConfigDigest: %+v\n", reportContext)

	key1Id := ocrtypes.OracleID(42) // example OracleID
	key1, err := signature.NewKeyBundle()
	require.NoError(t, err)

	key2Id := ocrtypes.OracleID(88) // example OracleID
	key2, err := signature.NewKeyBundle()
	require.NoError(t, err)

	observation1 := &types.Observation{Value: big.NewInt(100).Bytes()}
	observation2 := &types.Observation{Value: big.NewInt(101).Bytes()}

	observations := types.AttributedObservations{
		{
			Observation: observation1,
			Observer:    uint32(key1Id),
		},
		{
			Observation: observation2,
			Observer:    uint32(key2Id),
		},
	}
	report1, err := types.MakeAttestedReportOne(observations, reportContext, key1.SignOnChain)
	require.NoError(t, err)

	report2, err := types.MakeAttestedReportOne(observations, reportContext, key2.SignOnChain)
	require.NoError(t, err)

	reportFinal := &types.AttestedReportMany{
		AttributedObservations: observations,
		Signatures: [][]byte{
			report1.Signature,
			report2.Signature,
		},
	}

	result, err := Pack(reportContext, reportFinal)
	require.NoError(t, err)

	fmt.Println("ABI encoded result: " + hexutil.Encode(result))

	////////////////////// DECODE //////////////////////

	offchainReport, err := Unpack(result)
	require.NoError(t, err)

	fmt.Printf("%+v\n", offchainReport)

	for i, o := range offchainReport.Report.AttributedObservations {
		t.Run(fmt.Sprintf("observation:%d,observer:%d", observations[i].Observation.GoEthereumValue(), observations[i].Observer), func(t *testing.T) {
			require.EqualValues(t, o.Observation.GoEthereumValue(), observations[i].Observation.GoEthereumValue())
			require.EqualValues(t, o.Observer, observations[i].Observer)
		})
	}

	////////////////////// VERIFY //////////////////////

	whitelist := signature.Addresses{}
	whitelist[key1.PublicKeyAddressOnChain()] = key1Id
	whitelist[key2.PublicKeyAddressOnChain()] = key2Id

	err = offchainReport.GetReport().VerifySignatures(offchainReport.GetContext(), whitelist)
	require.NoError(t, err)
}
