package types

import "cosmossdk.io/collections"

const (
	ModuleName = "auction"
	StoreKey   = "auction"
	RESERVE    = "RESERVE"
)

var (
	NamesKey  = collections.NewPrefix(0)
	OwnersKey = collections.NewPrefix(1)
)

var (
	IDKey            = collections.NewPrefix(0)
	AuctionsKey      = collections.NewPrefix(1)
	OwnerAuctionsKey = collections.NewPrefix(2)
)
