package simulation

import (
	"fmt"

	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/comment"
	"github.com/coinexchain/dex/modules/comment/internal/keepers"
	"github.com/coinexchain/dex/modules/comment/internal/types"
	simulationx "github.com/coinexchain/dex/simulation"
)

// TODO
func SimulateCreateNewThread(k keepers.Keeper, ask asset.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {
		//create new token comment msg
		msg, err := createNewThread(r, ctx, k, ask, ak, accs)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		//get #token-comment
		lastCommnet := k.GetCommentCount(ctx, msg.Token)

		//handle msg
		handler := comment.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		//verify msg is correctly handled
		ok = verifyCreateNewThread(ctx, k, msg, lastCommnet)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("new token comment creation failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func createNewThread(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ask asset.Keeper, ak auth.AccountKeeper, accs []simulation.Account) (types.MsgCommentToken, error) {
	fromAcc := simulation.RandomAcc(r, accs)

	//randomly select token to comment
	token := randomToken(r, ctx, ask)
	if token == nil {
		return types.MsgCommentToken{}, fmt.Errorf("no token to comment")
	}

	//randomly select amount of cet to donate
	donation := simulationx.RandomCET(r, ctx, ak, fromAcc)

	//generate title
	title := randomUTF8OrBytes(r, r.Intn(types.MaxTitleSize), true)

	//randomly generate content & contentType
	contentType, content := randomContent(r)

	//construct msgCommentToken to create a new comment
	msg := types.NewMsgCommentToken(fromAcc.Address, token.GetSymbol(), donation, string(title), string(content), contentType, nil)
	if msg.ValidateBasic() != nil {
		return types.MsgCommentToken{}, fmt.Errorf("msg expected to pass validation check")
	}
	return *msg, nil

}

func verifyCreateNewThread(ctx sdk.Context, k keepers.Keeper, msg types.MsgCommentToken, lastComment uint64) bool {
	comment := k.GetCommentCount(ctx, msg.Token)
	return comment == lastComment+1

}

func SimulateCreateCommentRefs(k keepers.Keeper, ask asset.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {
		//create new token comment msg
		msg, err := createCommentReferences(r, ctx, k, ask, ak, accs)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		//get #token-comment
		lastCommnet := k.GetCommentCount(ctx, msg.Token)
		oldCoins := make(map[string]sdk.Coins)
		for _, ref := range msg.References {
			oldAcc := ak.GetAccount(ctx, ref.RewardTarget)
			oldCoins[oldAcc.GetAddress().String()] = oldAcc.GetCoins()
		}

		//handle msg
		handler := comment.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		//verify msg is correctly handled
		ok = verifyCreateCommentRefs(ctx, k, ak, msg, lastCommnet, oldCoins)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token comment references handle failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func createCommentReferences(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ask asset.Keeper, ak auth.AccountKeeper, accs []simulation.Account) (types.MsgCommentToken, error) {
	fromAcc := simulation.RandomAcc(r, accs)
	fromAuthAcc := ak.GetAccount(ctx, fromAcc.Address)
	//randomly select token to comment

	token, refs := randomCommentRef(r, ctx, k, ask, fromAuthAcc, accs)
	if token == nil {
		return types.MsgCommentToken{}, fmt.Errorf("no token comment to reference")
	}

	//randomly select amount of cet to donate
	donation := simulationx.RandomCET(r, ctx, ak, fromAcc)

	//construct msgCommentToken to create a new comment
	msg := types.NewMsgCommentToken(fromAcc.Address, token.GetSymbol(), donation, "", "", types.RawBytes, refs)
	if msg.ValidateBasic() != nil {
		return types.MsgCommentToken{}, fmt.Errorf("msg expected to pass validation check")
	}
	return *msg, nil

}

func randomCommentRef(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ask asset.Keeper, fromAcc auth.Account, accs []simulation.Account) (token asset.Token, refs []types.CommentRef) {

	totalComment := k.GetAllCommentCount(ctx)
	if totalComment == nil {
		return nil, []types.CommentRef{}
	}

	//generate len CommentRef
	token, ids := randomTokenCommentRef(r, ctx, k, ask)
	for i := 0; i < len(ids); i = i + 1 {

		rewardTarget := simulation.RandomAcc(r, accs)
		denom, amt := simulationx.RandomAccCoins(r, fromAcc)

		attitude := r.Intn(int(types.Condolences-types.Like)) + int(types.Like)

		refs = append(refs, types.CommentRef{
			ID:           ids[i],
			RewardTarget: rewardTarget.Address,
			RewardToken:  denom,
			RewardAmount: amt,
			Attitudes:    []int32{int32(attitude)},
		})
	}

	return
}

func verifyCreateCommentRefs(ctx sdk.Context, k keepers.Keeper, ak auth.AccountKeeper, msg types.MsgCommentToken, lastComment uint64, oldCoins map[string]sdk.Coins) bool {
	comment := k.GetCommentCount(ctx, msg.Token)
	if comment != lastComment+1 {
		return false
	}

	for _, ref := range msg.References {
		newCoins := ak.GetAccount(ctx, ref.RewardTarget).GetCoins()
		rewardCoins := sdk.NewCoins(sdk.NewCoin(ref.RewardToken, sdk.NewInt(ref.RewardAmount)))
		if !oldCoins[ref.RewardTarget.String()].Add(rewardCoins).IsEqual(newCoins) {
			return false
		}
	}

	return true
}
