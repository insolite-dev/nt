// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

// Version stores current version of app.
const Version = "v1.0.0"

// AppCommand is the root command of application and genesis of all sub-commands.
var appCommand = &cobra.Command{
	Use:     "notya",
	Version: Version,
	Short:   "Take notes quickly and expeditiously from terminal.",
	Long: `

	Usage: notya repo <command> [flags]

	Available commands:

	`,
}

// setSubCommands, inits all sub commands of app to AppCommand.
func setSubCommands() {
	// TOOD: Add command inits.
}

// RunApp executes appCommand, (sets all sub commands and flags of it).
// It'd be happend only once, on starting program in [main.go].
func RunApp() {
	_ = appCommand.Execute()
}
