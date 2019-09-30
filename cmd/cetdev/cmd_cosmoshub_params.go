package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/app"
)

type genesisDoc struct {
	AppState app.GenesisState `json:"app_state"`
}

func CosmosHubParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cosmos-hub-params",
		Short: "Print default params",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "https://raw.githubusercontent.com/cosmos/launch/master/genesis.json"
			fmt.Printf("downloading %s ...\n", url)

			resp, err := http.Get(url)
			if err != nil {
				return err
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			body = fixAddresses(body)
			genDoc := genesisDoc{}
			err = cdc.UnmarshalJSON(body, &genDoc)
			if err != nil {
				return err
			}

			printParams(genDoc.AppState)
			return nil
		},
	}
	cmd.Flags().Bool("include-sdk", false, "include params defined by cosmos-sdk modules")
	return cmd
}

func fixAddresses(body []byte) []byte {
	body2 := regexp.MustCompile(`"cosmosvalconspub[^"]*"`).ReplaceAllString(string(body), `"coinexvalconspub1addwnpepqtx9hr0sqk778yhdchdnzt6sfdqm3leg6x9yfjclnjc2g6eczrv75y8mcn5"`)
	body2 = regexp.MustCompile(`"cosmosvalcons[^"]*"`).ReplaceAllString(body2, `"`+sdk.ConsAddress("").String()+`"`)
	body2 = regexp.MustCompile(`"cosmosvaloper[^"]*"`).ReplaceAllString(body2, `"coinexvaloper1dj2m0nmwp7khdnltzmtfqexasx69hg5q385rlu"`)
	body2 = regexp.MustCompile(`"cosmos[^"]*"`).ReplaceAllString(body2, `"coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97"`)
	return []byte(body2)
}

func printParams(genState interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"package", "param name", "param value"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
	})

	gs := reflect.ValueOf(genState)
	for i := 0; i < gs.NumField(); i++ {
		gsFieldVal := gs.Field(i)
		if gsFieldVal.Type().Kind() == reflect.Struct {
			collectModuleParams(gsFieldVal, table)
		}
	}

	table.Render()
}

func collectModuleParams(modGenState reflect.Value, table *tablewriter.Table) {
	// fmt.Println(modGenState.Type().Name(), modGenState.Type().PkgPath())
	pkg := modGenState.Type().PkgPath()
	for i := 0; i < modGenState.NumField(); i++ {
		filedVal := modGenState.Field(i)
		if strings.HasSuffix(filedVal.Type().Name(), "Params") {
			collectModuleParamValues(pkg, filedVal, table)
		}
	}
}

func collectModuleParamValues(pkgName string, modParams reflect.Value, table *tablewriter.Table) {
	for i := 0; i < modParams.NumField(); i++ {
		// fmt.Println(pkg, modParams.Type().Field(i).Name, modParams.Field(i))
		paramName := modParams.Type().Field(i).Name
		paramVal := fmt.Sprintf("%v", modParams.Field(i))
		table.Append([]string{pkgName, paramName, paramVal})
	}
}
