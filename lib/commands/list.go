// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// listCommand is a command that used to list all exiting notes.
var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all notya notes",
	Run:     runListCommand,
}

// initListCommand adds listCommand to main application command.
func initListCommand() {
	appCommand.AddCommand(listCommand)
}

// runListCommand runs appropriate service functionalities to log all notes.
func runListCommand(cmd *cobra.Command, args []string) {
	// Generate a list of notes.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.PrintNotes(notes)
}
