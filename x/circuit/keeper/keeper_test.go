package keeper_test

import (
	"testing"

	cmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/circuit/keeper"
	"github.com/cosmos/cosmos-sdk/x/circuit/types"
	"gotest.tools/v3/assert"
)

type fixture struct {
	cdc        codec.Codec
	ctx        sdk.Context
	keeper     keeper.Keeper
	mockAddr   string
	mockPerms  types.Permissions
	mockMsgURL string
}

func initFixture(t *testing.T) *fixture {
	mockStoreKey := storetypes.NewKVStoreKey("test")
	mockAddr := "mock_address"
	keeperX := keeper.NewKeeper(mockStoreKey, mockAddr)
	mockMsgURL := "mock_url"
	mockCtx := testutil.DefaultContextWithDB(t, mockStoreKey, storetypes.NewTransientStoreKey("transient_test"))
	ctx := mockCtx.Ctx.WithBlockHeader(cmproto.Header{})
	mockPerms := types.Permissions{
		Level: 3,
	}

	return &fixture{
		ctx:        ctx,
		keeper:     keeperX,
		mockAddr:   mockAddr,
		mockPerms:  mockPerms,
		mockMsgURL: mockMsgURL,
	}
}

func TestGetAuthority(t *testing.T) {
	t.Parallel()
	f := initFixture(t)
	authority := f.keeper.GetAuthority()
	assert.Equal(t, string(f.mockAddr), authority)
}

// require.Equal(suite.T(), suite.mockAddr.String(), suite.keeper.GetAuthority())
func TestGetAndSetPermissions(t *testing.T) {
	t.Parallel()
	f := initFixture(t)
	// Set the permissions for the mock address.

	err := f.keeper.SetPermissions(f.ctx, f.mockAddr, &f.mockPerms)

	//// Retrieve the permissions for the mock address.
	perms, err := f.keeper.GetPermissions(f.ctx, f.mockAddr)
	require.NoError(t, err)

	//// Assert that the retrieved permissions match the expected value.
	require.Equal(t, &f.mockPerms, perms)

}

func TestIteratePermissions(t *testing.T) {
	t.Parallel()
	f := initFixture(t)
	// Define a set of mock permissions
	mockPerms := []types.Permissions{
		{Level: types.Permissions_LEVEL_SOME_MSGS, LimitTypeUrls: []string{"url1", "url2"}},
		{Level: types.Permissions_LEVEL_ALL_MSGS},
		{Level: types.Permissions_LEVEL_NONE_UNSPECIFIED},
	}

	// Set the permissions for a set of mock addresses
	mockAddrs := []string{
		"mock_address_1",
		"mock_address_2",
		"mock_address_3",
	}
	for i, addr := range mockAddrs {
		f.keeper.SetPermissions(f.ctx, addr, &mockPerms[i])
	}

	// Define a variable to store the returned permissions
	var returnedPerms []types.Permissions

	// Iterate through the permissions and append them to the returnedPerms slice
	f.keeper.IteratePermissions(f.ctx, func(address []byte, perms types.Permissions) (stop bool) {
		returnedPerms = append(returnedPerms, perms)
		return false
	})

	// Assert that the returned permissions match the set mock permissions
	require.Equal(t, mockPerms, returnedPerms)
}
