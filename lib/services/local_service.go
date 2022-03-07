// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"os"

	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

// LocalService is a class implementation of service repo.
// Which is connected to local storage of users machine.
// Uses ~notya/ as main root folder for notes and configuration files.
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

// generatePath returns non-zero-valuable string path from given note.
func (l *LocalService) generatePath(note models.Note) string {
	if note.Path != "" {
		return note.Path
	}

	return l.settings.LocalPath + note.Title
}

// Path returns current service's base working directory.
func (l *LocalService) Path() string {
	return l.notyaPath
}

// Init creates notya working directory into current machine.
func (l *LocalService) Init() error {
	// Generate the notya path.
	notyaPath, err := pkg.NotyaPWD(l.settings)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return err
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

	// Check settings validness.
	if !settings.IsValid() {
		return assets.InvalidSettingsData
	}

	if writeErr := pkg.WriteNote(settingsPath, settings.ToByte()); writeErr != nil {
		return writeErr
	}

	return nil
}

// Open, opens given note by editor.
func (l *LocalService) Open(note models.Note) error {
	notePath := l.generatePath(note)

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		return assets.NotExists(note.Title)
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
	notePath := l.generatePath(note)

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		return assets.NotExists(note.Title)
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
	notePath := l.generatePath(note)

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		return nil, assets.AlreadyExists(note.Title)
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
	notePath := l.generatePath(note)

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		return nil, assets.NotExists(note.Title)
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
	notePath := l.generatePath(note)

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		return nil, assets.NotExists(note.Title)
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
		return nil, assets.NotExists(editnote.Current.Title)
	}

	// Check if it's same titles.
	if editnote.Current.Title == editnote.New.Title {
		return nil, assets.SameTitles
	}

	// Check if file exists at new note path.
	if pkg.FileExists(editnote.New.Path) {
		return nil, assets.AlreadyExists(editnote.New.Title)
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
		return nil, assets.EmptyWorkingDirectory
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

	for _, note := range notes {
		// Remove note appropriate by default settings.
		if err := l.Remove(note); err != nil {
			continue
		}

		// Create note appropriate by updated settings.
		note.Path = settings.LocalPath + note.Title
		if _, err := l.Create(note); err != nil {
			continue
		}

	}

	return nil
}
