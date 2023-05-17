package cmd

import (
	"bode.fun/2fa/core"
	"github.com/spf13/cobra"
)

func NewListCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List all OTP code from your collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.DB().Reset()
			if err != nil {
				return err
			}

			identifiers, err := app.DB().Keys()
			if err != nil {
				return err
			}

			identifierCount := len(identifiers)

			if identifierCount == 0 {
				app.Logger().Warn("You haven't stored any OTP tokens, yet")
				return nil
			}

			for _, identifier := range identifiers {
				var code uint32
				code, _ = getOtpCode(app, string(identifier))
				app.Logger().Print(nil, "id", string(identifier), "code", code)
			}

			return nil
		},
	}

	return command
}
