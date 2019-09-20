package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/comment/internal/types"
)

type NewThreadReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Token       string       `json:"token"`
	Donation    string       `json:"donation"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	ContentType int8         `json:"content_type"`
}

var _ restutil.RestReq = &NewThreadReq{}

func (req *NewThreadReq) New() restutil.RestReq {
	return new(NewThreadReq)
}
func (req *NewThreadReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *NewThreadReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	donation, err := strconv.ParseInt(req.Donation, 10, 63)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Donation Amount.")
		return nil
	}

	if req.ContentType == types.ShortHanziLZ4 {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "ShortHanziLZ4 is not valid for rest.")
		return nil
	}

	return types.NewMsgCommentToken(sender, req.Token, donation, req.Title, req.Content, req.ContentType, nil)
}

func createNewThreadHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req NewThreadReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}

type FollowupCommentReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Token        string       `json:"token"`
	Donation     string       `json:"donation"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	ContentType  int8         `json:"content_type"`
	IDRewarded   string       `json:"id_rewarded"`
	RewardTarget string       `json:"reward_target"`
	RewardToken  string       `json:"reward_token"`
	RewardAmount string       `json:"reward_amount"`
	Attitudes    []int32      `json:"attitudes"`
}

var _ restutil.RestReq = &FollowupCommentReq{}

func (req *FollowupCommentReq) New() restutil.RestReq {
	return new(FollowupCommentReq)
}
func (req *FollowupCommentReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *FollowupCommentReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	donation, err := strconv.ParseInt(req.Donation, 10, 63)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Donation Amount.")
		return nil
	}

	crefs := make([]types.CommentRef, 1)
	idRewarded, err := strconv.ParseInt(req.IDRewarded, 10, 63)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Comment ID.")
		return nil
	}
	rewardTarget, err := sdk.AccAddressFromBech32(req.RewardTarget)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return nil
	}
	rewardAmount, err := strconv.ParseInt(req.RewardAmount, 10, 63)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Reward Amount.")
		return nil
	}
	crefs[0].ID = uint64(idRewarded)
	crefs[0].RewardTarget = rewardTarget
	crefs[0].RewardToken = req.RewardToken
	crefs[0].RewardAmount = rewardAmount
	crefs[0].Attitudes = req.Attitudes

	if req.ContentType == types.ShortHanziLZ4 {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "ShortHanziLZ4 is not valid for rest.")
		return nil
	}

	return types.NewMsgCommentToken(sender, req.Token, donation, req.Title, req.Content, req.ContentType, crefs)
}

func createFollowupCommentHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req FollowupCommentReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}

type CommentRef struct {
	ID           string  `json:"id"`
	RewardTarget string  `json:"reward_target"`
	RewardToken  string  `json:"reward_token"`
	RewardAmount string  `json:"reward_amount"`
	Attitudes    []int32 `json:"attitudes"`
}
type RewardCommentsReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Token      string       `json:"token"`
	References []CommentRef `json:"references"`
}

var _ restutil.RestReq = &RewardCommentsReq{}

func (req *RewardCommentsReq) New() restutil.RestReq {
	return new(RewardCommentsReq)
}
func (req *RewardCommentsReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *RewardCommentsReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	crefs := make([]types.CommentRef, len(req.References))
	for i, r := range req.References {
		idRewarded, err := strconv.ParseInt(r.ID, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Comment ID.")
			return nil
		}
		rewardTarget, err := sdk.AccAddressFromBech32(r.RewardTarget)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return nil
		}
		rewardAmount, err := strconv.ParseInt(r.RewardAmount, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Reward Amount.")
			return nil
		}
		crefs[i].ID = uint64(idRewarded)
		crefs[i].RewardTarget = rewardTarget
		crefs[i].RewardToken = r.RewardToken
		crefs[i].RewardAmount = rewardAmount
		crefs[i].Attitudes = r.Attitudes
	}

	msg := types.NewMsgCommentToken(sender, req.Token, 0, "", "", types.UTF8Text, crefs)
	if len(msg.References) <= 1 && len(msg.Title) == 0 {
		msg.Title = "reward-comments"
	}
	if len(msg.Content) == 0 {
		msg.Content = []byte("No-Content")
	}
	return msg
}

func createRewardCommentsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req RewardCommentsReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}
