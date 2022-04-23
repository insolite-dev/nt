// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"os"

	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

var (
	// stdargs is the global std arguments-state of application.
	stdargs models.StdArgs = models.StdArgs{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
)

var (
	service      services.ServiceRepo // default/active service of all commands.
	localService services.ServiceRepo // default/main service.
	fireService  services.ServiceRepo // firebase integrated service.
)

// appCommand is the root command of application and genesis of all sub-commands.
var appCommand = &cobra.Command{
	Use:     "notya",
	Version: pkg.Version,
	Long: assets.GenerateBanner(
		assets.MinimalisticBanner,
		assets.ShortSlog,
	),
}

// initCommands initializes all sub-commands of application.
func initCommands() {
	initSetupCommand()
	initSettingsCommand()
	initCreateCommand()
	initMkdirCommand()
	initRemoveCommand()
	initViewCommand()
	initEditCommand()
	initRenameCommand()
	initListCommand()
	initCopyCommand()
}

// ExecuteApp is a main function that app starts executing and working.
// Initializes all sub-commands and service for them.
//
// Usually used in [cmd/app.go].
func ExecuteApp() {
	initCommands()

	localService = services.NewLocalService(stdargs)
	if err := localService.Init(); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	config := localService.StateConfig()
	if config.IsFirebaseEnabled() {
		fireService = services.NewFirebaseService(stdargs, localService)
		if err := fireService.Init(); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}
	}

	service = fireService // FIXME: overwrite to local.

	_ = appCommand.Execute()
}
