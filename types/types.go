package types

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"time"
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

// NewRouter creates a new Router interface instance
func NewResolver() AuctionResolver {
	return &auctionResolver{
		handlers: make(map[string]AuctionHandler),
	}
}

type AuctionHandler interface {
	CreateAuction(ctx context.Context, id uint64, metadata AuctionMetadata) (Auction, error)
	ExecAuction(ctx context.Context, a Auction) error
}

// Seal seals the resolver which prohibits any additionsl route handlers to be
// added. Seal panics if called more than once.
func (ar *auctionResolver) Seal() {
	if ar.sealed {
		panic("resolver already sealed")
	}
	ar.sealed = true
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
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
