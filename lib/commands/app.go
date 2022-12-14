//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"os"

	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

var (
	// Main spin animator of application.
	loading = pkg.Spinner()

	// stdargs is the global std arguments-state of application.
	stdargs models.StdArgs = models.StdArgs{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
)

var (
	service      services.ServiceRepo // default/active service of all commands.
	localService services.ServiceRepo // default/main service.
	fireService  services.ServiceRepo // firebase integrated service.
)

// serviceFromType returns type appropriate service instance.
func serviceFromType(t string, enable bool) services.ServiceRepo {
	switch t {
	case services.LOCAL.ToStr():
		return localService
	case services.FIRE.ToStr():
		if enable {
			setupFirebaseService()
		}
		return fireService
	}

	return service
}

// Decides whether use firebase service as main service or not.
var firebaseF bool

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
	appCommand.PersistentFlags().BoolVarP(
		&firebaseF, "firebase", "f", false,
		"Run commands base on firebase service",
	)

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
	initFetchCommand()
	initPushCommand()
	initMigrateCommand()
	initCutCommand()
	initRemoteCommand()
}

// ExecuteApp is a main function that app starts executing and working.
// Initializes all sub-commands and service for them.
//
// Usually used in [cmd/app.go].
func ExecuteApp() {
	loading.Start()

	initCommands()

	setupLocalService()
	service = localService

	_ = appCommand.Execute()
}

// determineService checks user input service after execution main command.
// if user has provided a custom service for specific command-execution, it updates
// the [service] value with that custom-service[fireService ... etc].
func determineService() {
	if !firebaseF {
		return
	}

	setupFirebaseService()
	service = fireService

	//
	// TODO: implement other services.
	//
}

// setupLocalService initializes the local service.
// makes it able at [localService] instance.
func setupLocalService() {
	loading.Start()

	localService = services.NewLocalService(stdargs)
	err := localService.Init(nil)

	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		os.Exit(1)
	}
}

// setupFirebaseService initializes the firebase service.
// makes it able at [fireService] instance.
func setupFirebaseService() {
	loading.Start()

	fireService = services.NewFirebaseService(stdargs, localService)
	err := fireService.Init(nil)

	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		os.Exit(1)
	}
}
