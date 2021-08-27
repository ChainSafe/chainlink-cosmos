// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package ocr

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var transmitTypes = getTransmitTypes()

func Pack(context *types.ReportContext, report *types.AttestedReportMany) ([]byte, error) {
	serializedReport, rs, ss, vs, err := report.TransmissionArgs(context)
	if err != nil {
		return nil, err
	}

	result, err := transmitTypes.Pack(serializedReport, rs, ss, vs)
	if err != nil {
		return nil, err
	}

	return result, err
}

func Unpack(data []byte) (*types.OffChainReport, error) {
	rawArgs, err := transmitTypes.Unpack(data)
	if err != nil {
		return nil, err
	}

	// TODO check cast
	args := struct {
		SerializedReport []byte
		RS               [][32]byte
		SS               [][32]byte
		VS               [32]byte
	}{
		SerializedReport: rawArgs[0].([]byte),
		RS:               rawArgs[1].([][32]byte),
		SS:               rawArgs[2].([][32]byte),
		VS:               rawArgs[3].([32]byte),
	}

	rawReport, err := types.ReportTypes.Unpack(args.SerializedReport)
	if err != nil {
		return nil, err
	}

	// TODO check cast
	report := struct {
		RawReportContext [32]byte
		RawObservers     [32]byte
		Observations     []*big.Int
	}{
		RawReportContext: rawReport[0].([32]byte),
		RawObservers:     rawReport[1].([32]byte),
		Observations:     rawReport[2].([]*big.Int),
	}

	// RawObservers is list of Observer/Oracle ID

	// RawReportContext consists of:
	// 11-byte zero padding
	// 16-byte configDigest
	// 4-byte epoch
	// 1-byte round

	// TODO check report.RawReportContext length
	var reportContext types.ReportContext
	reportContext.ConfigDigest = &types.ConfigDigest{}
	reportContext.ConfigDigest.Value = report.RawReportContext[11 : 11+16]
	reportContext.Epoch = binary.BigEndian.Uint32(report.RawReportContext[11+16:])
	reportContext.Round = report.RawReportContext[11+16+4]

	var observerCount int
	for _, o := range report.RawObservers {
		if o != 0 { // TODO check if OracleID 0 is possible
			observerCount++
		}
	}
	if len(report.Observations) != observerCount {
		return nil, errors.New("the number of observations must be equal to number of observers")
	}

	var attributedObservations []*types.AttributedObservation
	for i, o := range report.Observations {
		attributedObservations = append(attributedObservations, &types.AttributedObservation{
			Observation: &types.Observation{Value: o.Bytes()},
			Observer:    uint32(report.RawObservers[i]),
		})
	}

	var signatures [][]byte
	for i := 0; i < len(args.RS); i++ {
		var sig []byte
		sig = append(sig, args.RS[i][:]...)
		sig = append(sig, args.SS[i][:]...)
		sig = append(sig, args.VS[i])
		signatures = append(signatures, sig)
	}

	attestedReportMany := &types.AttestedReportMany{
		AttributedObservations: attributedObservations,
		Signatures:             signatures,
	}

	offchainReport := &types.OffChainReport{
		Context: &reportContext,
		Report:  attestedReportMany,
	}

	return offchainReport, nil
}

func getTransmitTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return []abi.Argument{
		{Name: "_report", Type: mustNewType("bytes")},
		{Name: "_rs", Type: mustNewType("bytes32[]")},
		{Name: "_ss", Type: mustNewType("bytes32[]")},
		{Name: "_rawVs", Type: mustNewType("bytes32")},
	}
}
