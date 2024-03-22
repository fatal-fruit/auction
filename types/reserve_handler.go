package types

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var _ AuctionHandler = &ReserveAuctionHandler{}

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

func (ah *ReserveAuctionHandler) CreateAuction(ctx context.Context, id uint64, am AuctionMetadata) (Auction, error) {
	md, ok := am.(proto.Message)
	if !ok {
		return &ReserveAuction{}, fmt.Errorf("%T does not implement proto.Message", md)
	}

	a := &ReserveAuction{
		Id:          id,
		Status:      ACTIVE,
		AuctionType: sdk.MsgTypeURL(&ReserveAuction{}),
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

	strategy, err := NewSettleStrategy(ctx, ah.es, id)
	if err != nil {
		return &ReserveAuction{}, fmt.Errorf("error creating escrow contract for auction id :: %d", id)
	}
	a.Metadata.Strategy = strategy.ToProto()

	return a, nil
}

func (ah *ReserveAuctionHandler) SubmitBid(ctx context.Context, auction Auction, bidMsg *MsgNewBid) (Auction, error) {
	// Update auction with bid logic
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := auction.SubmitBid(sdkCtx.BlockTime(), bidMsg)
	if err != nil {
		return nil, fmt.Errorf("error submitting bid from auction handler")
	}

	am := auction.GetAuctionMetadata()
	switch resMd := am.(type) {
	case *ReserveAuctionMetadata:
		// Send bid amount to escrow contract
		s := resMd.GetStrategy()
		err = s.SubmitBid(ctx, bidMsg, ah.bk)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid auction metadata type")
	}

	return auction, nil
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
