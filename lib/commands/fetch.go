// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
	"github.com/fatih/color"
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
	selectedService := serviceFromType(selected, true)

	loading.Start()
	fetchedNodes, errs := service.Fetch(selectedService)
	loading.Stop()

	if len(fetchedNodes) == 0 && len(errs) == 0 {
		pkg.Print("Already up to date", color.FgHiGreen)
		return
	}

	pkg.PrintFPRes("Fetched", len(fetchedNodes), "nodes \n")
	pkg.PrintErrors("fetch", errs)
}
