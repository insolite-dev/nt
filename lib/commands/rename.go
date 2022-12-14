//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

// renameCommand is a command model which used to change name of nodes.
var renameCommand = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "mv"},
	Short:   "Change/Update node's name",
	Run:     runRenameCommand,
}

// initRenameCommand adds renameCommand to main application command.
func initRenameCommand() {
	appCommand.AddCommand(renameCommand)
}

// runRenameCommand runs appropriate service commands to rename a node.
func runRenameCommand(cmd *cobra.Command, args []string) {
	determineService()

	// Use arguments for old and new node names.
	if len(args) == 2 {
		rename(args[0], args[1])
		return
	}

	// Use first argument for old node name.
	if len(args) == 1 {
		askAndRename(args[0])
		return
	}

	loading.Start()

	// Generate array of all node names.
	_, nodeNames, err := service.GetAll("", "", models.NotyaIgnoreFiles)

	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for node selection.
	var selected string
	survey.AskOne(
		assets.ChooseNodePrompt("node", "rename", nodeNames),
		&selected,
	)

	askAndRename(selected)
}

// askAndRename asks user for new name,
// (for selected node), and changes its name.
func askAndRename(selected string) {
	var newname string
	survey.AskOne(assets.NewNamePrompt(selected), &newname)

	if len(newname) == 0 {
		os.Exit(-1)
		return
	}

	rename(selected, newname)
}

// rename takes selected and newname, then makes changes and alerts it.
func rename(selected string, newname string) {
	if len(selected) == 0 || len(newname) == 0 {
		os.Exit(-1)
		return
	}

	// Generate editable node by current node and updated node.
	editNode := models.EditNode{
		Current: models.Node{Title: selected},
		New:     models.Node{Title: newname},
	}

	loading.Start()
	err := service.Rename(editNode)
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
}
