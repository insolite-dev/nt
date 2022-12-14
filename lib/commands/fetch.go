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
	"github.com/fatih/color"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/services"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

var fetchCommand = &cobra.Command{
	Use:     "fetch",
	Aliases: []string{"pull"},
	Short:   "Fetch creates a clone of each node from [Y] service to [X] service",
	Run:     runFetchCommand,
}

func initFetchCommand() {
	appCommand.AddCommand(fetchCommand)
}

func runFetchCommand(cmd *cobra.Command, args []string) {
	determineService()
	loading.Start()

	availableServices := []string{}
	// Generate a list of availabe services
	// by not including current service.
	for _, s := range services.Services {
		if service.Type() == s {
			continue
		}

		availableServices = append(availableServices, s)
	}

	loading.Stop()

	// Ask for servie selection.
	var selected string
	survey.AskOne(
		assets.ChooseRemotePrompt(availableServices),
		&selected,
	)
	if len(selected) == 0 {
		os.Exit(-1)
		return
	}

	selectedService := serviceFromType(selected, true)

	loading.Start()
	fetchedNodes, errs := service.Fetch(selectedService)
	loading.Stop()

	if len(fetchedNodes) == 0 && len(errs) == 0 {
		pkg.Print("Already up to date", color.FgHiGreen)
		return
	}

	pkg.PrintErrors("fetch", errs)
	pkg.Alert(pkg.SuccessL, fmt.Sprintf("Fetched %v nodes", len(fetchedNodes)))
}
