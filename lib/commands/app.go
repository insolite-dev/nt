// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"os"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

var (
	// stdargs is the global std state of application.
	stdargs models.StdArgs = models.StdArgs{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
)

// service, is the default service of all commands.
var service services.ServiceRepo

// AppCommand is the root command of application and genesis of all sub-commands.
var appCommand = &cobra.Command{
	Use:     "notya",
	Version: pkg.Version,
	Short:   "\n 📝 Take notes quickly and expeditiously from terminal.",
	Long:    pkg.GetBanner(),
}

// initCommands sets all special commands to application command.
func initCommands() {
	initSetupCommand()
	initCreateCommand()
	initRemoveCommand()
	initViewCommand()
	initEditCommand()
	initRenameCommand()
	initListCommand()
}

// RunApp sets all special commands, checks notya initialization status,
// and then executes main app command.
func ExecuteApp() {
	initCommands()

	// Initialize new local service.
	service = services.NewLocalService(stdargs)

	// Initialize application.
	if err := service.Init(); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	_ = appCommand.Execute()
}
