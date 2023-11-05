package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	authority string

	ak           auctiontypes.AccountKeeper
	bk           auctiontypes.BankKeeper
	es           auctiontypes.EscrowService
	defaultDenom string

	// state management
	Schema        collections.Schema
	IDs           collections.Sequence
	Auctions      collections.Map[uint64, auctiontypes.ReserveAuction]
	OwnerAuctions collections.Map[sdk.AccAddress, auctiontypes.OwnerAuctions]

	// Queues
	ActiveAuctions    collections.KeySet[uint64]
	ExpiredAuctions   collections.KeySet[uint64]
	PendingAuctions   collections.KeySet[uint64]
	CancelledAuctions collections.KeySet[uint64]
}

func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ak auctiontypes.AccountKeeper, bk auctiontypes.BankKeeper, es auctiontypes.EscrowService, denom string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	ids := collections.NewSequence(sb, auctiontypes.IDKey, "auctionIds")
	auctions := collections.NewMap(sb, auctiontypes.AuctionsKey, "auctions", collections.Uint64Key, codec.CollValue[auctiontypes.ReserveAuction](cdc))
	ownerAuctions := collections.NewMap(sb, auctiontypes.OwnerAuctionsKey, "ownerAuctions", sdk.AccAddressKey, codec.CollValue[auctiontypes.OwnerAuctions](cdc))
	activeAuctions := collections.NewKeySet(sb, auctiontypes.ActiveAuctionsKey, "activeAuctions", collections.Uint64Key)
	expiredAuctions := collections.NewKeySet(sb, auctiontypes.ExpiredAuctionsKey, "expiredAuctions", collections.Uint64Key)
	cancelledAuctions := collections.NewKeySet(sb, auctiontypes.CancelledAuctionsKey, "cancelledAuctions", collections.Uint64Key)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		ak:           ak,
		bk:           bk,
		es:           es,
		defaultDenom: denom,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema
	k.IDs = ids
	k.Auctions = auctions
	k.OwnerAuctions = ownerAuctions
	k.ActiveAuctions = activeAuctions
	k.ExpiredAuctions = expiredAuctions
	k.CancelledAuctions = cancelledAuctions

	return k
}

func (k *Keeper) GetAuthority() string {
	return k.authority
}

func (k *Keeper) GetDefaultDenom() string {
	return k.defaultDenom
}

func (k *Keeper) GetModuleBalance(ctx context.Context, denom string) sdk.Coin {
	moduleAddress := k.ak.GetModuleAddress(auctiontypes.ModuleName)
	modAcc := k.ak.GetAccount(ctx, moduleAddress)
	if modAcc == nil {
		return sdk.Coin{}
	}
	return k.bk.GetBalance(ctx, modAcc.GetAddress(), denom)
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+auctiontypes.ModuleName)
}

func (k *Keeper) ProcessActiveAuctions(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var expired []uint64
	err := k.ActiveAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		auction, err := k.Auctions.Get(ctx, auctionId)
		if err != nil {
			return true, err
		}
		if auction.EndTime.Before(ctx.BlockTime()) {
			expired = append(expired, auctionId)
		}
		return false, nil
	})
	if err != nil {
		panic(err)
	}
	for _, exp := range expired {
		err = k.ActiveAuctions.Remove(goCtx, exp)
		if err != nil {
			panic(err)
		}
		err = k.ExpiredAuctions.Set(goCtx, exp)
		if err != nil {
			panic(err)
		}
	}
}

func (k *Keeper) ProcessExpiredAuctions(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var pending []uint64
	var cancelled []uint64
	err := k.ExpiredAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		auction, err := k.Auctions.Get(ctx, auctionId)
		if err != nil {
			return true, err
		}
		if len(auction.Bids) > 0 {
			pending = append(pending, auctionId)
		} else {
			cancelled = append(cancelled, auctionId)
		}
		return false, nil
	})
	if err != nil {
		panic(err)
	}
	// If no bids -> cancelled
	for _, c := range cancelled {
		err = k.ExpiredAuctions.Remove(goCtx, c)
		if err != nil {
			panic(err)
		}
		err = k.CancelledAuctions.Set(goCtx, c)
		if err != nil {
			panic(err)
		}
	}
	// If at least 1 bid -> pending
	for _, p := range pending {
		err = k.ExpiredAuctions.Remove(goCtx, p)
		if err != nil {
			panic(err)
		}
		err = k.PendingAuctions.Set(goCtx, p)
		if err != nil {
			panic(err)
		}
	}
}
