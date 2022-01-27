// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// createCommand is a command model that used to create new notes or files.
var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "make"},
	Short:   "Create new note",
	Run:     runCreateCommand,
}

// initCreateCommand adds it to the main application command.
func initCreateCommand() {
	appCommand.AddCommand(createCommand)
}

// runCreateCommand runs appropriate service commands to create new note.
func runCreateCommand(cmd *cobra.Command, args []string) {
	// Take new note's title from arguments, if it's provided.
	if len(args) > 0 {
		title := args[0]
		createAndFinish(title)
		return
	}

	// Ask for title of new note.
	var title string
	survey.Ask(assets.CreatePromptQuestion, &title)

	createAndFinish(title)
}

// createAndFinish asks to edit note and finishes creating loop.
func createAndFinish(title string) {
	// Create new note-file by given title.
	note, err := service.Create(models.Note{Title: title})
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for, open or not created note with editor.
	var openNote bool
	survey.AskOne(assets.OpenViaEditorPromt, &openNote)

	if openNote {
		// Open created note-file to edit it.
		if err := service.Open(*note); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}
	}

	pkg.Alert(pkg.SuccessL, "Created new note: "+note.Title)
}
