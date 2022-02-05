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
	settings  models.Settings
}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService creates new local service by given arguments.
func NewLocalService(stdargs models.StdArgs) *LocalService {
	return &LocalService{stdargs: stdargs}
}

// Init creates notya working directory into current machine.
func (l *LocalService) Init() error {
	// Generate the notya path.
	notyaPath, err := pkg.NotyaPWD(l.settings)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
	}

	l.notyaPath = *notyaPath + "/"
	settingsPath := l.notyaPath + models.SettingsName

	notyaDirSetted := pkg.FileExists(*notyaPath)
	settingsSetted := pkg.FileExists(settingsPath)

	// If settings exists, set it to state.
	if settingsSetted {
		settings, settingsErr := l.Settings()
		if settingsErr != nil {
			return settingsErr
		}

		l.settings = *settings
	}

	// Check if working directories already exists or not.
	if notyaDirSetted && settingsSetted {
		return nil
	}

	// Create new notya working directory, if it not exists.
	if !notyaDirSetted {
		if creatingErr := pkg.NewFolder(*notyaPath); creatingErr != nil {
			return creatingErr
		}
	}

	// Initialize settings file.
	newSettings := models.InitSettings(l.notyaPath)
	if settingsError := l.WriteSettings(newSettings); err != nil {
		return settingsError
	}

	l.settings = newSettings

	return nil
}

// Settings gets and returns current settings state data.
func (l *LocalService) Settings() (*models.Settings, error) {
	settingsPath := l.notyaPath + models.SettingsName

	// Get settings data.
	data, err := pkg.ReadBody(settingsPath)
	if err != nil {
		return nil, err
	}

	settings := models.FromJSON(*data)

	return &settings, nil
}

// WriteSettings, overwrites settings data by given settings model.
func (l *LocalService) WriteSettings(settings models.Settings) error {
	settingsPath := l.notyaPath + models.SettingsName
	if err := pkg.WriteNote(settingsPath, settings.ToByte()); err != nil {
		return err
	}

	return nil
}

// Open, opens given note by editor.
func (l *LocalService) Open(note models.Note) error {
	notePath := l.settings.LocalPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("Note not exists at: %v", note.Title)
		return errors.New(notExists)
	}

	// Open note-file with vi.
	openingErr := pkg.OpenViaEditor(notePath, l.stdargs, l.settings)
	if openingErr != nil {
		return openingErr
	}

	return nil
}

// Remove, deletes given note file, from [notya/note.title]
func (l *LocalService) Remove(note models.Note) error {
	notePath := l.settings.LocalPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("Note not exists at: %v", note.Title)
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
	notePath := l.settings.LocalPath + note.Title

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		alreadyExists := "A note with the name " + fmt.Sprintf("`%v`", note.Title) + " already exists"
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
	notePath := l.settings.LocalPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("Note not exists at: %v", note.Title)
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
	notePath := l.settings.LocalPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("Note not exists at: %v", note.Title)
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
	editnote.Current.Path = l.settings.LocalPath + editnote.Current.Title
	editnote.New.Path = l.settings.LocalPath + editnote.New.Title

	// Check if requested current file exists or not.
	if !pkg.FileExists(editnote.Current.Path) {
		notExists := fmt.Sprintf("Note not exists at: %v", editnote.Current.Title)
		return nil, errors.New(notExists)
	}

	// Check if it's same titles.
	if editnote.Current.Title == editnote.New.Title {
		return nil, errors.New("Current and new name are same")
	}

	// Check if file exists at new note path.
	if pkg.FileExists(editnote.New.Path) {
		alreadyExists := fmt.Sprintf("A note exists at: %v, please provide a unique name", editnote.New.Title)
		return nil, errors.New(alreadyExists)
	}

	// Rename given note.
	if err := os.Rename(editnote.Current.Path, editnote.New.Path); err != nil {
		return nil, err
	}

	return &editnote.New, nil
}

// GetAll, gets all note [names], and returns it as array list.
func (l *LocalService) GetAll() ([]models.Note, error) {
	// Generate array of all file names on LocalPath.
	files, err := pkg.ListDir(l.settings.LocalPath, models.SettingsName)
	if err != nil {
		return nil, err
	}

	if files == nil || len(files) == 0 {
		return nil, errors.New("Empty Directory: not created any note yet")
	}

	// Fetch notes by files.
	notes := []models.Note{}
	for _, name := range files {
		note, err := l.View(models.Note{Title: name})
		if err != nil {
			continue
		}

		notes = append(notes, *note)
	}

	return notes, nil
}

// MoveNote, moves all notes from "CURRENT" path to new path(given by settings parameter).
func (l *LocalService) MoveNotes(settings models.Settings) error {
	notes, err := l.GetAll()
	if err != nil {
		return err
	}

	// Remove notes at default settings' local path.
	for _, note := range notes {
		err := l.Remove(note)
		if err != nil {
			continue
		}
	}

	// Initialize new settings as default.
	l.settings = settings

	// Insert notes at param-settings' local path.
	for _, note := range notes {
		_, err := l.Create(note)
		if err != nil {
			continue
		}
	}

	return nil
}
