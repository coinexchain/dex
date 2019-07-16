package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/coinexchain/dex/modules/asset"
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
	flagTokenURL,
	flagTokenDescription,
}

// AddGenesisTokenCmd returns add-genesis-token cobra Command.
func AddGenesisTokenCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
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
	--url="www.coinex.org" \
	--description="A public chain built for the decentralized exchange"
`),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			token, err := parseTokenInfo()
			if err != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add genesis account to the app state
			var genesisState asset.GenesisState

			cdc.MustUnmarshalJSON(appState[asset.ModuleName], &genesisState)

			err = addGenesisToken(&genesisState, token)
			if err != nil {
				return err
			}

			genesisStateBz := cdc.MustMarshalJSON(genesisState)
			appState[asset.ModuleName] = genesisStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flagName, "", "token name is limited to 32 unicode characters")
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
	cmd.Flags().String(flagTokenURL, "", "url of token website")
	cmd.Flags().String(flagTokenDescription, "", "description of token info")

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

func getAddress(addrOrKeyName string) (addr sdk.AccAddress, err error) {
	addr, err = sdk.AccAddressFromBech32(addrOrKeyName)
	if err != nil {
		return getAddressFromKeyBase(addrOrKeyName)
	}
	return
}

func getAddressFromKeyBase(keyName string) (sdk.AccAddress, error) {
	kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
	if err != nil {
		return nil, err
	}

	info, err := kb.Get(keyName)
	if err != nil {
		return nil, err
	}

	addr := info.GetAddress()
	return addr, nil
}

func addGenesisToken(genesisState *asset.GenesisState, token asset.Token) error {
	for _, t := range genesisState.Tokens {
		if token.GetSymbol() == t.GetSymbol() {
			return fmt.Errorf("the application state already contains token %s", token.GetSymbol())
		}
	}
	genesisState.Tokens = append(genesisState.Tokens, token)

	return nil
}
