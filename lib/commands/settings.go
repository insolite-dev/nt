// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
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

// applyCommand is sub-command of settingsCommand.
// Which that used to apply HAND-MADE changes on settings file.
var applyCommand = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"-a"},
	Short:   "Apply HAND-MADE changes",
	Run:     runApplySettingsCommand,
}

// initSettingsCommand adds settingsCommand to main application command.
func initSettingsCommand() {
	settingsCommand.AddCommand(editSettingsCommand)
	settingsCommand.AddCommand(viewSettingsCommand)
	settingsCommand.AddCommand(applyCommand)

	appCommand.AddCommand(settingsCommand)
}

// runSettingsCommand runs appropriate service functionalities to manage settings.
func runSettingsCommand(cmd *cobra.Command, args []string) {
	_, err := service.Settings()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// TODO: log current settings
}

// runEditSettingsCommand runs appropriate service functionalities
// to edit the configuration file by best way.
func runEditSettingsCommand(cmd *cobra.Command, args []string) {
	settings, err := service.Settings()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	editedSettings := models.Settings{}
	survey.Ask([]*survey.Question{
		{
			Name: "editor",
			Prompt: &survey.Input{
				Default: settings.Editor,
				Message: "Editor",
				Help:    "Editor for notya. --> vi/vim/nvim/code/code-insiders ...",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "local_path",
			Prompt: &survey.Input{
				Default: settings.LocalPath,
				Message: "Local Path",
				Help:    "Local path of notya base working directory",
			},
		},
	},
		&editedSettings,
	)

	// TODO: apply changes, via --> runApplySettingsCommand(cmd, args)
}

// runViewSettingsCommand runs appropriate service functionalities
// to open settings file(json) with CURRENT editor.
func runViewSettingsCommand(cmd *cobra.Command, args []string) {
	settingsFile := models.Note{Title: models.SettingsName}
	if err := service.Open(settingsFile); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
	}

	runApplySettingsCommand(cmd, args)
}

// runApplySettingsCommand runs appropriate service functionalities
// to apply HAND-MADE changes on settings(configuration) file.
func runApplySettingsCommand(cmd *cobra.Command, args []string) {
	// TODO: implement apply functionality.
}
