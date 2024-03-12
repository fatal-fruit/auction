package types

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"time"
)

var (
	_ Auction         = &ReserveAuction{}
	_ AuctionMetadata = &ReserveAuctionMetadata{}
	_ AuctionHandler  = &ReserveAuctionHandler{}
)

type ReserveAuctionHandler struct {
	es EscrowService
	bk BankKeeper
}

func NewReserveAuctionHandler(es EscrowService, bk BankKeeper) *ReserveAuctionHandler {
	return &ReserveAuctionHandler{
		bk: bk,
		es: es,
	}
}

func (ra *ReserveAuction) GetType() string {
	return ra.AuctionType
}

func (ra *ReserveAuction) GetAuctionMetadata() AuctionMetadata {
	return ra.GetMetadata()
}

func (ra *ReserveAuction) SetOwner(owner sdk.AccAddress) {
	ra.Owner = owner.String()
}

func (ra *ReserveAuction) StartAuction(blockTime time.Time) {
	end := blockTime.Add(ra.Metadata.Duration)

	// Set Start and End time for auction
	ra.Metadata.StartTime = blockTime
	ra.Metadata.EndTime = end
}

// TODO: Implement safer logic to advance status
func (ra *ReserveAuction) UpdateStatus(newStatus string) {
	ra.Status = newStatus
}

// TODO: Implement logic to transfer funds
func (ra *ReserveAuction) SubmitBid(blockTime time.Time, bidMsg *MsgNewBid) error {
	// Validate bid price is over Reserve Price
	if bidMsg.Bid.IsLT(ra.Metadata.ReservePrice) {
		return fmt.Errorf("bid lower than reserve price :: %s", ra.Metadata.ReservePrice.String())
	}

	// Validate auction is active
	if blockTime.After(ra.Metadata.EndTime) {
		return fmt.Errorf("expired auction :: %s", ra.String())
	}

	// Validate bid price is competitive
	if len(ra.Metadata.Bids) > 0 && bidMsg.Bid.IsLTE(ra.Metadata.LastPrice) {
		return fmt.Errorf("bid lower than latest price :: %s", ra.Metadata.LastPrice)
	}

	ra.Metadata.Bids = append(ra.Metadata.Bids, &Bid{
		AuctionId: bidMsg.AuctionId,
		Bidder:    bidMsg.Owner,
		BidPrice:  bidMsg.Bid,
		Timestamp: blockTime,
	})

	ra.Metadata.LastPrice = bidMsg.Bid
	return nil
}

func (ra *ReserveAuction) IsExpired(blockTime time.Time) bool {
	return ra.Metadata.EndTime.Before(blockTime)
}

func (ra *ReserveAuction) HasBids() bool {
	return len(ra.Metadata.Bids) > 0
}

func (ah *ReserveAuctionHandler) CreateAuction(ctx context.Context, id uint64, am AuctionMetadata) (Auction, error) {
	md, ok := am.(proto.Message)
	if !ok {
		return &ReserveAuction{}, fmt.Errorf("%T does not implement proto.Message", md)
	}

	a := &ReserveAuction{
		Id:     id,
		Status: ACTIVE,
		Metadata: &ReserveAuctionMetadata{
			Bids: []*Bid{},
		},
	}

	switch m := am.(type) {
	case *ReserveAuctionMetadata:
		a.Metadata.Duration = m.Duration
		a.Metadata.ReservePrice = m.ReservePrice
	default:
		return &ReserveAuction{}, fmt.Errorf("invalid auction metadata :: %s", m.String())
	}

	strategy, err := BuildSettleStrategy(ctx, ah.es, id)
	if err != nil {
		return &ReserveAuction{}, fmt.Errorf("error creating escrow contract for auction id :: %d", id)
	}
	a.Metadata.Strategy = strategy.ToProto()

	return a, nil
}

func (ah *ReserveAuctionHandler) ExecAuction(ctx context.Context, auction Auction) error {
	switch a := auction.(type) {
	case *ReserveAuction:
		es := a.Metadata.GetStrategy()
		err := es.ExecuteStrategy(ctx, a, ah.es, ah.bk)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid auction metadata")
	}

	return nil
}

type Strategy interface {
	ExecuteStrategy(context.Context, ReserveAuction, EscrowService, BankKeeper) error
}

func (s *SettleStrategy) ExecuteStrategy(ctx context.Context, auction *ReserveAuction, es EscrowService, bk BankKeeper) error {
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

func GetWinner(auction *ReserveAuction) (*Bid, error) {
	var highestBid *Bid
	for _, b := range auction.Metadata.Bids {
		if highestBid.GetBidPrice().IsNil() || b.GetBidPrice().IsGTE(highestBid.GetBidPrice()) {
			highestBid = b
		}
	}

	return highestBid, nil
}
