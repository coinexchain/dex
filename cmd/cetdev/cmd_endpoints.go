package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/lcd"
)

func RestEndpointsCmd(registerRoutesFn func(*lcd.RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-endpoints",
		Short: "Show LCD REST endpoints",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Method", "Path"})

			router := prepareRouter(registerRoutesFn)
			router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
				path, _ := route.GetPathTemplate()
				method := getMethod(route)
				table.Append([]string{method, path})
				return nil
			})

			table.Render()
			return nil
		},
	}
	return cmd
}

func prepareRouter(registerRoutesFn func(*lcd.RestServer)) *mux.Router {
	rs := &lcd.RestServer{
		Mux:    mux.NewRouter(),
		CliCtx: context.CLIContext{},
	}

	registerRoutesFn(rs)
	return rs.Mux
}

func getMethod(route *mux.Route) string {
	methods, _ := route.GetMethods()
	if len(methods) > 0 {
		return methods[0]
	}
	return ""
}
