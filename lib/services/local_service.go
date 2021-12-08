// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"errors"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

// LocalService is a class implementation of service repo.
type LocalService struct{}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService, creates new local service.
func NewLocalService() *LocalService {
	return &LocalService{}
}

// CreateNote, creates new note at [notya notes path],
// and fills it's data by given note model.
func (l *LocalService) CreateNote(note models.Note) error {
	// Generate notya notes working directory path.
	notesPath, err := pkg.NotyaPWD()
	if err != nil {
		return err
	}

	notePath := *notesPath + note.Title

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		return errors.New("A file with this: " + note.Title + " name, already exists!")
	}

	// Create new file inside notes.
	creatingErr := pkg.NewFile(notePath, []byte(note.Body))
	if creatingErr != nil {
		return creatingErr
	}

	return nil
}
