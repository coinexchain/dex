package client

import (
	"fmt"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func GetAccountX(ctx context.CLIContext, address []byte) (authx.AccountX, error) {

	res, err := QueryAccountX(ctx, address)
	if err != nil {
		return authx.AccountX{}, err
	}

	var accountX authx.AccountX
	if err := ctx.Codec.UnmarshalJSON(res, &accountX); err != nil {
		return authx.AccountX{}, err
	}

	return accountX, nil
}

func EnsureAccountExistsFromAddr(ctx context.CLIContext, addr sdk.AccAddress) error {
	_, err := QueryAccountX(ctx, addr)
	return err
}

func QueryAccountX(ctx context.CLIContext, addr sdk.AccAddress) ([]byte, error) {
	bz, err := ctx.Codec.MarshalJSON(auth.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", authx.StoreKey, authx.QueryAccountX)

	res, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}
