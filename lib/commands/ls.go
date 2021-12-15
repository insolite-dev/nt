// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/spf13/cobra"
)

// lsCommand, is a application command, which used to list all
// notes/files from the notya folder.
var lsCommand = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all notya notes",
	Run:     runLsCommand,
}

// initLsCommand adds lsCommand to main application command.
func initLsCommand() {
	appCommand.AddCommand(lsCommand)
}

// runLsCommand runs appropriate service functionalities
// to list all notes from the notya folder.
func runLsCommand(cmd *cobra.Command, args []string) {
	// TODO: Add functionality.
}
