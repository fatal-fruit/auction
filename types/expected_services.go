package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type EscrowService interface {
	NewContract() (uint64, error)
	Release(address sdk.AccAddress) error
}
