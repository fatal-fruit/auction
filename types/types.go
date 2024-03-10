package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

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
}

type Auction interface {
	proto.Message

	GetId() uint64
	SetOwner(owner sdk.AccAddress)
	SubmitBid()
}

type AuctionMetadata interface {
	proto.Message
}

func (m *MsgNewAuction) SetMetadata(metadata AuctionMetadata) error {
	md, err := types.NewAnyWithValue(metadata)
	if err != nil {
		return err
	}
	m.AuctionMetadata = md
	return nil
}

// GetContent returns the proposal Content
