package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	simulation2 "github.com/coinexchain/dex/simulation"
)

func SimulateMsgIssueToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		tokenName, tokenSymbol, tokenURL, tokenDescrip, tokenIdentify := randomTokenAttrs(r)
		issueAmount := randomIssuance(r)
		msg := asset.NewMsgIssueToken(tokenName, tokenSymbol, issueAmount, acc.Address,
			true, true, true, true, tokenURL, tokenDescrip, tokenIdentify)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		ok := simulation2.SimulateHandleMsg(msg, asset.NewHandler(k), ctx)

		opMsg := simulation.NewOperationMsg(msg, ok, "")
		if !ok {
			return opMsg, nil, nil
		}
		ok = checkIssueTokenValid(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("issue token failed")
		}
		return opMsg, nil, nil
	}
}

func checkIssueTokenValid(ctx sdk.Context, k asset.Keeper, msg types.MsgIssueToken) bool {
	token := k.GetToken(ctx, msg.Symbol)
	return token.GetBurnable() == msg.Burnable &&
		token.GetMintable() == msg.Mintable &&
		token.GetOwner().Equals(msg.Owner) &&
		token.GetTotalBurn().IsZero() &&
		token.GetTotalMint().IsZero() &&
		token.GetSendLock().IsZero() &&
		!token.GetIsForbidden() &&
		token.GetTokenForbiddable() == msg.TokenForbiddable &&
		token.GetAddrForbiddable() == msg.AddrForbiddable &&
		token.GetDescription() == msg.Description &&
		token.GetURL() == msg.URL &&
		token.GetIdentity() == msg.Identity

}
func getOrGenSymbolOwner(r *rand.Rand, accounts []simulation.Account, ctx sdk.Context, k asset.Keeper) (types.Token, string, sdk.AccAddress) {
	token := RandomToken(r, ctx, k)

	var symbol string
	var owner sdk.AccAddress
	makeValid := simulation2.RandomBool(r)
	if makeValid && token != nil {
		symbol = token.GetSymbol()
		owner = token.GetOwner()

	} else {
		owner = simulation.RandomAcc(r, accounts).Address
		symbol = RandStringBytes(r, SymbolLenth)
	}

	return token, symbol, owner
}
func SimulateMsgTransferOwnership(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		token, symbol, originOwner := getOrGenSymbolOwner(r, accounts, ctx, k)
		if len(accounts) <= 1 || token == nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		newOwner := simulation.RandomAcc(r, accounts).Address
		for newOwner.Equals(originOwner) {
			newOwner = simulation.RandomAcc(r, accounts).Address
		}

		msg := asset.NewMsgTransferOwnership(symbol, originOwner, newOwner)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyTokenOwnerTransfer(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token ownership transfer falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func RandomToken(r *rand.Rand, ctx sdk.Context, k asset.Keeper) asset.Token {
	tokenList := k.GetAllTokens(ctx)
	if len(tokenList) == 0 {
		return nil
	}
	return tokenList[simulation2.GetRandomElemIndex(r, len(tokenList))]

}
func verifyTokenOwnerTransfer(ctx sdk.Context, k asset.Keeper, msg types.MsgTransferOwnership) bool {
	token := k.GetToken(ctx, msg.Symbol)
	return token.GetOwner().Equals(msg.NewOwner)
}
func SimulateMsgMintToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		var oldMint sdk.Int
		token := k.GetToken(ctx, symbol)
		if token != nil {
			oldMint = token.GetTotalMint()
		}

		mintAmount := r.Intn(50000000) + 100000
		msg := asset.NewMsgMintToken(symbol, sdk.NewInt(int64(mintAmount)), owner)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyTokenMint(ctx, k, msg, oldMint)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token mint falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func verifyTokenMint(ctx sdk.Context, k asset.Keeper, msg types.MsgMintToken, oldMint sdk.Int) bool {
	token := k.GetToken(ctx, msg.Symbol)
	return token.GetTotalMint().Sub(oldMint).Equal(msg.Amount)
}
func SimulateMsgBurnToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {
		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		var oldBurn sdk.Int
		token := k.GetToken(ctx, symbol)
		if token != nil {
			oldBurn = token.GetTotalBurn()
		}

		burnAmount := r.Intn(5000000) + 100000
		msg := asset.NewMsgBurnToken(symbol, sdk.NewInt(int64(burnAmount)), owner)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyTokenBurn(ctx, k, msg, oldBurn)

		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token burn falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func verifyTokenBurn(ctx sdk.Context, k asset.Keeper, msg types.MsgBurnToken, oldBurn sdk.Int) bool {
	token := k.GetToken(ctx, msg.Symbol)
	return token.GetTotalBurn().Sub(oldBurn).Equal(msg.Amount)

}
func SimulateMsgForbidToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		msg := asset.NewMsgForbidToken(symbol, owner)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyForbidToken(ctx, k, msg.Symbol)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token forbid falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func verifyForbidToken(ctx sdk.Context, k asset.Keeper, symbol string) bool {
	token := k.GetToken(ctx, symbol)
	return token.GetIsForbidden()
}

func SimulateMsgUnForbidToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {
		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		msg := asset.NewMsgUnForbidToken(symbol, owner)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = !verifyForbidToken(ctx, k, msg.Symbol)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token unforbid falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func SimulateMsgAddTokenWhitelist(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		accLen := r.Intn(len(accounts))
		whiteList := make([]sdk.AccAddress, accLen)
		for i := 0; i < accLen; i++ {
			whiteList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgAddTokenWhitelist(symbol, owner, whiteList)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyTokenWhitelist(ctx, k, msg.Symbol, msg.Whitelist, true)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token add whitelist falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func verifyTokenWhitelist(ctx sdk.Context, k asset.Keeper, symbol string, addrs []sdk.AccAddress, add bool) bool {
	whiteList := k.GetWhitelist(ctx, symbol)
	whiteAddrMap := make(map[string]struct{})
	for _, addr := range whiteList {
		whiteAddrMap[addr.String()] = struct{}{}
	}
	for _, addr := range addrs {
		_, existed := whiteAddrMap[addr.String()]
		if add {
			if !existed {
				return false
			}
		} else {
			if existed {
				return false
			}
		}
	}
	return true
}
func SimulateMsgRemoveTokenWhitelist(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		whiteList := k.GetWhitelist(ctx, symbol)
		removeList := randomSubsliceOfAddr(r, whiteList)

		msg := asset.NewMsgRemoveTokenWhitelist(symbol, owner, removeList)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyTokenWhitelist(ctx, k, msg.Symbol, msg.Whitelist, false)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token remove whitelist falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func randomSubsliceOfAddr(r *rand.Rand, addrs []sdk.AccAddress) []sdk.AccAddress {
	if len(addrs) == 0 {
		return []sdk.AccAddress{}
	}
	randLen := r.Intn(len(addrs))

	var returnList = make([]sdk.AccAddress, randLen)
	for i := 0; i < randLen; i = i + 1 {
		returnList = append(returnList, addrs[simulation2.GetRandomElemIndex(r, len(addrs))])
	}
	return returnList
}

func SimulateMsgForbidAddr(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		accLen := r.Intn(len(accounts))
		forbidList := make([]sdk.AccAddress, accLen)
		for i := 0; i < accLen; i++ {
			forbidList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgForbidAddr(symbol, owner, forbidList)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyForbiddenAddr(ctx, k, msg.Symbol, msg.Addresses, true)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token forbid addr falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func verifyForbiddenAddr(ctx sdk.Context, k asset.Keeper, symbol string, addrs []sdk.AccAddress, forbid bool) bool {
	forbidList := k.GetForbiddenAddresses(ctx, symbol)
	forbiddenAddrMap := make(map[string]struct{})
	for _, addr := range forbidList {
		forbiddenAddrMap[addr.String()] = struct{}{}
	}
	for _, addr := range addrs {
		_, existed := forbiddenAddrMap[addr.String()]
		if forbid {
			if !existed {
				return false
			}
		} else {
			if existed {
				return false
			}
		}
	}
	return true
}
func SimulateMsgUnForbidAddr(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		ForbiddenAddrs := k.GetForbiddenAddresses(ctx, symbol)
		unforbidList := randomSubsliceOfAddr(r, ForbiddenAddrs)

		msg := asset.NewMsgUnForbidAddr(symbol, owner, unforbidList)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyForbiddenAddr(ctx, k, msg.Symbol, msg.Addresses, false)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token unforbid addr falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func SimulateMsgModifyTokenInfo(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		_, symbol, owner := getOrGenSymbolOwner(r, accounts, ctx, k)

		url := fmt.Sprintf("www.%s.org", symbol)
		describe := fmt.Sprintf("simulation modify info %s", symbol)
		identity := types.TestIdentityString
		msg := asset.NewMsgModifyTokenInfo(symbol, url, describe, identity, owner,
			types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
			types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
			types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
		)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		handler := asset.NewHandler(k)
		ok := simulation2.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyModifyTokenInfo(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("token info modify falied")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func verifyModifyTokenInfo(ctx sdk.Context, k asset.Keeper, msg types.MsgModifyTokenInfo) bool {
	token := k.GetToken(ctx, msg.Symbol)
	return token.GetURL() == msg.URL &&
		token.GetDescription() == msg.Description
}
func RandStringBytes(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func randomTokenAttrs(r *rand.Rand) (name string, symbol string, url string, descrip string, identify string) {
	symbol = RandStringBytes(r, SymbolLenth)
	name = fmt.Sprintf("simulation-%s", symbol)
	url = fmt.Sprintf("www.%s.com", symbol)
	descrip = fmt.Sprintf("simulation issue token : %s", symbol)
	identify = fmt.Sprintf("%d-simulation", r.Int())
	return
}

func randomIssuance(r *rand.Rand) sdk.Int {
	total := 100000000000000000
	issue := r.Intn(100000000000000000)
	if issue < total/2 {
		issue += total / 4
	}
	return sdk.NewInt(int64(issue))
}
