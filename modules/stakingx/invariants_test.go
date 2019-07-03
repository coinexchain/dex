package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
)

type mockAssetKeeper struct {
	tokens []asset.BaseToken
}

func (k mockAssetKeeper) GetToken(ctx sdk.Context, symbol string) asset.Token {
	for _, token := range k.tokens {
		if token.Symbol == symbol {
			return &token
		}
	}
	return nil
}

func TestTotalSupplyInvariants(t *testing.T) {
	//intialize keeper & params
	defaultToken := asset.BaseToken{
		Symbol:      "cet",
		TotalSupply: 100e8,
	}
	ak := mockAssetKeeper{
		tokens: []asset.BaseToken{defaultToken},
	}
	sxk, ctx, _ := setUpInput()
	sxk.SetParams(ctx, DefaultParams())

	pool := staking.Pool{
		NotBondedTokens: sdk.NewInt(10e8),
		BondedTokens:    sdk.NewInt(90e8),
	}
	sxk.sk.SetPool(ctx, pool)

	//test TotalSupplyInvariants
	invariantFc := TotalSupplyInvariants(sxk, ak)
	require.Nil(t, invariantFc(ctx))
}

func TestTotalSupplyInvariantsFail(t *testing.T) {
	//intialize keeper & params
	defaultToken := asset.BaseToken{
		Symbol:      "cet",
		TotalSupply: 200e8,
	}
	ak := mockAssetKeeper{
		tokens: []asset.BaseToken{defaultToken},
	}
	sxk, ctx, _ := setUpInput()
	sxk.SetParams(ctx, DefaultParams())

	pool := staking.Pool{
		NotBondedTokens: sdk.NewInt(10e8),
		BondedTokens:    sdk.NewInt(90e8),
	}
	sxk.sk.SetPool(ctx, pool)

	//test TotalSupplyInvariants
	invariantFc := TotalSupplyInvariants(sxk, ak)
	require.NotNil(t, invariantFc(ctx))
}

func TestTotalSupplyInvariantsNil(t *testing.T) {
	//intialize keeper & params
	ak := mockAssetKeeper{
		tokens: []asset.BaseToken{},
	}
	sxk, ctx, _ := setUpInput()
	sxk.SetParams(ctx, DefaultParams())

	pool := staking.Pool{
		NotBondedTokens: sdk.NewInt(10e8),
		BondedTokens:    sdk.NewInt(90e8),
	}
	sxk.sk.SetPool(ctx, pool)

	//test TotalSupplyInvariants
	invariantFc := TotalSupplyInvariants(sxk, ak)
	require.NotNil(t, invariantFc(ctx))
}
