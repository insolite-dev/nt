package cmd

import (
	"github.com/spf13/cobra"
)

// Version stores current version of app.
const Version = "v1.0.0"

// AppCommand is the root command of application and genesis of all sub-commands.
var AppCommand = &cobra.Command{
	Use:     "notya",
	Version: Version,
	Short:   "Take notes quickly and expeditiously from terminal.",
	Long: `
	
	Call > notya -h/help
	To read help text(documentation) of Notya CLI.
	`,
}

// setSubCommands, inits all sub commands of app to AppCommand.
func setSubCommands() {
	// TOOD: Add command inits.
}

// RunApp executes appCommand, (sets all sub commands and flags of it).
// It'd be happend only once, on starting program in [main.go].
func RunApp() {
	_ = AppCommand.Execute()
}
