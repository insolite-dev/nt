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

// editCommand, is a command model which used to overwrite body of notes or files.
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

		// Check if file exists or not.
		if !pkg.FileExists(note.Path) {
			notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
			pkg.Alert(pkg.ErrorL, notExists)
			return
		}

		// Open note-file with vi, to edit it.
		openingErr := pkg.OpenFileWithVI(note.Path, StdArgs) // TODO: Pass full path
		if openingErr != nil {
			pkg.Alert(pkg.ErrorL, openingErr.Error())
		}

		pkg.Alert(pkg.SuccessL, "Note updated successfully: "+note.Title)
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
		Options: notes,
	}
	survey.AskOne(prompt, &selected)

	// Open created note-file with vi, to edit it.
	openingErr := pkg.OpenFileWithVI(selected, StdArgs) // TODO: Pass full path
	if openingErr != nil {
		pkg.Alert(pkg.ErrorL, openingErr.Error())
	}

	pkg.Alert(pkg.SuccessL, "Note updated successfully: "+selected)
}
