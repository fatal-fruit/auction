package types

import "cosmossdk.io/collections"

const (
	ModuleName = "auction"
	StoreKey   = "auction"
	RESERVE    = "RESERVE"
)

var (
	IDKey                = collections.NewPrefix(0)
	AuctionsKey          = collections.NewPrefix(1)
	OwnerAuctionsKey     = collections.NewPrefix(2)
	ActiveAuctionsKey    = collections.NewPrefix(3)
	ExpiredAuctionsKey   = collections.NewPrefix(4)
	CancelledAuctionsKey = collections.NewPrefix(5)
	PendingAuctionsKey   = collections.NewPrefix(6)
)
