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
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

// remoteCommand is a command model that used to manage remote connections.
var remoteCommand = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote connections",
	Run:   runRemoteCommand,
}

// connectToRemoteCommand is a command model that used
// to connect to new remote services.
var connectToRemoteCommand = &cobra.Command{
	Use:   "connect",
	Short: "Configure a connection to new remote service",
	Run:   runRemoteConnectCommand,
}

// disconnectFromRemoteCommand is a command model that used
// to disconnect from exiting remote services.
var disconnectFromRemoteCommand = &cobra.Command{
	Use:   "disconnect",
	Short: "Remove connection from concrete remote service",
	Run:   runRemoteDisconnectCommand,
}

// initRemoteCommand adds [remoteCommand] to the [appCommand].
func initRemoteCommand() {
	remoteCommand.AddCommand(connectToRemoteCommand)
	remoteCommand.AddCommand(disconnectFromRemoteCommand)

	appCommand.AddCommand(remoteCommand)
}

// runRemoteCommand lists all active remote connections of current application.
func runRemoteCommand(cmd *cobra.Command, args []string) {
	determineService()

	loading.Start()
	enabled, disabled := listAllRemote()
	loading.Stop()

	if len(enabled) > 0 {
		pkg.Print("\nConnected Services:", color.FgGreen)
		pkg.PrintServices(pkg.NOCOLOR, enabled)
	}

	if len(disabled) > 0 {
		pkg.Print("\nUnreachable Services:", color.FgYellow)
		pkg.PrintServices(pkg.NOCOLOR, disabled)
	}
}

// runRemoteConnectCommand connects to a new remote service connection.
func runRemoteConnectCommand(cmd *cobra.Command, args []string) {
	determineService()

	_, disabled := listAllRemote()
	loading.Stop()

	if len(disabled) == 0 {
		pkg.Alert(pkg.InfoL, "All remote service options are currently connected. You cannot establish additional connections at this time.")
		return
	}

	// Ask for service selection.
	var selected string
	survey.AskOne(
		assets.ChooseRemotePrompt(disabled),
		&selected,
	)
	if len(selected) == 0 {
		os.Exit(-1)
		return
	}

	switch selected {
	case services.FIRE.ToStr():
		promptResult := models.Settings{}

		// Ask for firebase prompt filling.
		survey.Ask(assets.FirebaseRemoteConnectPromptQuestion, &promptResult)

		loading.Start()

		s := service.StateConfig()
		updatedS := s.CopyWith(nil, nil, nil, nil, &promptResult.FirebaseProjectID, &promptResult.FirebaseAccountKey, &promptResult.FirebaseCollection)

		// Validate provided firebase connection:
		isEnabled := services.IsFirebaseEnabled(updatedS, &localService)

		loading.Stop()

		if !isEnabled {
			pkg.Alert(pkg.ErrorL, "Unable to connect to the specified Firebase project using the provided credentials. Please check your login details and try again.")
			return
		}

		loading.Start()
		service.WriteSettings(updatedS)
		loading.Stop()
	}

	pkg.Alert(pkg.SuccessL, fmt.Sprintf("Successfully connected to the specified %s project.", selected))
}

// runRemoteDisconnectCommand removes connection from concrete remove service
func runRemoteDisconnectCommand(cmd *cobra.Command, args []string) {
	determineService()

	loading.Start()
	enabled, _ := listAllRemote()
	loading.Stop()

	if len(enabled) == 0 {
		pkg.Alert(pkg.InfoL, "There are no active remote connections to disconnect from")
		return
	}

	// Ask for service selection.
	var selected string
	survey.AskOne(
		assets.ChooseRemotePrompt(enabled),
		&selected,
	)
	if len(selected) == 0 {
		os.Exit(-1)
		return
	}

	loading.Start()
	switch selected {
	case services.FIRE.ToStr():
		empty := ("")
		s := service.StateConfig()
		service.WriteSettings(s.CopyWith(nil, nil, nil, nil, &empty, &empty, &empty))
	}

	loading.Stop()

	pkg.Alert(pkg.SuccessL, fmt.Sprintf("Successfully disconnected from specified %s service", selected))
}

// Returns a list of all remote services by splitting them by their enabled or disabled level.
// first returned array includes "enabled" remote services, and second returned array includes "disabled" remote services.
func listAllRemote() ([]string, []string) {
	allEnabled, allDisabled := []string{}, []string{}

	for _, s := range services.RemoteServices {
		switch s {
		case services.FIRE.ToStr():
			if services.IsFirebaseEnabled(service.StateConfig(), &localService) {
				allEnabled = append(allEnabled, s)
			} else {
				allDisabled = append(allDisabled, s)
			}
		}
	}

	return allEnabled, allDisabled
}
