package main

import (
	"os"

	"bode.fun/otp/cmd"
	"bode.fun/otp/core"
	"bode.fun/otp/log"
	"github.com/pocketbase/dbx"
	"github.com/spf13/cobra"
)

var (
	Version = "(dev)"
	AppName = "2fa"
)

func main() {
	app := New()
	app.MustRun()
}

type App struct {
	rootCmd *cobra.Command
	logger  log.Logger
	db      *dbx.DB
}

func New() *App {
	rootCmd := &cobra.Command{
		Use:           AppName,
		Version:       Version,
		Short:         "Manage your otp tokens securely from your command line.",
		SilenceErrors: true,
		SilenceUsage:  true, // Invert this to enable usage printout when an error occurs
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	logger := log.New(os.Stderr, AppName)

	// TODO: load this from config and maybe close the db
	db, err := core.ConnectDB(":memory:")
	if err != nil {
		logger.Panic(err)
	}

	app := &App{
		rootCmd,
		logger,
		db,
	}

	app.registerCommands()

	return app
}

func (a *App) registerCommands() {
	a.rootCmd.AddCommand(cmd.AddCmd)
}

func (a *App) Run() error {
	return a.rootCmd.Execute()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.logger.Fatal(err)
	}
}
