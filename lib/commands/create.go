// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// note is early created and more-than-one-time usable empty note variable.
var note models.Note

// createCommand, is a command model which used to create new notes or files.
var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "make"},
	Short:   "Create new note",
	Run:     runCreateCommand,
}

// initCreateCommand sets flags of command, and adds it to main application command.
func initCreateCommand() {
	createCommand.Flags().StringVarP(
		&note.Title,
		"title", "t",
		"new_note.md",
		"name of new note/file",
	)

	appCommand.AddCommand(createCommand)
}

// runCreateCommand runs appropriate service commands to create new note.
func runCreateCommand(cmd *cobra.Command, args []string) {
	// Generate note model from arguments.
	note := models.Note{
		Title: note.Title,
		Path:  NotyaPath + note.Title,
	}

	// Create new note-file by [note].
	if err := service.CreateNote(note); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Alert success message.
	pkg.Alert(pkg.SuccessL, "Successfully created new note: "+note.Title)

	// Reset current note.
	note = models.Note{}
}
