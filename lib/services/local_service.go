// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"errors"
	"fmt"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

// LocalService is a class implementation of service repo.
type LocalService struct {
	notyaPath string
}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService, creates new local service.
func NewLocalService(notyapath string) *LocalService {
	return &LocalService{notyaPath: notyapath}
}

// Init creates notya working directory into running machine.
func (l *LocalService) Init() error {
	// Check if working directory already exists
	if pkg.FileExists(l.notyaPath) {
		return errors.New("Notya already initialized before")
	}

	// Create new notya working directory.
	creatingErr := pkg.NewFolder(l.notyaPath)
	if creatingErr != nil {
		return creatingErr
	}

	return nil
}

// CreateNote, creates new note at [notya notes path],
// and fills it's data by given note model.
func (l *LocalService) CreateNote(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		alreadyExists := "A file with the name " + fmt.Sprintf("`%v`", note.Title) + " already exists"
		return errors.New(alreadyExists)
	}

	// Create new file inside notes.
	creatingErr := pkg.NewFile(notePath, []byte(note.Body))
	if creatingErr != nil {
		return creatingErr
	}

	return nil
}
