package comment

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/coinexchain/dex/modules/comment/internal/keepers"
	"github.com/coinexchain/dex/modules/comment/internal/types"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var logStrList = make([]string, 0, 100)

func logStrClear() {
	logStrList = logStrList[:0]
}

func logStrAppend(s string) {
	logStrList = append(logStrList, s)
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func getRefs() []types.CommentRef {
	return []types.CommentRef{
		{
			ID:           900,
			RewardTarget: simpleAddr("00002"),
			RewardToken:  "cet",
			RewardAmount: 10000,
			Attitudes:    []int32{types.Like, types.Favorite},
		},
		{
			ID:           901,
			RewardTarget: simpleAddr("00003"),
			RewardToken:  "usdt",
			RewardAmount: 10,
			Attitudes:    []int32{types.Laugh, types.Favorite},
		},
	}
}

type mocBankxKeeper struct {
	maxAmount sdk.Int
}

func (k *mocBankxKeeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(amt))
	for i, coin := range amt {
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(types.CodeSpaceComment, 999, "Not enough coins")
		}
		coinStrList[i] = coin.Amount.String() + coin.Denom
	}
	s := "Send " + strings.Join(coinStrList, ",") + " from " + from.String() + " to " + to.String()
	logStrAppend(s)
	return nil
}

type mocAssetStatusKeeper struct {
	assets map[string]bool
}

func (k *mocAssetStatusKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	_, ok := k.assets[denom]
	return ok
}

type mocDistributionxKeeper struct {
	poolName  string
	maxAmount sdk.Int
}

func (k *mocDistributionxKeeper) DonateToCommunityPool(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(coins))
	for i, coin := range coins {
		coinStrList[i] = coin.Amount.String() + coin.Denom
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(types.CodeSpaceComment, 999, "Not enough coins")
		}
	}
	s := "Add " + strings.Join(coinStrList, ",") + " to " + k.poolName
	logStrAppend(s)
	return nil
}

func msgSend(key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s := "Msg(" + key + "): " + string(bytes)
	logStrAppend(s)
	return nil
}

func newContextAndKeeper(chainid string) (sdk.Context, *Keeper) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey(StoreKey)
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
	cck := keepers.NewCommentCountKeeper(key)
	k := keepers.NewKeeper(cck,
		&mocBankxKeeper{maxAmount: sdk.NewInt(100)},
		&mocAssetStatusKeeper{assets: map[string]bool{"usdt": true, "btc": true, "cet": true}},
		&mocDistributionxKeeper{poolName: "comPool", maxAmount: sdk.NewInt(100)},
		msgSend,
	)
	return ctx, k
}

func testGenesis(ctx sdk.Context, keeper *Keeper) {
	InitGenesis(ctx, *keeper, DefaultGenesisState())
	gns := ExportGenesis(ctx, *keeper)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", gns.CommentCount))
	gns = NewGenesisState(100)
	InitGenesis(ctx, *keeper, gns)
	if err := gns.Validate(); err != nil {
		logStrAppend("Genesis state is invalid")
	} else {
		logStrAppend("Genesis state is valid")
	}
	gns = ExportGenesis(ctx, *keeper)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", gns.CommentCount))
}

type MsgCreateTradingPair struct {
	Stock          string         `json:"stock"`
	Money          string         `json:"money"`
	Creator        sdk.AccAddress `json:"creator"`
	PricePrecision byte           `json:"price_precision"`
}

func (msg MsgCreateTradingPair) Route() string { return "market" }

func (msg MsgCreateTradingPair) Type() string { return "create_market_info" }

func (msg MsgCreateTradingPair) ValidateBasic() sdk.Error { return nil }

func (msg MsgCreateTradingPair) GetSignBytes() []byte { return nil }

func (msg MsgCreateTradingPair) GetSigners() []sdk.AccAddress { return nil }

func Test1(t *testing.T) {
	ctx, keeper := newContextAndKeeper("test-1")
	logStrClear()
	testGenesis(ctx, keeper)

	msgHandler := NewHandler(*keeper)
	msgCTP := &MsgCreateTradingPair{
		Stock:          "cet",
		Money:          "usdt",
		Creator:        simpleAddr("00200"),
		PricePrecision: 10,
	}
	res := msgHandler(ctx, msgCTP)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.Cck.GetCommentCount(ctx)))
	if res.IsOK() {
		t.Errorf("This should be a failed Result!")
	} else {
		logStrList = append(logStrList, res.Log)
	}

	s := "http://google.com"
	refs := getRefs()
	msg := types.NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, types.HTTP, refs)

	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.Cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if res.IsOK() {
		t.Errorf("This should be a fail result! " + res.Log)
	}

	msg.References[0].RewardAmount = 0
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.Cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if !res.IsOK() {
		t.Errorf("This should be a OK result! " + res.Log)
	}

	msg.Donation = 1000
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.Cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if res.IsOK() {
		t.Errorf("This should be a fail result! " + res.Log)
	}

	msg.Donation = 10
	msg.Token = "bnb"
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.Cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if res.IsOK() {
		t.Errorf("This should be a fail result! " + res.Log)
	}

	for _, s := range logStrList {
		//if refLogs[i]!=s {
		//	t.Errorf("Log String Mismatch!")
		//}
		fmt.Println(s)
	}
}
