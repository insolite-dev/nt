// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// settingsCommand is a command that used to manage settings of application.
var settingsCommand = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"config"},
	Short:   "Manage settings of notya",
	Run:     runSettingsCommand,
}

// editSettingsCommand is a sub-command of settingsCommand.
// That used as main way for editing notya settings.
var editSettingsCommand = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"-e"},
	Short:   "Edit settings of notya",
	Run:     runEditSettingsCommand,
}

// viewSettingsCommand is a sub-command of settingsCommand.
// Which that used to open settings file via editor.
var viewSettingsCommand = &cobra.Command{
	Use:     "view",
	Aliases: []string{"-v"},
	Short:   "View settings file of notya",
	Run:     runViewSettingsCommand,
}

// initSettingsCommand adds settingsCommand to main application command.
func initSettingsCommand() {
	settingsCommand.AddCommand(editSettingsCommand)
	settingsCommand.AddCommand(viewSettingsCommand)

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

// runEditSettingsCommand runs appropriate service functionalities
// to edit the configuration file by best way.
func runEditSettingsCommand(cmd *cobra.Command, args []string) {
	determineService()

	loading.Start()
	settings, err := service.Settings(nil)
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	editedSettings := models.Settings{}
	if err := survey.Ask(
		assets.SettingsEditPromptQuestions(*settings), &editedSettings,
	); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Breakdown function, if have no changes.
	if !models.IsUpdated(*settings, editedSettings) {
		pkg.Alert(pkg.InfoL, "No changes")
		return
	}

	loading.Start()
	writeErr := service.WriteSettings(editedSettings)
	loading.Stop()

	if writeErr != nil {
		pkg.Alert(pkg.ErrorL, writeErr.Error())
		return
	}

	// Finish process, if notes path not updated.
	if !models.IsPathUpdated(*settings, editedSettings, service.Type()) {
		return
	}

	// Ask to move notes, in case of notes-path updating.
	var moveNotes bool
	survey.AskOne(assets.MoveNotesPrompt, &moveNotes)

	if moveNotes {
		loading.Start()
		err := service.MoveNotes(editedSettings)
		loading.Stop()

		if err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		}
	}
}

// runViewSettingsCommand runs appropriate service functionalities
// to open settings file(json) with CURRENT editor.
func runViewSettingsCommand(cmd *cobra.Command, args []string) {
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
	if models.IsPathUpdated(*beforeSettings, *afterSettings, service.Type()) {
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
