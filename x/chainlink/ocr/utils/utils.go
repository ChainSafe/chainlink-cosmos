package utils

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr/signature"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type OCRAccount struct {
	OracleID ocrtypes.OracleID
	Key      *signature.KeyBundle
}

func GenerateFakeReport(roundId uint64, obs []int64) (*types.ReportContext, *types.AttestedReportMany, []*OCRAccount, error) {
	var configDigest [16]byte
	_, err := rand.Read(configDigest[:])
	if err != nil {
		return nil, nil, nil, err
	}

	reportContext := &types.ReportContext{
		ConfigDigest: &types.ConfigDigest{Value: configDigest[:]},
		Epoch:        uint32(time.Now().Unix()),
		Round:        uint8(roundId),
	}

	numberOfReport := len(obs)

	var accounts []*OCRAccount
	for i := 0; i < numberOfReport; i++ {
		keyId := ocrtypes.OracleID(i + 1)
		key, err := signature.NewKeyBundle()
		if err != nil {
			return nil, nil, nil, err
		}
		accounts = append(accounts, &OCRAccount{
			OracleID: keyId,
			Key:      key,
		})
	}

	observations := types.AttributedObservations{}
	for i := 0; i < numberOfReport; i++ {
		observations = append(observations, &types.AttributedObservation{
			Observation: &types.Observation{Value: big.NewInt(obs[i]).Bytes()},
			Observer:    uint32(accounts[i].OracleID),
		})
	}

	var signatures [][]byte
	for i := 0; i < numberOfReport; i++ {
		report1, err := types.MakeAttestedReportOne(observations, reportContext, accounts[i].Key.SignOnChain)
		if err != nil {
			return nil, nil, nil, err
		}

		signatures = append(signatures, report1.Signature)
	}

	reportFinal := &types.AttestedReportMany{
		AttributedObservations: observations,
		Signatures:             signatures,
	}

	return reportContext, reportFinal, accounts, nil
}

func MustGenerateFakeABIReport(roundId uint64, obs []int64) []byte {
	context, report, _, err := GenerateFakeReport(roundId, obs)
	if err != nil {
		panic(err)
	}

	result, err := ocr.Pack(context, report)
	if err != nil {
		panic(err)
	}

	return result
}
