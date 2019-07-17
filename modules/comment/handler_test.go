package comment

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coinexchain/dex/modules/comment/shorthanzi"
	"github.com/coinexchain/dex/modules/market"
	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"strings"
	"testing"
)

var logStrList = make([]string, 0, 100)

func logStrClear() {
	logStrList = logStrList[:0]
}

func logStrAppend(s string) {
	logStrList = append(logStrList, s)
}

type mocBankxKeeper struct {
	maxAmount sdk.Int
}

func (k *mocBankxKeeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(amt))
	for i, coin := range amt {
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(CodeSpaceComment, 999, "Not enough coins")
		}
		coinStrList[i] = coin.Amount.String() + coin.Denom
	}
	s := "Subtract " + strings.Join(coinStrList, ",") + " from " + addr.String()
	logStrAppend(s)
	return nil
}

func (k *mocBankxKeeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(amt))
	for i, coin := range amt {
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(CodeSpaceComment, 999, "Not enough coins")
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

type mocDistributionKeeper struct {
	poolName string
}

func (k *mocDistributionKeeper) AddCoinsToFeePool(ctx sdk.Context, coins sdk.Coins) {
	coinStrList := make([]string, len(coins))
	for i, coin := range coins {
		coinStrList[i] = coin.Amount.String() + coin.Denom
	}
	s := "Add " + strings.Join(coinStrList, ",") + " to " + k.poolName
	logStrAppend(s)
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
	cck := NewCommentCountKeeper(key)
	k := NewKeeper(cck,
		&mocBankxKeeper{maxAmount: sdk.NewInt(100)},
		&mocAssetStatusKeeper{assets: map[string]bool{"usdt": true, "btc": true, "cet": true}},
		&mocDistributionKeeper{poolName: "comPool"},
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

func testParseContentType() {
	inList := []string{"ipfs", "magnet", "http", "utf8text", "shorthanzi", "rawbytes", "fuck"}
	outList := make([]string, len(inList))
	for i, s := range inList {
		outList[i] = fmt.Sprintf("%s:%d", s, ParseContentType(s))
	}
	logStrAppend(strings.Join(outList, ","))
}

func testParseAttitude() {
	inList := []string{"like", "dislike", "laugh", "cry", "angry", "surprise", "heart", "sweat",
		"speechless", "favorite", "condolences", "fuck"}
	outList := make([]string, len(inList))
	for i, s := range inList {
		outList[i] = fmt.Sprintf("%s:%d", s, ParseAttitude(s))
	}
	logStrAppend(strings.Join(outList, ","))
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func getRefs() []CommentRef {
	return []CommentRef{
		{
			ID:           900,
			RewardTarget: simpleAddr("00002"),
			RewardToken:  "cet",
			RewardAmount: 10000,
			Attitudes:    []int32{Like, Favorite},
		},
		{
			ID:           901,
			RewardTarget: simpleAddr("00003"),
			RewardToken:  "usdt",
			RewardAmount: 10,
			Attitudes:    []int32{Laugh, Favorite},
		},
	}
}

func Test1(t *testing.T) {
	logStrClear()
	testParseContentType()
	testParseAttitude()
	refs := getRefs()
	////// ShortHanzi
	msg := NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text0, ShortHanzi, refs)

	if res := msg.ValidateBasic(); res != nil {
		fmt.Println(res.ABCILog())
		t.Errorf("This should be a valid Msg!")
	}

	tc := NewTokenComment(msg, 108)
	if msg.ContentType != ShortHanzi || shorthanzi.Text0 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// ShortHanziLZ4
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text1, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108)
	if msg.ContentType != ShortHanziLZ4 || shorthanzi.Text1 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text2, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108)
	if msg.ContentType != ShortHanziLZ4 || shorthanzi.Text2 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// RawBytes
	s := base64.StdEncoding.EncodeToString([]byte("大获全胜"))
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, RawBytes, refs)
	tc = NewTokenComment(msg, 108)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != RawBytes {
		t.Errorf("Invalid Token Comment!")
	}

	////// UTF8Text
	s = "孜孜不倦"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, UTF8Text, refs)
	tc = NewTokenComment(msg, 108)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// HTTP
	s = "http://google.com"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	tc = NewTokenComment(msg, 108)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != HTTP {
		t.Errorf("Invalid Token Comment!")
	}

	//len(msg.Sender) == 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Sender = nil
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if len(msg.Token) == 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Token = ""
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if msg.Donation < 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Donation = -1
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if len(msg.Title) == 0 || len(msg.References) <= 1 { return ErrNoTitle() }
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Title = ""
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if msg.ContentType < IPFS || msg.ContentType > ShortHanziLZ4
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.ContentType = 100
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if !utf8.Valid(msg.Content) { return ErrInvalidContent() }
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text2, ShortHanzi, refs)
	msg.ContentType = ShortHanzi
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}

	//if len(msg.Content) > MaxContentSize
	text := shorthanzi.Text3 + shorthanzi.Text3
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", text, UTF8Text, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf(fmt.Sprintf("This should be an invalid Msg %d", len(text)))
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if a < Like || a > Condolences { return ErrInvalidAttitude(a) }
	refs = getRefs()
	refs[0].Attitudes = []int32{100}
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if ref.RewardAmount < 0 { return ErrNegativeReward() }
	refs = getRefs()
	refs[1].RewardAmount = -1
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}

	refLogs := []string{
		`ipfs:0,magnet:1,http:2,utf8text:3,shorthanzi:4,rawbytes:6,fuck:-1`,
		`like:50,dislike:51,laugh:52,cry:53,angry:54,surprise:55,heart:56,sweat:57,speechless:58,favorite:59,condolences:60,fuck:-1`,
		`{"codespace":"sdk","code":7,"message":"missing sender address"}`,
		`{"codespace":"comment","code":901,"message":"Invalid Symbol"}`,
		`{"codespace":"comment","code":902,"message":"Donation can not be negative"}`,
		`{"codespace":"comment","code":903,"message":"No summary is provided"}`,
		`{"codespace":"comment","code":904,"message":"'100' is not a valid content type"}`,
		`{"codespace":"comment","code":905,"message":"Content has invalid format"}`,
		`{"codespace":"comment","code":906,"message":"Content is larger than 16384 bytes"}`,
		`{"codespace":"comment","code":907,"message":"'100' is not a valid attitude"}`,
		`{"codespace":"comment","code":908,"message":"Reward can not be negative"}`,
	}
	for i, s := range logStrList {
		if refLogs[i]!=s {
			t.Errorf("Log String Mismatch!")
		}
		fmt.Println(s)
	}
}

func Test2(t *testing.T) {
	ctx, keeper := newContextAndKeeper("test-1")
	logStrClear()
	testGenesis(ctx, keeper)

	msgHandler := NewHandler(*keeper)
	msgCTP := &market.MsgCreateTradingPair{
		Stock:          "cet",
		Money:          "usdt",
		Creator:        simpleAddr("00200"),
		PricePrecision: 10,
	}
	res := msgHandler(ctx, msgCTP)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.cck.GetCommentCount(ctx)))
	if res.IsOK() {
		t.Errorf("This should be a failed Result!")
	} else {
		logStrList = append(logStrList, res.Log)
	}

	s := "http://google.com"
	refs := getRefs()
	msg := NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)

	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if res.IsOK() {
		t.Errorf("This should be a fail result! " + res.Log)
	}

	msg.References[0].RewardAmount = 0
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if !res.IsOK() {
		t.Errorf("This should be a OK result! " + res.Log)
	}

	msg.Donation = 1000
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.cck.GetCommentCount(ctx)))
	logStrList = append(logStrList, res.Log)
	if res.IsOK() {
		t.Errorf("This should be a fail result! " + res.Log)
	}

	msg.Donation = 10
	msg.Token = "bnb"
	res = msgHandler(ctx, *msg)
	logStrAppend(fmt.Sprintf("Now comment count is: %d", keeper.cck.GetCommentCount(ctx)))
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
