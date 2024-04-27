package types

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type Auction interface {
	proto.Message

	GetId() uint64
	GetType() string
	GetAuctionMetadata() AuctionMetadata
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

type BidMetadata interface {
	proto.Message
}

type AuctionResolver interface {
	AddType(key string, h AuctionHandler) (rsv AuctionResolver)
	HasType(key string) bool
	GetHandler(key string) (h AuctionHandler)
	Seal()
	ListTypes() []string
}

type auctionResolver struct {
	handlers map[string]AuctionHandler
	sealed   bool
}

// NewResolver creates a new Auction Resolver interface instance
func NewResolver() AuctionResolver {
	return &auctionResolver{
		handlers: make(map[string]AuctionHandler),
	}
}

type AuctionHandler interface {
	CreateAuction(ctx context.Context, id uint64, metadata AuctionMetadata) (Auction, error)
	SubmitBid(ctx context.Context, auction Auction, bidMsg *MsgNewBid) (Auction, error)
	ExecAuction(ctx context.Context, a Auction) error
}

// Seal seals the resolver which prohibits any additionsl auction types to be
// registered. Seal panics if called more than once.
func (ar *auctionResolver) Seal() {
	if ar.sealed {
		panic("resolver already sealed")
	}
	ar.sealed = true
}

// AddType adds an auction type and its handler. It returns the Auction Type Resolver
// so AddType calls can be chained so long as it has not already been sealed.
func (ar *auctionResolver) AddType(key string, h AuctionHandler) AuctionResolver {
	if ar.sealed {
		panic("resolver sealed; cannot add auction type handler")
	}

	if ar.HasType(key) {
		panic(fmt.Sprintf("auction type %s has already been initialized", key))
	}

	ar.handlers[key] = h
	return ar
}

// HasType returns true if the auction type handler has been registered.
func (ar *auctionResolver) HasType(key string) bool {
	return ar.handlers[key] != nil
}

// GetHandler returns the auction type handler for a given key.
func (ar *auctionResolver) GetHandler(key string) AuctionHandler {
	if !ar.HasType(key) {
		panic(fmt.Sprintf("auction type handler \"%s\" does not exist", key))
	}

	return ar.handlers[key]
}

func (m *MsgNewAuction) SetMetadata(metadata AuctionMetadata) error {
	md, err := types.NewAnyWithValue(metadata)
	if err != nil {
		return err
	}
	m.AuctionMetadata = md
	return nil
}

// ListTypes returns all registered auction types.
func (ar *auctionResolver) ListTypes() []string {
	var types []string
	for key := range ar.handlers {
		types = append(types, key)
	}
	return types
}
