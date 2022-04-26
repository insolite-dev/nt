// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// listCommand is a command that used to list all exiting nodes.
var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all notya nodes(files & folders)",
	Run:     runListCommand,
}

// initListCommand adds listCommand to main application command.
func initListCommand() {
	appCommand.AddCommand(listCommand)
}

// runListCommand runs appropriate service functionalities to log all nodes.
func runListCommand(cmd *cobra.Command, args []string) {
	var additional string
	if len(args) > 0 {
		additional = args[0]
	}

	// Generate a list of nodes.
	nodes, _, err := service.GetAll(additional, models.NotyaIgnoreFiles)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	pkg.PrintNodes(nodes)
}
