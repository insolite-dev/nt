//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

// removeCommand is a command model that used to remove a file or folder.
var removeCommand = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove/Delete a notya element",
	Run:     runRemoveCommand,
}

var removeAll bool

// initRemoveCommand adds removeCommand to main application command.
func initRemoveCommand() {
	removeCommand.Flags().BoolVarP(
		&removeAll, "all", "a", false,
		"Remove all nodes (including nodes under the directories)",
	)

	appCommand.AddCommand(removeCommand)
}

// runRemoveCommand runs appropriate service commands to remove a file or folder.
func runRemoveCommand(cmd *cobra.Command, args []string) {
	determineService()

	if removeAll {
		loading.Start()
		clearedNodes, errs := service.ClearNodes()
		loading.Stop()

		pkg.PrintErrors("remove", errs)
		pkg.Alert(pkg.SuccessL, fmt.Sprintf("Removed %v nodes", len(clearedNodes)))
		return
	}

	// Take node title from arguments. If it's provided.
	if len(args) > 0 && args[0] != "." {
		removeAndFinish(models.Node{Title: args[0]})
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
		assets.ChooseNodePrompt("node", "remove", nodeNames),
		&selected,
	)

	removeAndFinish(models.Node{Title: selected})
}

// removeAndFinish removes given node and alerts success message if everything is OK.
func removeAndFinish(node models.Node) {
	if len(node.Title) == 0 {
		os.Exit(-1)
		return
	}

	loading.Start()

	err := service.Remove(node)

	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
	}
}
