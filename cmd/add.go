package cmd

import (
	"fmt"

	"bode.fun/2fa/core"
	"github.com/spf13/cobra"
)

func NewAddCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new account to your collection",
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		ArgAliases: []string{
			"key",
			"base32-secret",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("add is not implemented yet")
		},
	}

	command.Flags().UintP("digits", "d", 6, `The amount of digits your code should have.
You can pick between 6, 7 or 8 digits.`)

	return command
}
