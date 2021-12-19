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

// removeCommand, is a command model which used to remove a note or file.
var removeCommand = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove/Delete a notya file",
	Run:     runRemoveCommand,
}

// initRemoveCommand adds removeCommand to main application command.
func initRemoveCommand() {
	appCommand.AddCommand(removeCommand)
}

// runRemoveCommand runs appropriate service commands to remove note.
func runRemoveCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note := models.Note{Title: args[0]}

		removeAndFinish(note)
		return
	}

	// Generate array of all notes' names.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	prompt := &survey.Select{Message: "Choose a note to remove:", Options: notes}
	survey.AskOne(prompt, &selected)

	removeAndFinish(models.Note{Title: selected})
}

// removeAndFinish removes given note and alerts success message if everything is OK.
func removeAndFinish(note models.Note) {
	if err := service.Remove(note); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Note removed successfully: "+note.Title)
}
