package types

import "cosmossdk.io/collections"

const (
	ModuleName = "auction"
	StoreKey   = "auction"

	// Strategy Types
	SETTLE = "SETTLE"

	// Auction Status
	ACTIVE = "ACTIVE"
	CLOSED = "CLOSED"
)

var (
	IDKey                 = collections.NewPrefix(0)
	AuctionsKey           = collections.NewPrefix(1)
	OwnerAuctionsKey      = collections.NewPrefix(2)
	ActiveAuctionsKey     = collections.NewPrefix(3)
	ExpiredAuctionsKey    = collections.NewPrefix(4)
	CancelledAuctionsKey  = collections.NewPrefix(5)
	PendingAuctionsKey    = collections.NewPrefix(6)
	ContractAddressPrefix = collections.NewPrefix(7)
)
