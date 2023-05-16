package cmd

import (
	"fmt"
	"strings"

	"bode.fun/2fa/core"
	"bode.fun/otp/totp"
	"github.com/spf13/cobra"
)

func NewGetCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "get account",
		Aliases: []string{"g"},
		Short:   "Get a totp from your collection",
		Args:    cobra.MatchAll(cobra.MinimumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Support issuer and account properly
			app.Logger().Info("Account and Issuer are not supported yet")

			account := strings.Join(args, " ")

			otpUrlAsBytes, err := app.DB().Get([]byte(account))
			if err != nil {
				return err
			}

			totpInstance, err := totp.NewFromUrl(string(otpUrlAsBytes))
			if err != nil {
				return err
			}

			fmt.Println(totpInstance.Now())

			return nil
		},
	}

	command.Flags().UintP("digits", "d", 6, `The amount of digits your code should have.
You can pick between 6, 7 or 8 digits.`)

	return command
}
