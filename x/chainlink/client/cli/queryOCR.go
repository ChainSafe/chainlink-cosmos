// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr/utils"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func CmdGenerateFakeOCR() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-fake-ocr",
		Short: "Generate a fake OCR report with two observation (100, 101)",
		RunE: func(cmd *cobra.Command, args []string) error {
			context, report, _, err := utils.GenerateFakeReport(1, []int64{100, 101})
			if err != nil {
				return err
			}

			result, err := ocr.Pack(context, report)
			if err != nil {
				return err
			}

			fmt.Println(hexutil.Encode(result))

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
