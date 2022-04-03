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

// removeCommand is a command model that used to remove a file or folder.
var removeCommand = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove/Delete a notya element",
	Run:     runRemoveCommand,
}

// initRemoveCommand adds removeCommand to main application command.
func initRemoveCommand() {
	appCommand.AddCommand(removeCommand)
}

// runRemoveCommand runs appropriate service commands to remove a file or folder.
func runRemoveCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		removeAndFinish(models.Node{Title: args[0]})
		return
	}

	// Generate array of all node names.
	nodes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNotePrompt("remove", pkg.MapNodesList(nodes)),
		&selected,
	)

	removeAndFinish(models.Node{Title: selected})
}

// removeAndFinish removes given node and alerts success message if everything is OK.
func removeAndFinish(node models.Node) {
	if err := service.Remove(node); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.Alert(pkg.SuccessL, "Node removed: "+node.Title)
}
