//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/insolite-dev/notya/lib/commands"
)

// RunApp executes appCommand.
// It'd be happen only once, on starting program at [main.go].
func RunApp() {
	commands.ExecuteApp()
}
