package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"time"
)

type Auction interface {
	proto.Message

	GetId() uint64
	GetType() string
	SetOwner(owner sdk.AccAddress)
	UpdateStatus(string)
	StartAuction(blockTime time.Time)
	SubmitBid(blockTime time.Time, bidMsg *MsgNewBid) error
	IsExpired(blockTime time.Time) bool
	HasBids() bool
}

type AuctionMetadata interface {
	proto.Message
}

type AuctionResolver interface {
	AddType(key string, h AuctionHandler) (rsv AuctionResolver)
	HasType(key string) bool
	GetHandler(key string) (h AuctionHandler)
	Seal()
}

type auctionResolver struct {
	handlers map[string]AuctionHandler
	sealed   bool
}

type AuctionHandler interface {
	CreateAuction(ctx context.Context, id uint64, metadata AuctionMetadata) (Auction, error)
	ExecAuction(ctx context.Context, a Auction) error
}

func (m *MsgNewAuction) SetMetadata(metadata AuctionMetadata) error {
	md, err := types.NewAnyWithValue(metadata)
	if err != nil {
		return err
	}
	m.AuctionMetadata = md
	return nil
}
