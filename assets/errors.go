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

	EmptyWorkingDirectory       = errors.New(`Empty working directory, couldn't found any file`)
	InvalidSettingsData         = errors.New(`Invalid settings data, cannot complete operation`)
	InvalidFirebaseProjectID    = errors.New(`Providen firebase-project-id is invalid(or empty)`)
	FirebaseServiceKeyNotExists = errors.New(`Firebase service key file doesn't exists at given path`)
	InvalidFirebaseCollection   = errors.New(`Provided firebase-collection-id is invalid`)
)

// NotExists returns a formatted error message as data-not-exists error.
func NotExists(path, node string) error {
	var msg string
	if len(path) > 1 {
		msg = fmt.Sprintf("%v does not exists at: %v", node, path)
	} else {
		msg = fmt.Sprintf("%v does not exists", node)
	}

	return errors.New(msg)
}

// AlreadyExists returns a formatted error message as data-already-exists error.
func AlreadyExists(path, node string) error {
	msg := fmt.Sprintf("A %v already exists at: %v, please provide a unique title", node, path)
	return errors.New(msg)
}
