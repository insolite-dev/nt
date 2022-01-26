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

// editCommand is a command model which used to overwrite body of notes or files.
var editCommand = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"overwrite", "update"},
	Short:   "Edit/Update note data",
	Run:     runEditCommand,
}

// initEditCommand adds editCommand to main application command.
func initEditCommand() {
	appCommand.AddCommand(editCommand)
}

// runEditCommand runs appropriate service commands to edit/overwrite note data.
func runEditCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note := models.Note{Title: args[0]}

		if err := service.Open(note); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		}

		return
	}

	// Generate all note names.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	selected := ""
	prompt := &survey.Select{
		Message: "Choose a note to edit:",
		Options: pkg.MapNotesList(notes),
	}
	survey.AskOne(prompt, &selected)

	// Open selected note-file.
	if err := service.Open(models.Note{Title: selected}); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
}
