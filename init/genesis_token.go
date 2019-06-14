package init

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/cosmos/cosmos-sdk/client"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
)

const (
	flagName             = "name"
	flagSymbol           = "symbol"
	flagOwner            = "owner"
	flagTotalSupply      = "total-supply"
	flagMintable         = "mintable"
	flagBurnable         = "burnable"
	flagAddrForbiddable  = "addr-forbiddable"
	flagTokenForbiddable = "token-forbiddable"
	flagTotalBurn        = "total-burn"
	flagTotalMint        = "total-mint"
	flagIsForbidden      = "is-forbidden"
)

var tokenFlags = []string{
	flagName,
	flagSymbol,
	flagOwner,
	flagTotalSupply,
	flagMintable,
	flagBurnable,
	flagAddrForbiddable,
	flagTokenForbiddable,
	flagTotalBurn,
	flagTotalMint,
	flagIsForbidden,
}

// AddGenesisTokenCmd returns add-genesis-token cobra Command.
func AddGenesisTokenCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-token",
		Short: "Add genesis token to genesis.json",
		Long: strings.TrimSpace(
			`
Example:
$ cetd add-genesis-token --name="CoinEx Chain Native Token" \
	--symbol="cet" \
	--owner=ownerkey \
	--total-supply=588788547005740000 \
	--mintable=false \
	--burnable=true \
	--addr-forbiddable=false \
	--token-forbiddable=false \
	--total-burn=411211452994260000 \
	--total-mint=0 \
	--is-forbidden=false \
`),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			token, err := parseTokenInfo()
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `cetd init` first", genFile)
			}

			genDoc, err := LoadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			appState, err = addGenesisToken(appState, token)
			if err != nil {
				return err
			}

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			return ExportGenesisFile(genFile, genDoc.ChainID, nil, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(flagName, "", "tToken name is limited to 32 unicode characters")
	cmd.Flags().String(flagSymbol, "", "token symbol is limited to [a-z][a-z0-9]{1,7}")
	cmd.Flags().String(flagOwner, "", "token owner")
	cmd.Flags().Int64(flagTotalSupply, 0, "the total supply for token can have a maximum of "+
		"8 digits of decimal and is boosted by 1e8 in order to store as int64. "+
		"The amount before boosting should not exceed 90 billion.")
	cmd.Flags().Bool(flagMintable, false, "whether the token could be minted")
	cmd.Flags().Bool(flagBurnable, true, "whether hte token could be burned")
	cmd.Flags().Bool(flagAddrForbiddable, false, "whether the token holder address can be forbidden by token owner")
	cmd.Flags().Bool(flagTokenForbiddable, false, "whether the token can be forbidden")
	cmd.Flags().Int64(flagTotalBurn, 0, "the total burn amount")
	cmd.Flags().Int64(flagTotalMint, 0, "the total mint amount")
	cmd.Flags().Bool(flagIsForbidden, false, "whether the token is forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range tokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

func parseTokenInfo() (asset.Token, error) {
	token := &asset.BaseToken{}
	var err error

	owner, err := getAddress(viper.GetString(flagOwner))
	if err != nil {
		return nil, err
	}

	if err = token.SetName(viper.GetString(flagName)); err != nil {
		return nil, err
	}
	if err = token.SetSymbol(viper.GetString(flagSymbol)); err != nil {
		return nil, err
	}
	if err = token.SetOwner(owner); err != nil {
		return nil, err
	}
	if err = token.SetTotalSupply(viper.GetInt64(flagTotalSupply)); err != nil {
		return nil, err
	}
	if err = token.SetTotalBurn(viper.GetInt64(flagTotalBurn)); err != nil {
		return nil, err
	}
	if err = token.SetTotalMint(viper.GetInt64(flagTotalMint)); err != nil {
		return nil, err
	}

	token.SetMintable(viper.GetBool(flagMintable))
	token.SetBurnable(viper.GetBool(flagBurnable))
	token.SetAddrForbiddable(viper.GetBool(flagAddrForbiddable))
	token.SetTokenForbiddable(viper.GetBool(flagTokenForbiddable))
	token.SetIsForbidden(viper.GetBool(flagTokenForbiddable))

	return token, nil
}

func addGenesisToken(appState app.GenesisState, token asset.Token) (app.GenesisState, error) {
	for _, t := range appState.AssetData.Tokens {
		if token.GetSymbol() == t.GetSymbol() {
			return appState, fmt.Errorf("the application state already contains token %s", token.GetSymbol())
		}
	}
	appState.AssetData.Tokens = append(appState.AssetData.Tokens, token)

	return appState, nil
}
