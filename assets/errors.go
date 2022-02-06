// Copyright 2021-2022 present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package assets

import (
	"errors"
	"fmt"
)

// Constant and non modifiable errors.
var (
	SameTitles = errors.New(
		`Provided "current" and "new" title are the same, please provide a different title`,
	)

	EmptyWorkingDirectory = errors.New(`Empty working directory, couldn't found any file`)
)

// NotExists returns a formatted error message as data-not-exists error.
func NotExists(path string) error {
	msg := fmt.Sprintf("File does not exists at: %v", path)
	return errors.New(msg)
}

// AlreadyExists returns a formatted error message as data-already-exists error.
func AlreadyExists(path string) error {
	msg := fmt.Sprintf("A file already exists at: %v, please provide a unique title", path)
	return errors.New(msg)
}
