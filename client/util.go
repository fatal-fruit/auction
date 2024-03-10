package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"os"
)

func parseAuctionMetadata(cdc codec.Codec, auctionMetadataFile string) (auctiontypes.AuctionMetadata, error) {
	if auctionMetadataFile == "" {
		return nil, fmt.Errorf("invalid or missing auction metadata")
	}

	contents, err := os.ReadFile(auctionMetadataFile)
	if err != nil {
		return nil, err
	}

	var am auctiontypes.AuctionMetadata
	if err := cdc.UnmarshalInterfaceJSON(contents, &am); err != nil {
		return nil, fmt.Errorf("failed to parse auction metadata: %w", err)
	}

	return am, nil
}
