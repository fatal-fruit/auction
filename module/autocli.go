package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	auctionv1 "github.com/fatal-fruit/auction/api/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: auctionv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Auction",
					Use:       "auction [auction_id]",
					Short:     "Get auction by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "AllAuctions",
					Use:       "all-auctions",
					Short:     "Get all auctions",
				},
				{
					RpcMethod: "OwnerAuctions",
					Use:       "owner-auctions [owner-address]",
					Short:     "Query auctions by owner address",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:           auctionv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{},
		},
	}
}
