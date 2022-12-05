//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
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

	// TODO: list all active remote services.
}

// runRemoteConnectCommand connects to a new remote service connection.
func runRemoteConnectCommand(cmd *cobra.Command, args []string) {
	determineService()

	// TODO: add survey to add new remote service.
	// details: https://github.com/insolite-dev/notya/issues/83
}

// runRemoteDisconnectCommand removes connection from concrete remove service
func runRemoteDisconnectCommand(cmd *cobra.Command, args []string) {
	determineService()

	// TODO: add functionality to disconnect from remote service.
}
