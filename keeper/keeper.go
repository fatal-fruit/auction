package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec
	logger       log.Logger

	authority string

	ak           auctiontypes.AccountKeeper
	bk           auctiontypes.BankKeeper
	defaultDenom string

	// state management
	Schema        collections.Schema
	IDs           collections.Sequence
	Auctions      collections.Map[uint64, auctiontypes.Auction]
	OwnerAuctions collections.Map[sdk.AccAddress, auctiontypes.OwnerAuctions]

	// Queues
	ActiveAuctions    collections.KeySet[uint64]
	ExpiredAuctions   collections.KeySet[uint64]
	PendingAuctions   collections.KeySet[uint64]
	CancelledAuctions collections.KeySet[uint64]

	// Auction Type Registry
	resolver auctiontypes.AuctionResolver
}

// Todo: remove escrow, pass auction type resolver
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ak auctiontypes.AccountKeeper, bk auctiontypes.BankKeeper, resolver auctiontypes.AuctionResolver, denom string, logger log.Logger) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	ids := collections.NewSequence(sb, auctiontypes.IDKey, "auctionIds")
	auctions := collections.NewMap(sb, auctiontypes.AuctionsKey, "auctions", collections.Uint64Key, codec.CollInterfaceValue[auctiontypes.Auction](cdc))
	ownerAuctions := collections.NewMap(sb, auctiontypes.OwnerAuctionsKey, "ownerAuctions", sdk.AccAddressKey, codec.CollValue[auctiontypes.OwnerAuctions](cdc))
	activeAuctions := collections.NewKeySet(sb, auctiontypes.ActiveAuctionsKey, "activeAuctions", collections.Uint64Key)
	expiredAuctions := collections.NewKeySet(sb, auctiontypes.ExpiredAuctionsKey, "expiredAuctions", collections.Uint64Key)
	cancelledAuctions := collections.NewKeySet(sb, auctiontypes.CancelledAuctionsKey, "cancelledAuctions", collections.Uint64Key)
	pendingAuctions := collections.NewKeySet(sb, auctiontypes.PendingAuctionsKey, "pendingAuctions", collections.Uint64Key)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		ak:           ak,
		bk:           bk,
		resolver:     resolver,
		defaultDenom: denom,
		logger:       logger,
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
	k.PendingAuctions = pendingAuctions

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
func (keeper *Keeper) Logger() log.Logger {
	return keeper.logger.With("module", "x/"+auctiontypes.ModuleName)
}

func (k *Keeper) ProcessActiveAuctions(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()
	var expired []uint64

	var numActive int
	logger.Info("Processing-Active :: Checking for active auctions")
	err := k.ActiveAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		auction, err := k.Auctions.Get(ctx, auctionId)
		if err != nil {
			return true, err
		}
		// TODO: Auction checks itself for expiration
		if auction.IsExpired(ctx.BlockTime()) {
			expired = append(expired, auctionId)
		} else {
			numActive++
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Processing-Active :: Number of active auctions: %d", numActive))
	for _, exp := range expired {
		err = k.ActiveAuctions.Remove(goCtx, exp)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Active :: Removed Auction ID: %d", exp))

		err = k.ExpiredAuctions.Set(goCtx, exp)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Active :: Pushed to Expired: %d", exp))
	}
	return nil
}

func (k *Keeper) ProcessExpiredAuctions(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()
	logger.Info("Processing-Expired :: Checking for expired auctions")
	var numExpired int

	var pending []uint64
	var cancelled []uint64
	err := k.ExpiredAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		auction, err := k.Auctions.Get(ctx, auctionId)
		if err != nil {
			return true, err
		}

		// TODO: Auction executes own logic for this
		if auction.HasBids() {
			pending = append(pending, auctionId)
		} else {
			cancelled = append(cancelled, auctionId)
		}
		numExpired++
		return false, nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Processing-Expired :: Number of expired auctions: %d", numExpired))

	// If no bids -> cancelled
	logger.Info("Processing-Expired :: Checking for cancelled auctions")
	for _, c := range cancelled {
		err = k.ExpiredAuctions.Remove(goCtx, c)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Expired :: Removed Auction ID without bids from expired: %d", c))

		err = k.CancelledAuctions.Set(goCtx, c)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Expired :: Pushed Auction ID without bids to cancelled: %d", c))

	}
	// If at least 1 bid -> pending
	logger.Info("Processing-Expired :: Checking for pending auctions")
	for _, p := range pending {
		err = k.ExpiredAuctions.Remove(goCtx, p)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Expired :: Removed Auction ID with bids from expired: %d", p))

		err = k.PendingAuctions.Set(goCtx, p)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Expired :: Pushed Auction ID with bids to pending: %d", p))

	}
	return nil
}

func (k *Keeper) GetPending(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()
	logger.Info("Processing-Pending :: Checking for pending auctions")
	var numPending int

	var pending []uint64
	err := k.PendingAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		if err != nil {
			return true, err
		}
		pending = append(pending, auctionId)
		numPending++
		return false, nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Processing-Pending :: Number of pending auctions: %d", numPending))
	logger.Info(fmt.Sprintf("Processing-Pending :: Pending Auctions: %d", pending))
	return nil
}

func (k *Keeper) PurgeCancelledAuctions(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.CancelledAuctions.Walk(sdkCtx, nil, func(auctionId uint64) (stop bool, err error) {
		err = k.CancelledAuctions.Remove(ctx, auctionId)
		if err != nil {
			return true, err
		}
		return false, nil
	})

	if err != nil {
		k.Logger().Error("Failed to purge cancelled auctions", "error", err)
		return err
	}

	k.Logger().Info("All cancelled auctions have been purged")
	return nil
}

func (k *Keeper) GetCancelledAuctions(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()
	logger.Info("Processing-Cancelled :: Checking for cancelled auctions")
	var numCancelled int

	var cancelled []uint64
	err := k.CancelledAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		if err != nil {
			return true, err
		}
		cancelled = append(cancelled, auctionId)
		numCancelled++
		return false, nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Processing-Cancelled :: Number of cancelled auctions: %d", numCancelled))
	logger.Info(fmt.Sprintf("Processing-Cancelled :: Cancelled Auctions: %v", cancelled))
	return nil
}

// CancelAuction marks an auction as cancelled by its ID.
func (k *Keeper) CancelAuction(ctx context.Context, auctionId uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_, err := k.Auctions.Get(sdkCtx, auctionId)
	if err != nil {
		return fmt.Errorf("auction with ID %d not found: %v", auctionId, err)
	}

	err = k.CancelledAuctions.Set(sdkCtx, auctionId)
	if err != nil {
		return fmt.Errorf("failed to cancel auction with ID %d: %v", auctionId, err)
	}

	k.Logger().Info("Auction cancelled", "auctionId", auctionId)
	return nil
}

func (k *Keeper) GetAllAuctions(ctx sdk.Context) []auctiontypes.ReserveAuction {
	//var auctions []auctiontypes.ReserveAuction
	//
	//err := k.Auctions.Walk(ctx, nil, func(id uint64, auction auctiontypes.ReserveAuction) (stop bool, err error) {
	//	auctions = append(auctions, auction)
	//	return false, nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//return auctions
	return nil
}
