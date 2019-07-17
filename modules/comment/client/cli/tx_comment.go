package cli

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/coinexchain/dex/modules/comment"
)

const (
	FlagToken       = "token"
	FlagDonation    = "donation"
	FlagTitle       = "title"
	FlagContent     = "content"
	FlagContentType = "content-type"
	FlagRewardTo    = "reward-to"
	FlagFollow      = "follow"
)

var createNewThreadFlags = []string{
	FlagToken,
	FlagDonation,
	FlagTitle,
	FlagContent,
	FlagContentType,
}

var createFollowupCommentFlags = []string{
	FlagToken,
	FlagDonation,
	FlagTitle,
	FlagContent,
	FlagContentType,
	FlagFollow,
}

var rewardCommentsFlags = []string{
	FlagToken,
	FlagRewardTo,
}

func CreateNewThreadCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-thread",
		Short: "Create a new thread of comments under some token",
		Long: `Post a comment under some token, which creates a new thread, instead of following any other comments.

Example: 
	 cetcli tx comment new-thread --token=cet --donation=2 
	 --title="I love cet." --content="CET to da moon!!!" --content-type=UTF8Text
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadcastComment(cdc, "new-thread", nil)
		},
	}

	markCreateNewThreadFlags(cmd)
	return cmd
}

func markCreateNewThreadFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagToken, "cet", "The token you want to comment about")
	cmd.Flags().Int(FlagDonation, 0, "The donation to community pool. The more you donate, the more your comment weights.")
	cmd.Flags().String(FlagTitle, "", "The comment's title")
	cmd.Flags().String(FlagContent, "", "The comment's content")
	cmd.Flags().String(FlagContentType, "UTF8Text", "The type of the comment's content (IPFS, Magnet, HTTP, UTF8Text, ShortHanzi or RawBytes)")

	for _, flag := range createNewThreadFlags {
		cmd.MarkFlagRequired(flag)
	}
}

func CreateFollowupCommentCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow-up",
		Short: "Create a follow-up comment in a thread",
		Long: `Post a comment to follow another comment in a thread.

Example: 
	 cetcli tx comment follow-up --token=cet --donation=0 --follow="10001;coinex1qw508d6qejxtdg4y5r2zarvary0c0xw9kv8f3t4;2;cet;like,favorite"
	 --title="I love cet too." --content="CET to da moon!!!" --content-type=UTF8Text
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadcastComment(cdc, "follow-up", nil)
		},
	}

	markCreateFollowupCommentFlags(cmd)
	return cmd
}

func markCreateFollowupCommentFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagToken, "cet", "The token you want to comment about")
	cmd.Flags().Int(FlagDonation, 0, "The donation to community pool. If you have negative opinion against the comment you are following, please donate some cet tokens to prove your dissatisfaction.")
	cmd.Flags().String(FlagFollow, "", "Some information about the comment you are following. Should be like this: \"<comment-id>;<the-sender-of-the-comment>;<reward-amount>;<reward-token>;<comma-separated-attitued-list>\". Valid attitudes include: like, dislike, laugh, cry, angry, surprise, heart, sweat, speechless, favorite, condolences.")
	cmd.Flags().String(FlagTitle, "", "The comment's title")
	cmd.Flags().String(FlagContent, "", "The comment's content")
	cmd.Flags().String(FlagContentType, "UTF8Text", "The type of the comment's content (IPFS, Magnet, HTTP, UTF8Text, ShortHanzi or RawBytes)")

	for _, flag := range createFollowupCommentFlags {
		cmd.MarkFlagRequired(flag)
	}
}

func RewardCommentsCmd(cdc *codec.Codec) *cobra.Command {
	var rewardsArray []string
	cmd := &cobra.Command{
		Use:   "reward-comments",
		Short: "Reward some comments that you like",
		Long: `Reward the senders and some comments that you like, while showing why you like them individually.

Example: 
	 cetcli tx comment reward-comments --token=cet --reward-to="10001;coinex1qi598e62ejitdg4yur3zarvary0c5xw7kv8f3t4;2;cet;like,favorite" --reward-to="20021;coinex1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4;1;cet;like"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadcastComment(cdc, "reward-comments", &rewardsArray)
		},
	}

	cmd.Flags().String(FlagToken, "cet", "The token you want to comment about")
	cmd.Flags().StringArrayVar(&rewardsArray, FlagRewardTo, nil, "You can use this option multiple times to reward multiple comments. This option specify some information about one comment you want to reward. Should be like this: \"<comment-id>;<the-sender-of-the-comment>;<reward-amount>;<reward-token>;<comma-separated-attitued-list>\". Valid attitudes include: like, dislike, laugh, cry, angry, surprise, heart, sweat, speechless, favorite, condolences.")

	cmd.MarkFlagRequired(FlagToken)
	cmd.MarkFlagRequired(FlagRewardTo)
	return cmd
}

func createAndBroadcastComment(cdc *codec.Codec, subcmd string, rewardsArrayPtr *[]string) error {
	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

	sender := cliCtx.GetFromAddress()

	msg, err := parseFlags(sender, rewardsArrayPtr)
	if err != nil {
		return errors.Errorf("tx flag is error, please see help : " +
			"$ cetcli tx comment " + subcmd + " -h")
	}
	if err = msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

func parseRewardLine(line string) (*comment.CommentRef, error) {
	symbols := strings.Split(line, ";")
	if len(symbols) != 5 {
		return nil, errors.Errorf("invalid format: " + line)
	}

	id, err := strconv.ParseInt(symbols[0], 10, 63)
	if err != nil {
		return nil, errors.Errorf("Not a valid comment id: " + symbols[0])
	}

	target, err := sdk.AccAddressFromBech32(symbols[1])
	if err != nil {
		return nil, errors.Errorf("Not a valid address: " + symbols[1])
	}

	amt, err := strconv.ParseInt(symbols[3], 10, 63)
	if err != nil {
		return nil, errors.Errorf("Not a valid amount: " + symbols[3])
	}

	attitudes := strings.Split(symbols[4], ",")
	attList := make([]int32, len(attitudes))
	for i, a := range attitudes {
		attList[i] = comment.ParseAttitude(a)
		if attList[i] < 0 {
			return nil, errors.Errorf("invalid attitude: " + a)
		}
	}

	cref := &comment.CommentRef{
		ID:           uint64(id),
		RewardTarget: target,
		RewardToken:  symbols[2],
		RewardAmount: amt,
		Attitudes:    attList,
	}
	return cref, nil
}

func parseFlags(sender sdk.AccAddress, rewardsArrayPtr *[]string) (*comment.MsgCommentToken, error) {
	ctstr := viper.GetString(FlagContentType)
	ct := comment.ParseContentType(ctstr)
	if ct < 0 {
		return nil, errors.Errorf(ctstr + " is not a valid content type.")
	}

	var references []comment.CommentRef
	followup := viper.GetString(FlagFollow)
	if len(followup) != 0 {
		cref, err := parseRewardLine(followup)
		if err != nil {
			return nil, err
		}
		references = []comment.CommentRef{*cref}
	} else {
		references = make([]comment.CommentRef, 0, len(*rewardsArrayPtr))
		for _, line := range *rewardsArrayPtr {
			cref, err := parseRewardLine(line)
			if err != nil {
				return nil, err
			}
			references = append(references, *cref)
		}
	}

	token := viper.GetString(FlagToken)
	donation := viper.GetInt64(FlagDonation)
	title := viper.GetString(FlagTitle)
	content := viper.GetString(FlagContent)
	return comment.NewMsgCommentToken(sender, token, donation, title, content, ct, references), nil
}
