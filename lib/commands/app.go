// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// service, is the default service of commands.
var service services.ServiceRepo = &services.LocalService{}

// version is current version of application.
const version = "v1.0.0"

// AppCommand is the root command of application and genesis of all sub-commands.
var appCommand = &cobra.Command{
	Use:     "notya",
	Version: version,
	Short:   "\n 📝 Take notes quickly and expeditiously from terminal.",
}

// initCommands sets all special commands to application command.
func initCommands() {
	initCreateCommand()
	initSetupCommand()
}

// RunApp sets all special commands, then executes app command.
func ExecuteApp() {
	initCommands()

	// Check initialization status of notya,
	// Setup working directories, if it's not initialized before.
	err := initializeIfNotExists()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	_ = appCommand.Execute()
}
