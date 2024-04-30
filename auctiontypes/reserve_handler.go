package auctiontypes

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/fatal-fruit/auction/types"
)

var _ types.AuctionHandler = &ReserveAuctionHandler{}

type ReserveAuctionHandler struct {
	es types.EscrowService
	bk types.BankKeeper
}

func NewReserveAuctionHandler(es types.EscrowService, bk types.BankKeeper) *ReserveAuctionHandler {
	return &ReserveAuctionHandler{
		bk: bk,
		es: es,
	}
}

func (ah *ReserveAuctionHandler) CreateAuction(ctx context.Context, id uint64, am types.AuctionMetadata) (types.Auction, error) {

	// if _, ok := am.(types.Auction); !ok {
	//     return nil, fmt.Errorf("provided data does not implement fatal_fruit.auction.v1.Auction interface")
	// }

	md, ok := am.(proto.Message)
	if !ok {
		return &ReserveAuction{}, fmt.Errorf("%T does not implement proto.Message", md)
	}

	a := &ReserveAuction{
		Id:          id,
		Status:      types.ACTIVE,
		AuctionType: sdk.MsgTypeURL(&ReserveAuction{}),
		Metadata: &ReserveAuctionMetadata{
			Bids: []*types.Bid{},
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

func (ah *ReserveAuctionHandler) SubmitBid(ctx context.Context, auction types.Auction, bidMsg *types.MsgNewBid) (types.Auction, error) {
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

	log.Printf("Checking if any bids have been recorded: auctionId=%d", auction.GetId())

	if len(auction.GetAuctionMetadata().(*ReserveAuctionMetadata).Bids) == 0 {
		return nil, fmt.Errorf("no bids recorded after submission")
	}

	return auction, nil
}

func (ah *ReserveAuctionHandler) ExecAuction(ctx context.Context, auction types.Auction) error {
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
