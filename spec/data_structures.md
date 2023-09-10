# Data Structures & Methods

### Interfaces
```go
type AuctionStatus int

const (
    ACTIVE AuctionStatus = iota + 1
    PENDING
    CANCELLED
    SETTLED
)

type AuctionType string

type ExecutionStrategy interface {
    // TODO: Define specific execution strategies 
    Execute(a Auction)    bool, error
}

// TODO: Change Name
type SimpleSettleStrategy type {}

type Auction interface {
    Initialize()         bool, error
    UpdateStatus()       bool, error // Update AuctionStatus
    PlaceBid()           bool, error
    AmmendBid()          bool, error
    Execute()            bool, error
}

```

### Concrete Types
```go
type Bid struct {
    Id          UUID
    Bidder      sdk.AccAddress
    Amount      []sdk.Coin
    AuctionId   UUID
}

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

```
### Methods

```go
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

func (ss *SimpleSettleStrategy) Execute(ra ReserveAuction) (bool, error){
    // Send Deposit to ra.HighestBidder
    // Send ra.HighestBidder.Amount to Auctioneer
    // Return all other bids to bidders
    // Return success
}
```