//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

// settingsCommand is a general command that used to manage settings of application.
// implements a sub command to edit the settings file.
//
// Default functionality of running settingsCommand is just printing current settings data.
// to edit it should include running editSettingsCommand.
var settingsCommand = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"config"},
	Short:   "Manage settings of notya",
	Run:     runSettingsCommand,
}

// editSettingsCommand is a sub-command of settingsCommand.
// that opens the setting (configuration) file with your current editor.
var editSettingsCommand = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"-e"},
	Short:   "Opens the configuration file of notya with your current editor",
	Run:     runEditSettingsCommand,
}

// initSettingsCommand adds settingsCommand to main application command.
func initSettingsCommand() {
	settingsCommand.AddCommand(editSettingsCommand)

	appCommand.AddCommand(settingsCommand)
}

// runSettingsCommand runs appropriate service functionalities to manage settings.
func runSettingsCommand(cmd *cobra.Command, args []string) {
	determineService()

	loading.Start()
	settings, err := service.Settings(nil)
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Print settings' current values.
	pkg.PrintSettings(*settings)
	pkg.Print("\n > [notya settings -h/help] for more", color.FgGreen)
}

// runViewSettingsCommand runs appropriate service functionalities
// to open settings file(json) with CURRENT editor.
func runEditSettingsCommand(cmd *cobra.Command, args []string) {
	determineService()

	loading.Start()
	beforeSettings, err := service.Settings(nil)
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	openErr := service.OpenSettings(*beforeSettings)
	if openErr != nil {
		pkg.Alert(pkg.ErrorL, openErr.Error())
		return
	}

	loading.Start()
	afterSettings, err := service.Settings(&beforeSettings.ID)
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask to move notes if path were updated.
	if pkg.IsPathUpdated(*beforeSettings, *afterSettings, service.Type()) {
		var moveNotes bool
		if survey.AskOne(assets.MoveNotesPrompt, &moveNotes); !moveNotes {
			return
		}

		loading.Start()
		err := service.MoveNotes(*afterSettings)
		loading.Stop()

		if err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		}
	}
}
