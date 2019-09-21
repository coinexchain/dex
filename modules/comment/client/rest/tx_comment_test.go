package rest

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/comment/internal/types"
)

func Test1(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	respWr := restutil.NewResponseWriter4UT()
	//respWr.ClearBody()

	newThreadReq := &NewThreadReq{
		Token:       "cet",
		Donation:    "100",
		Title:       "I love BTC.",
		Content:     "This is the content.",
		ContentType: types.UTF8Text,
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	msg, _ := newThreadReq.GetMsg(nil, addr)
	_ = msg.(*types.MsgCommentToken)
	correct, _ := json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "cet",
		Donation:    100,
		Title:       "I love BTC.",
		Content:     []byte("This is the content."),
		ContentType: types.UTF8Text,
	})
	msgStr, _ := json.Marshal(msg)
	assert.Equal(t, string(correct), string(msgStr))
	assert.Equal(t, 0, len(respWr.GetBody()))

	newThreadReq = &NewThreadReq{
		Token:       "cet",
		Donation:    "100a",
		Title:       "I love BTC.",
		Content:     "This is the content.",
		ContentType: types.UTF8Text,
	}
	_, err := newThreadReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid donation amount", err.Error())

	newThreadReq = &NewThreadReq{
		Token:       "cet",
		Donation:    "100",
		Title:       "I love BTC.",
		Content:     "This is the content.",
		ContentType: types.ShortHanziLZ4,
	}
	_, err = newThreadReq.GetMsg(nil, addr)
	assert.Equal(t, "ShortHanziLZ4 is not valid for rest", err.Error())

	followupCommentReq := &FollowupCommentReq{
		Token:        "btc",
		Donation:     "0",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.UTF8Text,
		IDRewarded:   "10001",
		RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
		RewardToken:  "cet",
		RewardAmount: "10",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	msg, _ = followupCommentReq.GetMsg(nil, addr)
	_ = msg.(*types.MsgCommentToken)
	addr2, _ := sdk.AccAddressFromBech32("coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa")
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "btc",
		Donation:    0,
		Title:       "I love cet too.",
		Content:     []byte("CET to da mars!!!"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr2,
				RewardToken:  "cet",
				RewardAmount: 10,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
		},
	})
	msgStr, _ = json.Marshal(msg)
	assert.Equal(t, string(correct), string(msgStr))
	assert.Equal(t, 0, len(respWr.GetBody()))

	followupCommentReq = &FollowupCommentReq{
		Token:        "btc",
		Donation:     "0x100",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.UTF8Text,
		IDRewarded:   "10001",
		RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
		RewardToken:  "cet",
		RewardAmount: "10",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	_, err = followupCommentReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid donation amount", err.Error())

	followupCommentReq = &FollowupCommentReq{
		Token:        "btc",
		Donation:     "100",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.ShortHanziLZ4,
		IDRewarded:   "10001",
		RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
		RewardToken:  "cet",
		RewardAmount: "10",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	_, err = followupCommentReq.GetMsg(nil, addr)
	assert.Equal(t, "ShortHanziLZ4 is not valid for rest", err.Error())

	followupCommentReq = &FollowupCommentReq{
		Token:        "btc",
		Donation:     "100",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.UTF8Text,
		IDRewarded:   "a10001",
		RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
		RewardToken:  "cet",
		RewardAmount: "10",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	_, err = followupCommentReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid comment ID", err.Error())

	followupCommentReq = &FollowupCommentReq{
		Token:        "btc",
		Donation:     "100",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.UTF8Text,
		IDRewarded:   "10001",
		RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
		RewardToken:  "cet",
		RewardAmount: "1a0",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	_, err = followupCommentReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid reward amount", err.Error())

	followupCommentReq = &FollowupCommentReq{
		Token:        "btc",
		Donation:     "100",
		Title:        "I love cet too.",
		Content:      "CET to da mars!!!",
		ContentType:  types.UTF8Text,
		IDRewarded:   "10001",
		RewardTarget: "coinex1jv65s3grqf6",
		RewardToken:  "cet",
		RewardAmount: "10",
		Attitudes:    []int32{types.Like, types.Favorite},
	}
	_, err = followupCommentReq.GetMsg(nil, addr)
	assert.Equal(t, "decoding bech32 failed: checksum failed. Expected 6akalk, got 3grqf6.", err.Error())

	rewardCommentsReq := &RewardCommentsReq{
		Token: "cet",
		References: []CommentRef{
			{
				ID:           "10001",
				RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
				RewardToken:  "cet",
				RewardAmount: "20",
				Attitudes:    []int32{types.Like},
			},
		},
	}
	msg, _ = rewardCommentsReq.GetMsg(nil, addr)
	_ = msg.(*types.MsgCommentToken)
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "cet",
		Donation:    0,
		Title:       "reward-comments",
		Content:     []byte("No-Content"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr2,
				RewardToken:  "cet",
				RewardAmount: 20,
				Attitudes:    []int32{types.Like},
			},
		},
	})
	msgStr, _ = json.Marshal(msg)
	assert.Equal(t, string(correct), string(msgStr))
	assert.Equal(t, 0, len(respWr.GetBody()))

	rewardCommentsReq = &RewardCommentsReq{
		Token: "cet",
		References: []CommentRef{
			{
				ID:           "10a001",
				RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
				RewardToken:  "cet",
				RewardAmount: "20",
				Attitudes:    []int32{types.Like},
			},
		},
	}
	_, err = rewardCommentsReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid comment ID", err.Error())

	rewardCommentsReq = &RewardCommentsReq{
		Token: "cet",
		References: []CommentRef{
			{
				ID:           "10001",
				RewardTarget: "coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa",
				RewardToken:  "cet",
				RewardAmount: "20a",
				Attitudes:    []int32{types.Like},
			},
		},
	}
	_, err = rewardCommentsReq.GetMsg(nil, addr)
	assert.Equal(t, "invalid reward amount", err.Error())

	rewardCommentsReq = &RewardCommentsReq{
		Token: "cet",
		References: []CommentRef{
			{
				ID:           "10001",
				RewardTarget: "coinex1jv65s3grqf6",
				RewardToken:  "cet",
				RewardAmount: "20",
				Attitudes:    []int32{types.Like},
			},
		},
	}
	_, err = rewardCommentsReq.GetMsg(nil, addr)
	assert.Equal(t, "decoding bech32 failed: checksum failed. Expected 6akalk, got 3grqf6.", err.Error())
}
