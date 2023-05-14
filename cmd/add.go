package cmd

import (
	"fmt"

	"bode.fun/otp/core"
	"github.com/spf13/cobra"
)

// TODO: Make this into a function which receives an app instance like https://github.com/pocketbase/pocketbase/blob/c6d599244239ed17b2f2f7ce892b1279ddabf5ac/cmd/serve.go#L27

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
