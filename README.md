# Auction Module

`x/auction` is a general purpose Cosmos-SDK module that provides an agnostic framework for on-chain auctions.  

### API

A minimal API is exposed for creating, updating, and executing application specific auctions. The auction module expects all auction types to be registered via an auction handler for each specific implementation.
Auction 

For an in depth understanding of auction state mechanics, see the most recent [spec](/spec/spec_v2.md).

### Integrating

**App Wiring**
`app.go`

Ensure that the `AuctionKeeper` is instantiated after the `AccountKeeper` and `BankKeeper`
```go
import (
	auctiontypes "github.com/fatal-fruit/auction/types"
	auctionkeeper "github.com/fatal-fruit/auction/keeper"
	auction "github.com/fatal-fruit/auction/module"
)

type App struct {
    keys            map[string]*storetypes.KVStoreKey
    mm              *module.Manager
    configurator    module.Configurator
}

func NewApp() *App {
	
    // Add store keys
    keys := storetypes.NewKVStoreKeys(auctiontypes.StoreKey)
    app := &App{
        keys: keys,
    }

    // Initialize keeper
    auctionKeeper := auctionkeeper.NewKeeper(
        encConfig.Codec,
        addresscodec.NewBech32Codec("cosmos"),
        storeService,
        authority.String(),
        accountKeeper,
        bankKeeper,
        sdk.DefaultBondDenom,
        app.Logger(),
    )
	
    /*
        Configure Auction Handlers Here
     */
    
    // Configure module
    app.mm = module.NewManager(auction.NewAppModule(appCodec, auctionKeeper))
    
    // Configure endblockers
    app.mm.SetOrderEndBlockers(auctiontypes.ModuleName)
    app.SetEndBlocker(app.EndBlocker)	
}

```

**Set Auction Handlers**

```golang
import (
    auctiontypes "github.com/fatal-fruit/auction/types"
    auctionkeeper "github.com/fatal-fruit/auction/keeper"
    auction "github.com/fatal-fruit/auction/module"
)
func NewApp() *App {
	
    k := auctionkeeper.NewKeeper()

    /*
       Configure Auction Handlers Here
    */
	
    // Instantiate a new auction resolve
    resolver := auctiontypes.NewResolver()
	
    // Create a new auction type handler that implements AuctionHandler
    // The basic concrete type is the ReserveAuction handler
    handler := auctiontypes.NewReserveAuctionHandler(escrowService, bankService)
    
    // Set the auction type on the resolver. 
    // AddType() returns an instance of the resolver so calls can be chained. 
    resolver.AddType(sdk.MsgTypeURL(&auctiontypes.ReserveAuction{}), handler)
	
    // Seal and set the resolver on the auction keeper
    resolver.seal()
    k.SetAuctionTypesResolver(resolver)

}
```

### Acknowlegements
This work was made possible by funding from the [AADAO](https://www.atomaccelerator.com/).