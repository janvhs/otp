package main

import (
	"os"

	"bode.fun/2fa/cmd"
	"github.com/charmbracelet/charm/kv"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	Version = "(dev)"
	AppName = "2fa"
	DBName  = AppName
)

func main() {
	app := New()
	app.MustRun()
}

type App struct {
	rootCmd *cobra.Command
	logger  *log.Logger
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

	logger := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: AppName,
	})

	// TODO: load this from config and maybe close the db
	// TODO: Lazy load this
	db, err := kv.OpenWithDefaults(DBName)
	if err != nil {
		logger.Fatal("can't open the database", "err", err)
	}

	app := &App{
		rootCmd,
		logger,
		db,
	}

	app.registerCommands()

	return app
}

// TODO: Add initialization on access
func (a *App) DB() *kv.KV {
	return a.db
}

func (a *App) Logger() *log.Logger {
	return a.logger
}

func (a *App) registerCommands() {
	a.rootCmd.AddCommand(
		cmd.NewAddCommand(a),
		cmd.NewGetCommand(a),
		cmd.NewListCommand(a),
		cmd.NewRemoveCommand(a),
		cmd.NewSyncCommand(a),
	)
}

func (a *App) Run() error {
	defer a.DB().Close()
	return a.rootCmd.Execute()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.Logger().Fatal(err)
	}
}
