package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"os"
)

func parseAuctionType(cdc codec.Codec, auctionType string) (string, error) {
	_, err := cdc.InterfaceRegistry().Resolve(auctionType)
	if err != nil {
		return "", fmt.Errorf("auction type %s not registered: %v", auctionType, err)
	}
	return auctionType, nil
}

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
