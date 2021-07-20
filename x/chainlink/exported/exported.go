package exported

import (
	"github.com/gogo/protobuf/proto"
)

// RoundDataI TODO
type RoundDataI interface {
	proto.Message

	GetFeedId() string
	GetFeedData() OCRAbiEncodedI
}

// ObservationI TODO
type ObservationI interface {
	proto.Message

	GetData() []byte
}

// OCRAbiEncodedI TODO
type OCRAbiEncodedI interface {
	proto.Message

	GetContext() []byte
	GetOracles() []byte
	GetObservations() []ObservationI
}
