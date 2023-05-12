package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a new account to your collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("add is not implemented yet")
	},
}
