# `x/auction`

Canonical implementation of general purpose permisionless auctions for the Cosmos SDK.

## Concepts

`Reserve Auction`
- Bid Protocol

`Dutch Auction`
- Bid Protocol

`Execution Strategies`
- Simple Settle

## State

### Interface
```go
type AuctionStatus int

const (
    ACTIVE AuctionStatus = iota + 1
    PENDING
    CANCELLED
    SETTLED
)

type AuctionType string

type Bid struct {
    Id          UUID
    Bidder      sdk.AccAddress
    Amount      []sdk.Coin
    AuctionId   UUID
}

type ExecutionStrategy interface {
    // TODO: Define specific execution strategies 
    Execute(a Auction)    bool, error
}

// TODO: Change Name
type SimpleSettleStrategy type {}

func (ss *SimpleSettleStrategy) Execute(ra ReserveAuction) (bool, error){
    // Send Deposit to ra.HighestBidder
    // Send ra.HighestBidder.Amount to Auctioneer
    // Return all other bids to bidders
    // Return success
}

type Auction interface {
    Initialize()         bool, error
    UpdateStatus()       bool, error // Update AuctionStatus
    PlaceBid()           bool, error
    AmmendBid()          bool, error
    Execute()            bool, error
}

// concrete
type ReserveAuction struct {
    Id                      UUID
    Auctioneer              sdk.AccAddress
    Duration                time.Duration
    EscrowAcc               sdk.AccAddress
    Bids                    []Bids
    HighestBid              Bid
    Status                  AuctionStatus
    Strategy                ExecutionStrategy
	
    Extended                bool
    ExtensionDuration       time.Duration 
    // Starting Price
    ReservePrice            []sdk.Coin
}

type DutchAuction struct {
    Id                      UUID
    Auctioneer              sdk.AccAddress
    Duration                time.Duration
    EscrowAcc               sdk.AccAddress
    Bids                    []Bids
    BestBid                 Bid
    Status                  AuctionStatus
    Strategy                ExecutionStrategy

    //Initialized as start price
    CurrentPrice            []sdk.Coin //Initialized as start price
    // ReservePrice            []sdk.Coin
    // Cannot save reserve price on chain without threshold decryption might be able to leverage VE
    // Might be interesting to include automatic price adjustment at a cadence in strategy
}

func (da *DutchAuction) PlaceBid(currTime time.Duration){
	// TODO
}

func (ra *ReserveAuction) PlaceBid(currTime time.Duration){
    // If currTime is past ra.Duration, then return error
    //      Else check if bid is higher than ra.HighestBid
    //      If bid is highest effective bid, add to ra.Bids and update ra.HighestBid
    // Else reject bid
	
	// Check for extension (if already extended, ignore)
	// If currTime is within ra.ExtensionDuration, update ra.Duration and set ra.Extended to true
}

func (ra *ReserveAuction) AmmendBid(currTime time.Duration){
    // If bid is not in ra.Bids, return error	
    // If currTime is past ra.Duration, then return error
    // If bid is higher than ra.HighestBid
    //      If bid is highest effective bid, update ra.Bids and update ra.HighestBid
    // Else reject update
	
    //  Check for extension logic as in PlaceBid
}

func (da *DutchAuction) AmmendBid(currTime time.Duration){
    // TODO
}

func (da *DutchAuction) Execute(){
    da.Strategy.Execute(da)
}

func (ra *ReserveAuction) Execute(){
    ra.Strategy.Execute(ra)
}


```
### Storage

Tables
Auctions :: Map <UUID, Auction>

Queues
Active :: Map <UUID, Auction>
Pending :: Map <UUID, Auction>
Expired :: Map <UUID, Auction>
Cancelled :: Map <UUID, Auction>

Indexes
AuctionsByOwner
AuctionByBidderAddress
ActiveAuctions
BidsByAddress // TBD can filter on status (can access auction via bid)

## State Transitions

### Auctions
- New Auction
  - Initializes a new Auction in the Auction Table, and add to the `Active` auction queue
- Canceled Auction
  - Only applicable to auctions without bids
  - Pop auction from `Active` queue and push to `Cancelled` queue
- Updated Auction
  - New Bid
    - Only applicable if auction is in `Active` queue
    - Will call `PlaceBid()` on auction which will determine the bid's acceptance
  - Ammended Bid
    - Only applicable if auction is in `Active` queue
    - Bid must represent best effective price for specific `AuctionType`
- Extended Auction
  - Bid
    - For certain auction types like `ReserveAuction`, if a bid is submitted within the `ExtensionDuration`, the auction `Duration` will be extended by the preset `DurationAmount`
- Execute Auction
  - Every auction has a custom execution strategy that specifies how to settle assets betweent the Auctioneer and winner

### Endblock
- Process Auction Queues
  - Active 
    - Iterate through all `Active` auctions that have officially elapsed their `Duration` and push to the `Expired` queue. 
  - Expired
    - Iterate through all `Expired` auctions.
    - Auctions with at least 1 bid are pushed to `Pending`.
    - Auctions with no bids are pushed to `Cancelled`.
  - Pending
    - Iterate through all `Pending` auctions.
    - Auctions that have not been executed and have elapsed the chain's `auction_expire_time` parameter are pushed to `Cancelled`
  - Cancelled
    - Iterate through all `Cancelled` auctions.
    - Return deposited assets to Auctioneer address
    - Retrieve all bids and return amounts.
    - Remove auction from queue and update state in `Auction` table.

## Invariants
- Auction may only be cancelled while no bids have been placed
- If a bid is placed within the `ExtensionDuration` period, the auction is extended the same amount of time as the `ExtensionDuration` // TODO
- Bids cannot be placed after `Duration` has elapsed
- Bid amount must be above the auction's reserve price // TODO: Verify if in CheckTx
- `Duration` cannot be ammended

## Messages

### CreateAuction
```protobuf

message MsgCreateAuctionMessage {
  option (cosmos.msg.v1.signer) = "auctioneer";
  
  
}

```

### Bid
```protobuf

message MsgNewBid {
  option (cosmos.msg.v1.signer) = "bidder";
  
}

message MsgAmmendBid {
  option (cosmos.msg.v1.signer) = "bidder";

}
```

### CancelAuction
```protobuf

message MsgCancelAuctionMessage {
  option (cosmos.msg.v1.signer) = "auctioneer";
  
  
}

```

### ExecuteAuction
```protobuf

message MsgExecuteAuctionMessage {
  option (cosmos.msg.v1.signer) = "creator";
  
  
}

```

## End Block
See [Endblock](#endblock)

## Hooks

Describe available hooks to be called by/from this module.

## Events

List and describe event tags used.

## Client

List and describe CLI commands and gRPC and REST endpoints.

## Params

**Auction Expire Time**

Expiration duration for pending auctions to be executed. If an auction in `PENDING` state exceeds this time, it will be cancelled and all bids refunded.
```json
{
  "auction_expire_time": "1209600" // 2 weeks
}
```

## Future Improvements

Describe future improvements of this module.

## Tests

Acceptance tests.

## Appendix

Supplementary details referenced elsewhere within the spec.