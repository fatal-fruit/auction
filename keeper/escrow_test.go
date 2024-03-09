package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewContract(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	f.MockAcctKeeper.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	f.MockAcctKeeper.EXPECT().NewAccount(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	escrowMod := keeper.NewEscrowModule(f.MockAcctKeeper, f.MockBankKeeper)

	var previousAddress types.AccAddress
	for i := 0; i < 10; i++ {
		id, err := f.K.IDs.Next(f.Ctx)
		require.NoError(err)

		contract, err := escrowMod.NewContract(f.Ctx, id)
		require.NoError(err)

		address := contract.GetAddress()
		require.NotEmpty(t, address)
		require.NotEqual(t, address, f.ModAddr)
		if i > 0 {
			require.NotEqual(t, address, previousAddress)
		}
		previousAddress = address
	}
}
