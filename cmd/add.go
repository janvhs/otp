package cmd

import (
	"fmt"
	"strings"

	"bode.fun/2fa/core"
	"bode.fun/otp/totp"
	"github.com/spf13/cobra"
)

func NewAddCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "add account base32-secret",
		Aliases: []string{"a"},
		Short:   "Add a new account to your collection",
		Args:    cobra.MatchAll(cobra.MinimumNArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			digits, err := cmd.Flags().GetUint("digits")
			if err != nil {
				return err
			}

			if !(digits >= 6 && digits <= 8) {
				return fmt.Errorf("the code has to be between 6 and 8 digits")
			}

			// TODO: Support issuer and account properly
			app.Logger().Info("Account and Issuer are not supported yet")

			argLen := len(args)

			account := args[0]
			secret := args[1]

			if argLen > 2 {
				accountPats := args[:argLen-1]
				account = strings.Join(accountPats, " ")
				secret = args[argLen-1]
			}

			otpions := []totp.TotpOption{
				totp.WithAccount(account),
				totp.WithDigits(digits),
			}

			totpInstance, err := totp.NewFromBase32(secret, otpions...)
			if err != nil {
				return err
			}

			otpUrl := totpInstance.ToUrl()
			err = app.DB().Set([]byte(account), []byte(otpUrl))
			if err != nil {
				return err
			}

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
