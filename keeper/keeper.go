package keeper

import (
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
}

func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ak auctiontypes.AccountKeeper, bk auctiontypes.BankKeeper, es auctiontypes.EscrowService, denom string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	ids := collections.NewSequence(sb, auctiontypes.IDKey, "auctionIds")
	auctions := collections.NewMap(sb, auctiontypes.AuctionsKey, "auctions", collections.Uint64Key, codec.CollValue[auctiontypes.ReserveAuction](cdc))
	ownerAuctions := collections.NewMap(sb, auctiontypes.OwnerAuctionsKey, "ownerAuctions", sdk.AccAddressKey, codec.CollValue[auctiontypes.OwnerAuctions](cdc))

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

	return k
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetDefaultDenom() string {
	return k.defaultDenom
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+auctiontypes.ModuleName)
}
