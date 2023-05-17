package cmd

import (
	"bode.fun/2fa/core"
	"github.com/spf13/cobra"
)

func NewSyncCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "sync",
		Short: "Sync your collection with the encrypted cloud database",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.DB().Sync()
			if err != nil {
				return err
			}

			app.Logger().Info("Database synced successfully")

			return nil
		},
	}

	return command
}
