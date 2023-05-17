package cmd

import (
	"fmt"

	"bode.fun/2fa/core"
	"bode.fun/otp/totp"
	"github.com/spf13/cobra"
)

func NewAddCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "add issuer account base32-secret",
		Aliases: []string{"a"},
		Short:   "Add a new account to your collection",
		Args:    cobra.MatchAll(cobra.ExactArgs(3)),
		RunE: func(cmd *cobra.Command, args []string) error {
			digits, err := cmd.Flags().GetUint("digits")
			if err != nil {
				return err
			}

			if !(digits >= 6 && digits <= 8) {
				return fmt.Errorf("the digit flag has to be between 6 and 8")
			}

			issuer := args[0]
			account := args[1]
			secret := args[2]

			otpions := []totp.TotpOption{
				totp.WithIssuer(issuer),
				totp.WithAccount(account),
				totp.WithDigits(digits),
			}

			totpInstance, err := totp.NewFromBase32(secret, otpions...)
			if err != nil {
				return err
			}

			otpUrl := totpInstance.ToUrl()

			err = app.DB().Set([]byte(totpInstance.Label()), []byte(otpUrl))
			if err != nil {
				return err
			}

			app.Logger().Info(
				"successfully added totp.",
				"account", totpInstance.Account(),
				"issuer",
				totpInstance.Issuer(),
				"id",
				totpInstance.Label(),
			)

			err = app.DB().Sync()
			if err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().UintP("digits", "d", 6, `The amount of digits your code should have.
You can pick between 6, 7 or 8 digits.`)

	return command
}
