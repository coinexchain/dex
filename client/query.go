package client

import (
	"errors"
	"fmt"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetAccount(ctx context.CLIContext, address []byte) (authx.AccountX, error) {
	if ctx.AccDecoder == nil {
		return authx.AccountX{}, errors.New("account decoder required but not provided")
	}

	res, err := QueryAccountx(ctx, address)
	if err != nil {
		return authx.AccountX{}, err
	}

	var accountx authx.AccountX
	if err := ctx.Codec.UnmarshalJSON(res, &accountx); err != nil {
		return authx.AccountX{}, err
	}

	return accountx, nil
}

func EnsureAccountExistsFromAddr(ctx context.CLIContext, addr sdk.AccAddress) error {
	_, err := QueryAccountx(ctx, addr)
	return err
}

func QueryAccountx(ctx context.CLIContext, addr sdk.AccAddress) ([]byte, error) {
	bz, err := ctx.Codec.MarshalJSON(authx.NewQueryAccountxParams(addr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", authx.StoreKey, authx.QueryAccountx)

	res, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}
