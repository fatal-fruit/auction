package keeper_test

import (
	"testing"

	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewContract(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	id, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)

	escrowServ := f.MockEscrowService
	// escrowMod := keeper.NewEscrowModule(authkeeper.NewAccountKeeper(f.Ctx, authtypes.StoreKey, ), bankkeeper.NewBaseKeeper(f.Ctx))

	require.NoError(err)
	contract, err := escrowServ.NewContract(f.Ctx, id)
	require.NoError(err)

	require.Equal(t, contract.GetId(), id)

}
