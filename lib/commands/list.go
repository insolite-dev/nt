// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// listCommand, is a application command, which used to list all
// notes/files from the notya folder.
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

// runListCommand runs appropriate service functionalities
// to list all notes from the notya folder.
func runListCommand(cmd *cobra.Command, args []string) {
	// Generate a list of notes.
	list, err := pkg.ListDir(NotyaPath)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.ShowListOfNotes(list, 3)
}
