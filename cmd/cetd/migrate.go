package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/market"
	"github.com/coinexchain/dex/app"
)

const (
	flagOutput         = "output"
	flagListValidators = "list-validators"
	GenesisBlockHeight = "genesis-block-height"
	flagGenesisTime    = "genesis-time"
)

func migrateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [from]",
		Short: "Migrate genesis.json (coinexdex -> coinexdex2)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]
			outputFile := viper.GetString(flagOutput)
			return migrateGenesisFile(cdc, inputFile, outputFile)
		},
	}

	cmd.Flags().Int64(GenesisBlockHeight, 0, "node's genesis block height")
	cmd.Flags().Int64(flagGenesisTime, 0, "The unix timestamp for genesis time, in seconds")
	cmd.Flags().String(flagOutput, "", "New genesis.json file")
	cmd.Flags().Bool(flagListValidators, false, "List validators in genesis.json file")

	cmd.MarkFlagRequired(flagGenesisTime)
	return cmd
}

func migrateGenesisFile(cdc *codec.Codec, inputFile, outputFile string) error {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	genDoc := &tm.GenesisDoc{}
	cdc.MustUnmarshalJSON(data, genDoc)
	if viper.GetBool(flagListValidators) {
		listValidators(genDoc)
		return nil
	}
	genesisTime := viper.GetInt64(flagGenesisTime)

	genState := &app.GenesisState{}
	cdc.MustUnmarshalJSON(genDoc.AppState, genState)

	upgradeGenesisState(genState)

	genDoc.ChainID = "coinexdex2"
	genDoc.GenesisBlockHeight = viper.GetInt64(GenesisBlockHeight)
	genDoc.GenesisTime = time.Unix(genesisTime, 0)
	genDoc.AppState = cdc.MustMarshalJSON(genState)
	data = cdc.MustMarshalJSON(genDoc)

	if outputFile == "" {
		fmt.Println(string(data))
		return nil
	}
	return ioutil.WriteFile(outputFile, data, 0644)
}

func upgradeGenesisState(genState *app.GenesisState) {
	genState.GovData.VotingParams.VotingPeriod = app.VotingPeriod
	genState.StakingXData.Params.MinSelfDelegation = app.MinSelfDelegation
	genState.AuthXData.Params = authx.DefaultParams()
	genState.AssetData.Params = asset.DefaultParams()
	genState.MarketData.Params = market.DefaultParams()
	for _, v := range genState.MarketData.Orders {
		if v.FrozenFee != 0 {
			v.FrozenCommission = v.FrozenFee
			v.FrozenFee = 0
		}
	}
	for k, v := range genState.BancorData.BancorInfoMap {
		if v.AR == 0 {
			v.MaxMoney = sdk.ZeroInt()
			genState.BancorData.BancorInfoMap[k] = v
		}
	}
	genState.Incentive.State.HeightAdjustment = 0
	// TODO: more upgrades
}
