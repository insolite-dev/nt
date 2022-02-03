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

// copyCommand is a command model which used to copy notes or files.
var copyCommand = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"c"},
	Short:   "Copy the full note",
	Run:     runCopyCommand,
}

// initCopyCommand initializes copyCommand to the main application command.
func initCopyCommand() {
	appCommand.AddCommand(copyCommand)
}

// rrunCopyCommand runs appropriate service commands copy note data.
func runCopyCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note := models.Note{Title: args[0]}

		copyAndFinish(note)
		return
	}

	// Generate array of all note names.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNotePrompt("Choose a note to copy:", pkg.MapNotesList(notes)),
		&selected,
	)

	copyAndFinish(models.Note{Title: selected})
}

func copyAndFinish(note models.Note) {
	if _, err := service.Copy(note); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Note copied to clipboard")
}
