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

The Auctions module keeps state on all Auctions and corresponding Bids.

Before an auction has closed, it is also persisted in one of four queues:
- Active Auctions
- Expired Auctions
- Pending Auctions
- Cancelled Auctions

When an auction gets executed, it is removed from the Pending queue updated a final time in the Auctions table.

### Storage

Tables
- Auctions :: Map <UUID, Auction>

Queues
- Active :: Map <UUID, Auction>
- Pending :: Map <UUID, Auction>
- Expired :: Map <UUID, Auction>
- Cancelled :: Map <UUID, Auction>

Indexes
- AuctionsByOwner
- AuctionByBidderAddress
- ActiveAuctions
- BidsByAddress // TBD can filter on status (can access auction via bid)

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