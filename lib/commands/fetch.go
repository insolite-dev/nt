// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
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
	_, err := service.Fetch(selectedService)
	loading.Stop()

	// TODO: log got errors.
	// TODO: log fetched nodes.

	if err != nil {
		for i, e := range err {
			pkg.Alert(pkg.ErrorL, fmt.Sprintf("%v | %v", i, e.Error()))
		}
	}
}
