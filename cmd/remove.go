package cmd

import (
	"bode.fun/2fa/core"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "rm identifier",
		Short: "Remove an OTP token from your collection",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			identifier := args[0]

			err := app.DB().Delete([]byte(identifier))
			if err != nil {
				return err
			}

			// TODO: Only print if there actually is a key with this id
			app.Logger().Info("Removed OTP token", "id", identifier)

			err = app.DB().Reset()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return command
}
