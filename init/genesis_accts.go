package init

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/app"
)

type accountInfo struct {
	addr         sdk.AccAddress
	coins        sdk.Coins
	vestingAmt   sdk.Coins
	vestingStart int64
	vestingEnd   int64
}

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			accInfo, err := collectAccInfo(args)
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

			appState, err = addGenesisAccount(appState, accInfo)
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
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")

	return cmd
}

func collectAccInfo(args []string) (*accountInfo, error) {
	addr, err := getAddress(args[0])
	if err != nil {
		return nil, err
	}

	coins, err := sdk.ParseCoins(args[1])
	if err != nil {
		return nil, err
	}

	vestingStart := viper.GetInt64(flagVestingStart)
	vestingEnd := viper.GetInt64(flagVestingEnd)
	vestingAmt, err := sdk.ParseCoins(viper.GetString(flagVestingAmt))
	if err != nil {
		return nil, err
	}

	accInfo := &accountInfo{
		addr:         addr,
		coins:        coins,
		vestingAmt:   vestingAmt,
		vestingStart: vestingStart,
		vestingEnd:   vestingEnd}
	return accInfo, nil
}

func getAddress(addrOrKeyName string) (addr sdk.AccAddress, err error) {
	addr, err = sdk.AccAddressFromBech32(addrOrKeyName)
	if err != nil {
		kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
		if err != nil {
			return nil, err
		}

		info, err := kb.Get(addrOrKeyName)
		if err != nil {
			return nil, err
		}

		addr = info.GetAddress()
	}
	return
}

func addGenesisAccount(appState app.GenesisState, accInfo *accountInfo) (app.GenesisState, error) {
	for _, stateAcc := range appState.Accounts {
		if stateAcc.Address.Equals(accInfo.addr) {
			return appState, fmt.Errorf("the application state already contains account %v", accInfo.addr)
		}
	}

	acc, err := newGenesisAccount(accInfo)
	if err != nil {
		return appState, err
	}

	appState.Accounts = append(appState.Accounts, acc)
	return appState, nil
}

func newGenesisAccount(accInfo *accountInfo) (genAcc app.GenesisAccount, err error) {
	acc := auth.NewBaseAccountWithAddress(accInfo.addr)
	acc.Coins = accInfo.coins

	if !accInfo.vestingAmt.IsZero() {
		var vacc auth.VestingAccount

		bvacc := &auth.BaseVestingAccount{
			BaseAccount:     &acc,
			OriginalVesting: accInfo.vestingAmt,
			EndTime:         accInfo.vestingEnd,
		}

		if bvacc.OriginalVesting.IsAllGT(acc.Coins) {
			return genAcc, fmt.Errorf("vesting amount cannot be greater than total amount")
		}
		if accInfo.vestingStart >= accInfo.vestingEnd {
			return genAcc, fmt.Errorf("vesting start time must before end time")
		}

		if accInfo.vestingStart != 0 {
			vacc = &auth.ContinuousVestingAccount{
				BaseVestingAccount: bvacc,
				StartTime:          accInfo.vestingStart,
			}
		} else {
			vacc = &auth.DelayedVestingAccount{
				BaseVestingAccount: bvacc,
			}
		}

		return app.NewGenesisAccountI(vacc), nil
	}

	return app.NewGenesisAccount(&acc), nil
}
