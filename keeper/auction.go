package keeper

import auctiontypes "github.com/fatal-fruit/auction/types"

func GetWinner(auction auctiontypes.ReserveAuction) (*auctiontypes.Bid, error) {
	var highestBid *auctiontypes.Bid
	for _, b := range auction.Bids {
		if highestBid.GetBidPrice().IsNil() || b.GetBidPrice().IsGTE(highestBid.GetBidPrice()) {
			highestBid = b
		}
	}

	return highestBid, nil
}
