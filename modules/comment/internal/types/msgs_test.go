package types

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/coinexchain/dex/modules/comment/shorthanzi"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var logStrList = make([]string, 0, 100)

func logStrClear() {
	logStrList = logStrList[:0]
}

func logStrAppend(s string) {
	logStrList = append(logStrList, s)
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

	tc := NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanzi || shorthanzi.Text0 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// ShortHanziLZ4
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text1, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanziLZ4 || shorthanzi.Text1 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", shorthanzi.Text2, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanziLZ4 || shorthanzi.Text2 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// RawBytes
	s := base64.StdEncoding.EncodeToString([]byte("大获全胜"))
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, RawBytes, refs)
	tc = NewTokenComment(msg, 108, 1000)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != RawBytes {
		t.Errorf("Invalid Token Comment!")
	}

	////// UTF8Text
	s = "孜孜不倦"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, UTF8Text, refs)
	tc = NewTokenComment(msg, 108, 1000)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// HTTP
	s = "http://google.com"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	tc = NewTokenComment(msg, 108, 1000)
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
	//if len(msg.Title) == 0 && len(msg.References) <= 1 { return ErrNoTitle() }
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "", s, HTTP, []CommentRef{refs[0]})
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
		`{"codespace":"comment","code":903,"message":"No title is provided"}`,
		`{"codespace":"comment","code":904,"message":"'100' is not a valid content type"}`,
		`{"codespace":"comment","code":905,"message":"Content has invalid format"}`,
		`{"codespace":"comment","code":906,"message":"Content is larger than 16384 bytes"}`,
		`{"codespace":"comment","code":907,"message":"'100' is not a valid attitude"}`,
		`{"codespace":"comment","code":908,"message":"Reward can not be negative"}`,
	}
	for i, s := range logStrList {
		if refLogs[i] != s {
			t.Errorf("Log String Mismatch!")
		}
		fmt.Println(s)
	}
}
