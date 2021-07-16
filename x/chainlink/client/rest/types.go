// SPDX-License-Identifier: MIT

package rest

import "github.com/cosmos/cosmos-sdk/types/rest"

type createFeedRequest struct {
	BaseReq    rest.BaseReq `json:"baseReq"`
	FeedId     string       `json:"feedId"`
	FeedData   []byte       `json:"feedData"`
	Signatures [][]byte     `json:"signature"`
}
