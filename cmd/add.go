package cmd

import (
	"fmt"

	"bode.fun/2fa/core"
	"github.com/spf13/cobra"
)

func NewAddCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new account to your collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("add is not implemented yet")
		},
	}
}
