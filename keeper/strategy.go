package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

type Strategy interface {
	ExecuteStrategy(auctiontypes.EscrowService) error
}

type SettleStrategy struct {
	*auctiontypes.SettleStrategy
}

func (s *SettleStrategy) ExecuteStrategy(es auctiontypes.EscrowService) error {
	addr := sdk.MustAccAddressFromBech32(s.EscrowContractAddress)
	err := es.Release(s.EscrowContractId, addr)
	if err != nil {
		return err
	}
	return nil
}

func (s *SettleStrategy) ToProto() *auctiontypes.SettleStrategy {
	return &auctiontypes.SettleStrategy{
		StrategyType:          s.GetStrategyType(),
		EscrowContractId:      s.GetEscrowContractId(),
		EscrowContractAddress: s.GetEscrowContractAddress(),
	}
}

// Use generics
func BuildSettleStrategy(ctx context.Context, es auctiontypes.EscrowService) (*SettleStrategy, error) {
	contract, err := es.NewContract(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("Unable to create escrow contract")
	}
	s := &auctiontypes.SettleStrategy{
		StrategyType:          auctiontypes.SETTLE,
		EscrowContractId:      contract.GetId(),
		EscrowContractAddress: contract.GetAddress().String(),
	}
	return &SettleStrategy{s}, nil

}
