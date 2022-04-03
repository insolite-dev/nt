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
	Aliases: []string{"rn", "mv"},
	Short:   "Change/Update note's name",
	Run:     runRenameCommand,
}

// initRenameCommand adds renameCommand to main application command.
func initRenameCommand() {
	appCommand.AddCommand(renameCommand)
}

// runRenameCommand runs appropriate service commands to rename note.
func runRenameCommand(cmd *cobra.Command, args []string) {
	// Use arguments for old and new note names.
	if len(args) == 2 {
		rename(args[0], args[1])
		return
	}

	// Use first argument for old note name.
	if len(args) == 1 {
		askAndRename(args[0])
		return
	}

	// Generate array of all note names, if arguments is empty.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNotePrompt("rename:", pkg.MapNodesList(notes)),
		&selected,
	)

	askAndRename(selected)
}

// askAndRename asks user for new name,
// (for selected note), and changes its name.
func askAndRename(selected string) {
	var newname string
	survey.AskOne(assets.NewNamePrompt(selected), &newname)

	rename(selected, newname)
}

// rename takes selected and newname, then makes changes and alerts it.
func rename(selected string, newname string) {
	// Generate editable node by current note and updated note.
	editNode := models.EditNode{
		Current: models.Node{Title: selected},
		New:     models.Node{Title: newname},
	}

	if err := service.Rename(editNode); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Node renamed successfully: "+newname)
}
