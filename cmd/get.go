package cmd

import (
	"fmt"

	"bode.fun/2fa/core"
	"bode.fun/otp/totp"
	"github.com/spf13/cobra"
)

// TODO: This is not a good experience
func NewGetCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "get identifier",
		Aliases: []string{"g"},
		Short:   "Get an OTP code from your collection",
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			identifier := args[0]

			err := app.DB().Reset()
			if err != nil {
				return err
			}

			// TODO: Add prefixed zeros
			// TODO: Add a custom error message when key is not found
			code, err := getOtpCode(app, identifier)
			if err != nil {
				return err
			}
			fmt.Println(code)

			return nil
		},
	}

	return command
}

func getOtpCode(app core.App, identifier string) (uint32, error) {
	var otpCode uint32
	otpUrlAsBytes, err := app.DB().Get([]byte(identifier))
	if err != nil {
		return otpCode, err
	}

	totpInstance, err := totp.NewFromUrl(string(otpUrlAsBytes))
	if err != nil {
		return otpCode, err
	}

	otpCode = totpInstance.Now()
	return otpCode, err
}
