package keeper

import "github.com/fatal-fruit/auction/types"

// Escrow Implementation
type EscrowModule struct {
}

func NewEscrowModule() types.EscrowService {
	return &EscrowModule{}
}

func (em *EscrowModule) NewContract() (uint64, error) {
	return 1, nil
}
