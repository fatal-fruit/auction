package testutil

import (
	"context"
	// "crypto/rand"
	"encoding/binary"
	"fmt"

	// "math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/fatal-fruit/auction/types"
)

/*
Escrow Service Test Implementation

**Not intended for use in production**

BasicEscrowService is a testing util for the purpose of facilitating simple escrow contracts.
The ReserveAuction implementation specifies a strategy to settle the auction. This strategy is instantiated
with an escrow contract before the auction is started. When the auction is executed, it releases the contract
and sends its balance to the winner of the auction.

EscrowService specifies two functions, NewContract, and Release. NewContract is called when the Auction
implementation calls CreateAuction. Release is called when the Auction calls Exec.

EscrowModContract consists of a single Id and Module Account address.
*/
type BasicEscrowService struct {
	ak types.AccountKeeper
	bk types.BankKeeper
}

type EscrowModContract struct {
	Id      uint64
	Address sdk.AccAddress
}

func NewTestEscrowModule(ak types.AccountKeeper, bk types.BankKeeper) types.EscrowService {
	return &BasicEscrowService{
		ak: ak,
		bk: bk,
	}
}

func (em *BasicEscrowService) NewContract(ctx context.Context, id uint64) (types.EscrowContract, error) {
	// Generate account address of contract
	accountAddr, err := em.generateUniqueAccountAddress(ctx, em.ak, id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique account address: %w", err)
	}

	return &EscrowModContract{
		Id:      id,
		Address: accountAddr,
	}, nil
}

func (em *BasicEscrowService) Release(ctx context.Context, id uint64, sender sdk.AccAddress, recipient sdk.AccAddress) error {
	//TODO: Refactor to get contract address

	//TODO: Extend for all balances
	balances := em.bk.GetAllBalances(ctx, sender)

	err := em.bk.SendCoins(ctx, sender, recipient, balances)
	if err != nil {
		return err
	}

	return nil
}

func (ec *EscrowModContract) GetId() uint64 {
	return ec.Id
}

func (ec *EscrowModContract) GetAddress() sdk.AccAddress {
	return ec.Address
}

func (em *BasicEscrowService) generateUniqueAccountAddress(ctx context.Context, ak types.AccountKeeper, id uint64) (sdk.AccAddress, error) {
	var accountAddr sdk.AccAddress
	nextAccVal := id + 1
	derivationKey := make([]byte, 8)

	binary.BigEndian.PutUint64(derivationKey, nextAccVal)

	ac, err := authtypes.NewModuleCredential(types.ModuleName, types.ContractAddressPrefix, derivationKey)
	if err != nil {
		return nil, fmt.Errorf("could not create contract account :: %w", err)
	}
	accountAddr = sdk.AccAddress(ac.Address())

	if em.ak.GetAccount(ctx, accountAddr) != nil {
		return nil, fmt.Errorf("could not create get account :: %w", err)
	}

	account, err := authtypes.NewBaseAccountWithPubKey(ac)
	if err != nil {
		return nil, fmt.Errorf("could not create contract account :: %w", err)
	}

	_ = em.ak.NewAccount(ctx, account)

	return accountAddr, nil
}
