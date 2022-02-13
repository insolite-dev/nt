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

// renameCommand is a command model which used to change name of notes or files.
var renameCommand = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "r"},
	Short:   "Change/Update note's name",
	Run:     runRenameCommand,
}

// initRenameCommand adds renameCommand to main application command.
func initRenameCommand() {
	appCommand.AddCommand(renameCommand)
}

// runRenameCommand runs appropriate service commands to rename note.
func runRenameCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) == 2 {
		editableNote := models.EditNote{
			Current : models.Note{Title: args[0]},
			New: models.Note{Title: args[1]},
		}

		rename(editableNote.Current.Title, editableNote.New.Title)
		return
	} 

	if len(args) == 1 {
		askAndRename(args[0])
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
		assets.ChooseNotePrompt("Choose a note to rename:", pkg.MapNotesList(notes)),
		&selected,
	)

	askAndRename(selected)
}

// askAndRename asks user for new name,
// (for selected note), and changes its name.
func askAndRename(selected string) {
	var newname string
	survey.AskOne(assets.NewNamePrompt(selected), &newname)

	// Generate editable note by current note and updated note.
	editableNote := models.EditNote{
		Current: models.Note{Title: selected},
		New:     models.Note{Title: newname},
	}

	if _, err := service.Rename(editableNote); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Note renamed successfully: "+newname)
}

func rename(selected string, newname string) {
	// Generate editable note by current note and updated note.
	editableNote := models.EditNote{
		Current: models.Note{Title: selected},
		New:     models.Note{Title: newname},
	}

	if _, err := service.Rename(editableNote); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Note renamed successfully: "+newname)
}