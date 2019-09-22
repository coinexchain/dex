package keepers

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/comment/internal/types"
)

func newContextAndKeeper(chainid string) (sdk.Context, *Keeper) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey(types.StoreKey)

	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	keeper := NewKeeper(key, nil, nil, nil, nil, "")

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())

	return ctx, keeper
}

func Test1(t *testing.T) {
	ctx, keeper := newContextAndKeeper("Test-1")
	keeper.SetCommentCount(ctx, "cet", 1)
	keeper.SetCommentCount(ctx, "btc", 2)
	keeper.IncrCommentCount(ctx, "btc")
	count := keeper.GetCommentCount(ctx, "btc")
	require.Equal(t, uint64(3), count)
	count = keeper.GetCommentCount(ctx, "cet")
	require.Equal(t, uint64(1), count)
	m := keeper.GetAllCommentCount(ctx)
	require.Equal(t, uint64(3), m["btc"])
	require.Equal(t, uint64(1), m["cet"])
}

