// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

// LocalService is a class implementation of service repo.
type LocalService struct {
	notyaPath string
	stdargs   models.StdArgs
}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService creates new local service by given arguments.
func NewLocalService(stdargs models.StdArgs) *LocalService {
	return &LocalService{stdargs: stdargs}
}

// Init creates notya working directory into current machine.
func (l *LocalService) Init() error {
	// Generate notya path.
	notyaPath, err := pkg.NotyaPWD()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
	}

	l.notyaPath = *notyaPath

	// Check if working directory already exists.
	if pkg.FileExists(*notyaPath) {
		return nil
	}

	// Create new notya working directory.
	if creatingErr := pkg.NewFolder(*notyaPath); creatingErr != nil {
		return creatingErr
	}

	return nil
}

// Open, opens given note with VI.
func (l *LocalService) Open(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
		return errors.New(notExists)
	}

	// Open note-file with vi.
	openingErr := pkg.OpenFileWithVI(notePath, l.stdargs)
	if openingErr != nil {
		return openingErr
	}

	return nil
}

// Remove, deletes given note file, from [notya/note.title]
func (l *LocalService) Remove(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
		return errors.New(notExists)
	}

	// Delete the note from [notePath].
	if err := pkg.Delete(notePath); err != nil {
		return err
	}

	return nil
}

// Create, creates new note file at [notya notes path],
// and fills it's data by given note model.
func (l *LocalService) Create(note models.Note) (*models.Note, error) {
	notePath := l.notyaPath + note.Title

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		alreadyExists := "A file with the name " + fmt.Sprintf("`%v`", note.Title) + " already exists"
		return nil, errors.New(alreadyExists)
	}

	// Create new file.
	if creatingErr := pkg.WriteNote(notePath, []byte(note.Body)); creatingErr != nil {
		return nil, creatingErr
	}

	return &models.Note{Title: note.Title, Path: notePath}, nil
}

// View, opens note-file from given [note.Name], then takes it body,
// and returns new fully-filled note.
func (l *LocalService) View(note models.Note) (*models.Note, error) {
	notePath := l.notyaPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
		return nil, errors.New(notExists)
	}

	// Open and read body of note.
	res, err := pkg.ReadBody(notePath)
	if err != nil {
		return nil, err
	}

	// Re-generate note with full body.
	modifiedNote := models.Note{Title: note.Title, Path: notePath, Body: *res}

	return &modifiedNote, nil
}

// Edit, overwrites exiting file's content-body.
func (l *LocalService) Edit(note models.Note) (*models.Note, error) {
	notePath := l.notyaPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
		return nil, errors.New(notExists)
	}

	// Overwrite note's body.
	if writingErr := pkg.WriteNote(notePath, []byte(note.Body)); writingErr != nil {
		return nil, writingErr
	}

	return &models.Note{Title: note.Title, Path: notePath, Body: note.Body}, nil
}

// Rename, changes given note's name.
func (l *LocalService) Rename(editnote models.EditNote) (*models.Note, error) {
	editnote.Current.Path = l.notyaPath + editnote.Current.Title
	editnote.New.Path = l.notyaPath + editnote.New.Title

	// Check if requested current file exists or not.
	if !pkg.FileExists(editnote.Current.Path) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", editnote.Current.Title)
		return nil, errors.New(notExists)
	}

	// Check if it's same titles.
	if editnote.Current.Title == editnote.New.Title {
		return nil, errors.New("Current and new name are same")
	}

	// Check if file exists at new note path.
	if pkg.FileExists(editnote.New.Path) {
		alreadyExists := fmt.Sprintf("A file exists at: notya/%v, please provide a unique name", editnote.New.Title)
		return nil, errors.New(alreadyExists)
	}

	// Rename given note.
	if err := os.Rename(editnote.Current.Path, editnote.New.Path); err != nil {
		return nil, err
	}

	return &editnote.New, nil
}

// GetAll, gets all note [names], and returns it as array list.
func (l *LocalService) GetAll() ([]string, error) {
	// Generate array of all notes' names.
	notes, err := pkg.ListDir(l.notyaPath)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
