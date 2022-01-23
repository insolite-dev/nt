// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
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
}

// initSettingsCommand adds settingsCommand to main application command.
func initSettingsCommand() {
	settingsCommand.AddCommand(applyCommand)
	settingsCommand.AddCommand(viewSettingsCommand)

	appCommand.AddCommand(settingsCommand)
}

// runSettingsCommand runs appropriate service functionalities to manage settings.
func runSettingsCommand(cmd *cobra.Command, args []string) {
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
