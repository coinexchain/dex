package cli

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/comment/internal/types"
)

var ResultMsg *types.MsgCommentToken

func CliRunCommandForTest(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
	ResultMsg = msg.(*types.MsgCommentToken)
	return nil
}

func Test1(t *testing.T) {
	cliutil.CliRunCommand = CliRunCommandForTest

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	args := []string{
		"new-thread",
		"--token=cet",
		"--donation=2000000",
		`--title=I love cet.`,
		`--content=CET to da moon!!!`,
		"--content-type=UTF8Text",
	}
	cmd := CreateNewThreadCmd(nil)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	correct, _ := json.Marshal(&types.MsgCommentToken{
		Sender:      []byte{},
		Token:       "cet",
		Donation:    2000000,
		Title:       "I love cet.",
		Content:     []byte("CET to da moon!!!"),
		ContentType: types.UTF8Text,
	})
	msgStr, _ := json.Marshal(ResultMsg)
	assert.EqualValues(t, string(correct), string(msgStr))

	args = []string{
		"follow-up",
		"--token=cet",
		"--donation=0",
		`--title=I love cet too.`,
		`--content=CET to da mars!!!`,
		`--follow=10001;coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;cet;like,favorite`,
		"--content-type=UTF8Text",
	}
	cmd = CreateFollowupCommentCmd(nil)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      []byte{},
		Token:       "cet",
		Donation:    0,
		Title:       "I love cet too.",
		Content:     []byte("CET to da mars!!!"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr,
				RewardToken:  "cet",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.EqualValues(t, string(correct), string(msgStr))

	args = []string{
		"reward-comments",
		"--token=cet",
		`--reward-to=10001;coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;cet;like,favorite`,
		`--reward-to=20021;coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa;1;cet;like`,
	}
	cmd = RewardCommentsCmd(nil)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	addr2, _ := sdk.AccAddressFromBech32("coinex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa")
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      []byte{},
		Token:       "cet",
		Donation:    0,
		Title:       "",
		Content:     []byte("No-Content"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr,
				RewardToken:  "cet",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
			{
				ID:           20021,
				RewardTarget: addr2,
				RewardToken:  "cet",
				RewardAmount: 1,
				Attitudes:    []int32{types.Like},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.EqualValues(t, string(correct), string(msgStr))

	args = []string{
		"reward-comments",
		"--token=cet",
		`--reward-to=10001;coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;cet;like,favorite`,
	}
	cmd = RewardCommentsCmd(nil)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      []byte{},
		Token:       "cet",
		Donation:    0,
		Title:       "reward-comments",
		Content:     []byte("No-Content"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr,
				RewardToken:  "cet",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.EqualValues(t, string(correct), string(msgStr))
}
