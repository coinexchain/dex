package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/comment"
)

type NewThreadReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Token       string       `json:"token"`
	Donation    string       `json:"donation"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	ContentType int8         `json:"content_type"`
}

func createNewThreadHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NewThreadReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence := req.BaseReq.Sequence
		if sequence == 0 {
			_, sequence, err = auth.NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
				return
			}
		}
		req.BaseReq.Sequence = sequence

		donation, err := strconv.ParseInt(req.Donation, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Donation Amount.")
			return
		}

		if req.ContentType == comment.ShortHanziLZ4 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "ShortHanziLZ4 is not valid for rest.")
			return
		}

		msg := comment.NewMsgCommentToken(sender, req.Token, donation, req.Title, req.Content, req.ContentType, nil)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

func createFollowupCommentHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FollowupCommentReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence := req.BaseReq.Sequence
		if sequence == 0 {
			_, sequence, err = auth.NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
				return
			}
		}
		req.BaseReq.Sequence = sequence

		donation, err := strconv.ParseInt(req.Donation, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Donation Amount.")
			return
		}

		crefs := make([]comment.CommentRef, 1)
		idRewarded, err := strconv.ParseInt(req.IDRewarded, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Comment ID.")
			return
		}
		rewardTarget, err := sdk.AccAddressFromBech32(req.RewardTarget)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rewardAmount, err := strconv.ParseInt(req.RewardAmount, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Reward Amount.")
			return
		}
		crefs[0].ID = uint64(idRewarded)
		crefs[0].RewardTarget = rewardTarget
		crefs[0].RewardToken = req.RewardToken
		crefs[0].RewardAmount = rewardAmount
		crefs[0].Attitudes = req.Attitudes

		if req.ContentType == comment.ShortHanziLZ4 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "ShortHanziLZ4 is not valid for rest.")
			return
		}

		msg := comment.NewMsgCommentToken(sender, req.Token, donation, req.Title, req.Content, req.ContentType, crefs)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

func createRewardCommentsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RewardCommentsReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence := req.BaseReq.Sequence
		if sequence == 0 {
			_, sequence, err = auth.NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
				return
			}
		}
		req.BaseReq.Sequence = sequence

		crefs := make([]comment.CommentRef, len(req.References))
		for i, r := range req.References {
			idRewarded, err := strconv.ParseInt(r.ID, 10, 63)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Comment ID.")
				return
			}
			rewardTarget, err := sdk.AccAddressFromBech32(r.RewardTarget)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			rewardAmount, err := strconv.ParseInt(r.RewardAmount, 10, 63)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Reward Amount.")
				return
			}
			crefs[i].ID = uint64(idRewarded)
			crefs[i].RewardTarget = rewardTarget
			crefs[i].RewardToken = r.RewardToken
			crefs[i].RewardAmount = rewardAmount
			crefs[i].Attitudes = r.Attitudes
		}

		msg := comment.NewMsgCommentToken(sender, req.Token, 0, "", "", comment.UTF8Text, crefs)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
