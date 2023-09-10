package types

import "cosmossdk.io/collections"

const (
	ModuleName = "auction"
	StoreKey   = "auctionKey"
)

var (
	NamesKey  = collections.NewPrefix(0)
	OwnersKey = collections.NewPrefix(1)
)
