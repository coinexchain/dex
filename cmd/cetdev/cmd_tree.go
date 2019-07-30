package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

func ShowCommandTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command-tree",
		Short: "Show Cetd command tree",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tree := getCmdTree(cmd.Root(), 0)
			fmt.Println(tree)
			return nil
		},
	}
	return cmd
}

func getCmdTree(cmd *cobra.Command, level int) string {
	if cmd == client.LineBreak {
		return ""
	}

	tree := strings.Repeat("  ", level) + cmd.Name() + "\n"
	if len(cmd.Commands()) > 0 {
		for _, subCmd := range cmd.Commands() {
			tree = tree + getCmdTree(subCmd, level+1)
		}
	}

	return tree
}
