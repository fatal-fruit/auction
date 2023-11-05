package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/fatal-fruit/auction/types"
)

// Escrow Implementation
type EscrowModule struct {
	ak types.AccountKeeper
}

type EscrowModContract struct {
	Id      uint64
	Address sdk.AccAddress
}

func (em *EscrowModule) NewContract(ctx context.Context, id uint64) (types.EscrowContract, error) {
	// Generate account address of contract
	var accountAddr sdk.AccAddress
	for {
		nextAccVal := id
		derivationKey := make([]byte, 8)
		binary.BigEndian.PutUint64(derivationKey, nextAccVal)

		ac, err := authtypes.NewModuleCredential(types.ModuleName, types.ContractAddressPrefix, derivationKey)
		if err != nil {
			return nil, fmt.Errorf("could not create contract account :: %w", err)
		}
		accountAddr = sdk.AccAddress(ac.Address())
		if em.ak.GetAccount(ctx, accountAddr) != nil {
			continue
		}

		account, err := authtypes.NewBaseAccountWithPubKey(ac)
		if err != nil {
			return nil, fmt.Errorf("could not create contract account :: %w", err)
		}

		_ = em.ak.NewAccount(ctx, account)

		break
	}

	return &EscrowModContract{
		Id:      id,
		Address: accountAddr,
	}, nil
}

func (em *EscrowModule) Release(address sdk.AccAddress) error {
	// TODO: Implement
	return nil
}

func (ec *EscrowModContract) GetId() uint64 {
	return ec.Id
}

func (ec *EscrowModContract) GetAddress() sdk.AccAddress {
	return ec.Address
}
