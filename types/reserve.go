package types

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

var _ Auction = &ReserveAuction{}
var _ AuctionMetadata = &ReserveAuctionMetadata{}
var _ AuctionHandler = &ReserveAuctionHandler{}

type ReserveAuctionHandler struct {
	es EscrowService
}

func (ra *ReserveAuction) SetOwner(owner sdk.AccAddress) {
	ra.Owner = owner.String()
}

func (ra *ReserveAuction) SubmitBid() {

}

func (ah *ReserveAuctionHandler) CreateAuction(ctx context.Context, id uint64, am AuctionMetadata) (Auction, error) {
	// TODO: derive concrete type from m
	md, ok := am.(proto.Message)
	if !ok {
		return &ReserveAuction{}, fmt.Errorf("%T does not implement proto.Message", md)
	}

	a := &ReserveAuction{
		Id:     id,
		Status: ACTIVE,
		Bids:   []*Bid{},
	}

	switch m := md.(type) {
	case *ReserveAuctionMetadata:
		a.Duration = m.Duration
		a.ReservePrice = m.ReservePrice
		a.AuctionType = m.AuctionType
	default:
		return &ReserveAuction{}, fmt.Errorf("invalid auction metadata", m)

	}

	strategy, err := BuildSettleStrategy(ctx, ah.es, id)
	if err != nil {
		return &ReserveAuction{}, fmt.Errorf("error creating escrow contract for auction")
	}
	a.Strategy = strategy.ToProto()

	return a, nil
}

type Strategy interface {
	ExecuteStrategy(context.Context, ReserveAuction, EscrowService, BankKeeper) error
}

func (s *SettleStrategy) ExecuteStrategy(ctx context.Context, auction ReserveAuction, es EscrowService, bk BankKeeper) error {
	// Select Winner
	winningBid, err := GetWinner(auction)
	if err != nil {
		return err
	}

	bidder := sdk.MustAccAddressFromBech32(winningBid.Bidder)
	auctioneer := sdk.MustAccAddressFromBech32(auction.Owner)

	// Send bid amount to auction owner
	err = bk.SendCoins(ctx, bidder, auctioneer, sdk.Coins{winningBid.BidPrice})
	if err != nil {
		return err
	}

	// Release escrowed auction bounty
	escrowAddr := sdk.MustAccAddressFromBech32(s.EscrowContractAddress)
	err = es.Release(ctx, s.EscrowContractId, escrowAddr, bidder)
	if err != nil {
		return err
	}
	return nil
}

func (s *SettleStrategy) ToProto() *SettleStrategy {
	return &SettleStrategy{
		StrategyType:          s.GetStrategyType(),
		EscrowContractId:      s.GetEscrowContractId(),
		EscrowContractAddress: s.GetEscrowContractAddress(),
	}
}

func BuildSettleStrategy(ctx context.Context, es EscrowService, id uint64) (*SettleStrategy, error) {
	contract, err := es.NewContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Unable to create escrow contract")
	}
	s := &SettleStrategy{
		StrategyType:          SETTLE,
		EscrowContractId:      contract.GetId(),
		EscrowContractAddress: contract.GetAddress().String(),
	}
	return s, nil
}

func GetWinner(auction ReserveAuction) (*Bid, error) {
	var highestBid *Bid
	for _, b := range auction.Bids {
		if highestBid.GetBidPrice().IsNil() || b.GetBidPrice().IsGTE(highestBid.GetBidPrice()) {
			highestBid = b
		}
	}

	return highestBid, nil
}
