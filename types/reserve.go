package types

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

var (
	_ Auction         = &ReserveAuction{}
	_ AuctionMetadata = &ReserveAuctionMetadata{}
)

func (ra *ReserveAuction) GetType() string {
	return ra.AuctionType
}

func (ra *ReserveAuction) GetAuctionMetadata() AuctionMetadata {
	return ra.GetMetadata()
}

func (ra *ReserveAuction) HasBids() bool {
	return len(ra.Metadata.Bids) > 0
}

func (ra *ReserveAuction) IsExpired(blockTime time.Time) bool {
	return ra.Metadata.EndTime.Before(blockTime)
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

// TODO: Implement logic to transfer funds
func (ra *ReserveAuction) SubmitBid(blockTime time.Time, bidMsg *MsgNewBid) error {
	// Validate bid price is over Reserve Price
	if bidMsg.BidAmount.IsLT(ra.Metadata.ReservePrice) {
		return fmt.Errorf("bid lower than reserve price :: %s", ra.Metadata.ReservePrice.String())
	}

	// Validate auction is active
	if blockTime.After(ra.Metadata.EndTime) {
		return fmt.Errorf("expired auction :: %s", ra.String())
	}

	// Validate bid price is competitive
	if len(ra.Metadata.Bids) > 0 && bidMsg.BidAmount.IsLTE(ra.Metadata.LastPrice) {
		return fmt.Errorf("bid lower than latest price :: %s", ra.Metadata.LastPrice)
	}

	ra.Metadata.Bids = append(ra.Metadata.Bids, &Bid{
		AuctionId: bidMsg.AuctionId,
		Bidder:    bidMsg.Owner,
		BidPrice:  bidMsg.BidAmount,
		Timestamp: blockTime,
	})

	ra.Metadata.LastPrice = bidMsg.BidAmount
	return nil
}

// TODO: Implement safer logic to advance status
func (ra *ReserveAuction) UpdateStatus(newStatus string) {
	ra.Status = newStatus
}

type Strategy interface {
	UpdateBid(ctx context.Context, bid Bid, service EscrowService, keeper BankKeeper) error
	ExecuteStrategy(context.Context, ReserveAuction, EscrowService, BankKeeper) error
}

func NewSettleStrategy(ctx context.Context, es EscrowService, id uint64) (*SettleStrategy, error) {
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

func (s *SettleStrategy) ExecuteStrategy(ctx context.Context, auction *ReserveAuction, es EscrowService, bk BankKeeper) error {
	// Select Winner
	winningBid, err := s.GetWinner(auction)
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

func (s *SettleStrategy) SubmitBid(ctx context.Context, bid *MsgNewBid, bk BankKeeper) error {
	// Handle funds
	bidder := sdk.MustAccAddressFromBech32(bid.GetOwner())
	amt := bid.GetBidAmount()
	escrowAddr := sdk.MustAccAddressFromBech32(s.EscrowContractAddress)

	// Send bid amount to escrow account
	return bk.SendCoins(ctx, bidder, escrowAddr, sdk.Coins{amt})
}

func (s *SettleStrategy) GetWinner(auction *ReserveAuction) (*Bid, error) {
	var highestBid *Bid
	for _, b := range auction.Metadata.Bids {
		if highestBid.GetBidPrice().IsNil() || b.GetBidPrice().IsGTE(highestBid.GetBidPrice()) {
			highestBid = b
		}
	}

	return highestBid, nil
}

func (s *SettleStrategy) ToProto() *SettleStrategy {
	return &SettleStrategy{
		StrategyType:          s.GetStrategyType(),
		EscrowContractId:      s.GetEscrowContractId(),
		EscrowContractAddress: s.GetEscrowContractAddress(),
	}
}
