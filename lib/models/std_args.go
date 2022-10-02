//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

import "io"

// StdArgs is a global std state model for application.
// Makes is easy to test functionalities by specifying std state.
type StdArgs struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
