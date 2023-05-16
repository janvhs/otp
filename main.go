package main

import (
	"os"

	"bode.fun/2fa/cmd"
	"bode.fun/2fa/log"
	"github.com/charmbracelet/charm/kv"
	"github.com/spf13/cobra"
)

var (
	Version = "(dev)"
	AppName = "2fa"
	DBName  = AppName
)

func main() {
	app := New()
	defer app.DB().Close()
}

type App struct {
	rootCmd *cobra.Command
	logger  log.Logger
	db      *kv.KV
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
	db, err := kv.OpenWithDefaults(DBName)
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

func (a *App) DB() *kv.KV {
	return a.db
}

func (a *App) Logger() log.Logger {
	return a.logger
}

func (a *App) registerCommands() {
	a.rootCmd.AddCommand(cmd.NewAddCommand(a))
}

func (a *App) Run() error {
	defer a.DB().Close()
	return a.rootCmd.Execute()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.logger.Fatal(err)
	}
}
