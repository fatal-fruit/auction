# `x/auction`

Canonical implementation of general purpose permisionless auctions for the Cosmos SDK.

## Concepts

Describe specialized concepts and definitions used throughout the spec.

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

type Bid interface {
	Id          UUID
	Bidder      sdk.AccAddress
	Amount      []sdk.Coin
	AuctionId   UUID
}

type ExecutionStrategy interface {
	// TODO: Define specific execution strategies
}

type Auction struct {
	Id                      UUID
	Auctioneer              sdk.AccAddress
    Duration                time.Duration     
	EscrowAcc               sdk.AccAddress
	Bids                    [] Bids
	HighestBid              Bid
	Status                  AuctionStatus
	Strategy                ExecutionStrategy
}

// concrete
type ReserveAuction struct {
    ExtensionDuration       time.Duration
	ReservePrice            []sdk.Coin
}

type DutchAuction struct {
    StartPrice            []sdk.Coin
	// Cannot save reserve price on chain without threshold decryption might be able to leverage VE
	// Might be interesting to include automatic price adjustment at a cadence in strategy
}


```
### Storage

Tables
Auctions :: Map <UUID, Auction>
Bids :: Map <UUID, Bids>

Indexes
AuctionsByOwner
AuctionByBidderAddress
ActiveAuctions
BidsByAddress // TBD can filter on status (can access auction via bid)

## State Transitions

```protobuf
service AuctionService {
  rpc CreateAuction(MsgCreateAuction) returns (MsgCreateAuctionResponse);
  
  rpc NewBid(MsgNewBid) returns (MsgNewBidResponse);
}
```


### Auctions
State transitions for auctions occur in `EndBlock` 

#### Canceled Auction

#### Extended Auction

#### Ended Auction

#### Pending Auction


### Bids
When a bid is placed or updated, it affects the auction state, and can trigger a state transition if the bid is placed after the auction has expired but before the `EndBlocker`. 

#### Bid
If the application attempts to place a bid with a submission timestamp after the auctions's duration, the transaction will be rejected, and the Auction will be transferred to the `PENDING` or `EXPIRED` queue depending if any bids have been placed.

For certain auctions (Reserve), if the bid is placed within the `ExtensionDuration`, it will trigger an update to the auction's `Duration`. 

#### Update Bid
Similarly to a first bid, if an existing bid is updated with a new price, it will trigger the same state transition as a regular bid if it 

### Module
Endlblock protocol

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
  option (cosmos.msg.v1.signer) = "creator";
  
  
}

```

### Bid
```protobuf

message MsgNewBid {
  option (cosmos.msg.v1.signer) = "bidder";
  
}

```
#### State Modifications
- Generate new `AuctionId`
- TBD

### Bid

## End Block
- At `EndBlock`, all active auctions that have elapsed their duration are set to `PENDING`. 
- Any `PENDING` auctions that have since expired in the last block are added to `CANCELLED` queue for processing.
- Process all auctions in `CANCELLED` queue to return bids.

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