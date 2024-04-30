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
	Resolver auctiontypes.AuctionResolver
}

// Todo: pass denom and authority as configs
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ak auctiontypes.AccountKeeper, bk auctiontypes.BankKeeper, denom string, logger log.Logger) Keeper {
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

func (k *Keeper) SetAuctionTypesResolver(resolver auctiontypes.AuctionResolver) {
	k.Resolver = resolver
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

		logger.Info(fmt.Sprintf("Processing-Active :: Number of  bids: %v", auction.HasBids()))

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

		err = k.ExpiredAuctions.Set(goCtx, exp)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Processing-Active :: Removed Auction ID: %d", exp))

		err = k.ActiveAuctions.Remove(goCtx, exp)
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

	err := k.ExpiredAuctions.Walk(goCtx, nil, func(auctionId uint64) (stop bool, err error) {
		auction, err := k.Auctions.Get(ctx, auctionId)
		if err != nil {
			return true, err
		}
		logger.Info("Auction details", "auctionId:", auctionId, "hasBids:", auction.HasBids())

		if auction.HasBids() {
			err = k.PendingAuctions.Set(goCtx, auctionId)
			if err != nil {
				logger.Error("Failed to set auction to pending", "auctionId", auctionId, "error", err)
				return true, err
			}
			logger.Info(fmt.Sprintf("Successfully moved Auction ID with bids to pending: %d", auctionId))

			err = k.ExpiredAuctions.Remove(goCtx, auctionId)
			if err != nil {
				logger.Error("Failed to remove auction from expired", "auctionId", auctionId, "error", err)
				return true, err
			}
		} else {
			err = k.ExpiredAuctions.Remove(goCtx, auctionId)
			if err != nil {
				return true, err
			}
			err = k.CancelledAuctions.Set(goCtx, auctionId)
			if err != nil {
				return true, err
			}
			logger.Info(fmt.Sprintf("Processing-Expired :: Moved Auction ID without bids to cancelled: %d", auctionId))
		}
		numExpired++
		return false, nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Processing-Expired :: Number of expired auctions processed: %d", numExpired))
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

func (k *Keeper) GetAllAuctions(ctx sdk.Context) ([]auctiontypes.Auction, error) {
	var auctions []auctiontypes.Auction

	err := k.Auctions.Walk(ctx, nil, func(id uint64, auction auctiontypes.Auction) (stop bool, err error) {
		auctions = append(auctions, auction)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (k *Keeper) Validate() error {
	if k.cdc == nil {
		return fmt.Errorf("Codec (cdc) is not initialized")
	}
	if k.addressCodec == nil {
		return fmt.Errorf("AddressCodec is not initialized")
	}
	if k.authority == "" {
		return fmt.Errorf("Authority is not initialized")
	}
	if k.ak == nil {
		return fmt.Errorf("AccountKeeper is not initialized")
	}
	if k.bk == nil {
		return fmt.Errorf("BankKeeper is not initialized")
	}
	if k.defaultDenom == "" {
		return fmt.Errorf("DefaultDenom is not initialized")
	}
	if k.logger == nil {
		return fmt.Errorf("Logger is not initialized")
	}
	return nil
}
