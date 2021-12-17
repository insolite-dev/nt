// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// rmCommand, is a command model which used to remove a note or file.
var rmCommand = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove/Delete a notya file",
	Run:     runRmCommand,
}

// initRmCommand adds rmCommand to main application command.
func initRmCommand() {
	appCommand.AddCommand(rmCommand)
}

// runRmCommand runs appropriate service commands to remove note.
func runRmCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note := models.Note{Title: args[0], Path: NotyaPath + args[0]}

		// Check if file exists or not.
		if !pkg.FileExists(note.Path) {
			notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
			pkg.Alert(pkg.ErrorL, notExists)
			return
		}

		// TODO: Remove note, at [note]'s path.

		return
	}

	// Generate array of all notes' names.
	notes, err := pkg.ListDir(NotyaPath)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	prompt := &survey.Select{Message: "Choose a note to rename:", Options: notes}
	survey.AskOne(prompt, &selected)

	// TODO: Remove note, at [note]'s path.
}
