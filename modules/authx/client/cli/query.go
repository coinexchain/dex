package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx/types"
)

func GetAccountX(ctx context.CLIContext, address []byte) (types.AccountX, error) {
	res, err := QueryAccountX(ctx, address)
	if err != nil {
		return types.AccountX{}, err
	}

	var accountX types.AccountX
	if err := ctx.Codec.UnmarshalJSON(res, &accountX); err != nil {
		return types.AccountX{}, err
	}

	return accountX, nil
}

func QueryAccountX(ctx context.CLIContext, addr sdk.AccAddress) ([]byte, error) {
	bz, err := ctx.Codec.MarshalJSON(auth.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAccountX)

	res, _, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}
