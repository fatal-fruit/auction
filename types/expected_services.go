package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EscrowService interface {
	NewContract(context.Context, uint64) (EscrowContract, error)
	Release(uint64, sdk.AccAddress) error
}

type EscrowContract interface {
	GetId() uint64
	GetAddress() sdk.AccAddress
}
