// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"os"
	"sort"
	"strings"

	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/atotto/clipboard"
)

// LocalService is a class implementation of service repo.
// Which is connected to local storage of users machine.
// Uses ~notya/ as main root folder for notes and configuration files.
type LocalService struct {
	Stdargs   models.StdArgs
	NotyaPath string
	Config    models.Settings
}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService creates new local service by given arguments.
func NewLocalService(stdargs models.StdArgs) *LocalService {
	return &LocalService{Stdargs: stdargs}
}

// GeneratePath returns non-zero-valuable string path from given additional sub-path(title of node).
func (l *LocalService) GeneratePath(title string) string {
	local := l.Config.LocalPath

	if string(local[len(local)-1]) != "/" {
		local += "/"
	}

	return local + title
}

// Type returns type of LocalService - LOCAL
func (l *LocalService) Type() string {
	return LOCAL.ToStr()
}

// Path returns current service's base working directory.
func (l *LocalService) Path() string {
	return l.NotyaPath
}

// StateConfig returns current configuration of state i.e [l.Config].
func (l *LocalService) StateConfig() models.Settings {
	return l.Config
}

// Init creates notya working directory into current machine.
func (l *LocalService) Init() error {
	// Generate the notya path.
	notyaPath, err := pkg.NotyaPWD(l.Config)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return err
	}

	l.NotyaPath = *notyaPath + "/"
	settingsPath := l.NotyaPath + models.SettingsName

	notyaDirSetted := pkg.FileExists(*notyaPath)
	settingsSetted := pkg.FileExists(settingsPath)

	// If settings exists, set it to state.
	if settingsSetted {
		settings, settingsErr := l.Settings(nil)
		if settingsErr != nil {
			return settingsErr
		}

		l.Config = *settings
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
	newSettings := models.InitSettings(l.NotyaPath)
	if settingsError := l.WriteSettings(newSettings); err != nil {
		return settingsError
	}

	l.Config = newSettings

	return nil
}

// Settings gets and returns current settings state data.
func (l *LocalService) Settings(p *string) (*models.Settings, error) {
	var settingsPath string
	if p != nil && len(*p) != 0 {
		settingsPath = l.GeneratePath(*p)
	} else {
		settingsPath = l.NotyaPath + models.SettingsName
	}

	data, err := pkg.ReadBody(settingsPath)
	if err != nil {
		return nil, err
	}

	settings := models.DecodeSettings(*data)

	return &settings, nil
}

// WriteSettings overwrites settings data by given settings model.
func (l *LocalService) WriteSettings(settings models.Settings) error {
	settingsPath := l.NotyaPath + models.SettingsName

	if !settings.IsValid() {
		return assets.InvalidSettingsData
	}

	if writeErr := pkg.WriteNote(settingsPath, settings.ToByte()); writeErr != nil {
		return writeErr
	}

	return nil
}

// OpenSettings opens given settings via editor.
func (l *LocalService) OpenSettings(settings models.Settings) error {
	path := models.SettingsName
	if len(settings.ID) > 0 {
		path = settings.ID
	}

	// We could use open func of node, in local service.
	return l.Open(models.Node{Title: path})
}

// Open opens given node(file or folder) via editor.
func (l *LocalService) Open(node models.Node) error {
	nodePath := l.GeneratePath(node.Title)

	if len(strings.Trim(node.Title, " ")) < 1 || !pkg.FileExists(nodePath) {
		return assets.NotExists(node.Title, "File or Directory")
	}

	openingErr := pkg.OpenViaEditor(nodePath, l.Stdargs, l.Config)
	if openingErr != nil {
		return openingErr
	}

	return nil
}

// Remove deletes given node.
func (l *LocalService) Remove(node models.Node) error {
	nodePath := l.GeneratePath(node.Title)

	if len(strings.Trim(node.Title, " ")) < 1 || !pkg.FileExists(nodePath) {
		return assets.NotExists(node.Title, "File or Directory")
	}

	// Check for directory, to remove sub nodes of it.
	if pkg.IsDir(nodePath) {
		subNodes, _, err := l.GetAll(node.StructAsFolder().Title, []string{})
		if err != nil && err != assets.EmptyWorkingDirectory {
			return err
		}

		// Sort subNodes via decreasing order.
		sort.Slice(
			subNodes,
			func(i, j int) bool { return len(subNodes[i].Title) > len(subNodes[j].Title) },
		)

		// Remove all sub nodes of directory that're based at [nodePath].
		for _, subNode := range subNodes {
			title := node.StructAsFolder().Title + subNode.StructAsNote().Title
			if err := l.Remove(models.Node{Title: title}); err != nil {
				return err
			}
		}
	}

	if err := pkg.Delete(nodePath); err != nil {
		return err
	}

	return nil
}

// Rename changes given note's name.
func (l *LocalService) Rename(editNode models.EditNode) error {
	editNode.Current.Path = l.Config.LocalPath + editNode.Current.Title
	editNode.New.Path = l.Config.LocalPath + editNode.New.Title

	if len(strings.Trim(editNode.Current.Title, " ")) < 1 || !pkg.FileExists(editNode.Current.Path) {
		return assets.NotExists(editNode.Current.Title, "File or Directory")
	}

	if editNode.Current.Title == editNode.New.Title {
		return assets.SameTitles
	}

	if pkg.FileExists(editNode.New.Path) {
		return assets.AlreadyExists(editNode.New.Title, "File or Directory")
	}

	if err := os.Rename(editNode.Current.Path, editNode.New.Path); err != nil {
		return err
	}

	return nil
}

// Create creates new note file.
// and fills it's data by given note model.
func (l *LocalService) Create(note models.Note) (*models.Note, error) {
	notePath := l.GeneratePath(note.Title)

	if pkg.FileExists(notePath) {
		return nil, assets.AlreadyExists(note.Title, "file")
	}

	if creatingErr := pkg.WriteNote(notePath, []byte(note.Body)); creatingErr != nil {
		return nil, creatingErr
	}

	return &models.Note{Title: note.Title, Path: notePath}, nil
}

// View opens note-file from given [note.Name], then takes it body,
// and returns new fully-filled note.
func (l *LocalService) View(note models.Note) (*models.Note, error) {
	notePath := l.GeneratePath(note.Title)

	if len(strings.Trim(note.Title, " ")) < 1 || !pkg.FileExists(notePath) {
		return nil, assets.NotExists(note.Title, "File")
	}

	res, err := pkg.ReadBody(notePath)
	if err != nil {
		return nil, err
	}

	// Re-generate note with full body.
	modifiedNote := models.Note{Title: note.Title, Path: notePath, Body: *res}

	return &modifiedNote, nil
}

// Edit overwrites exiting file's content-body.
func (l *LocalService) Edit(note models.Note) (*models.Note, error) {
	notePath := l.GeneratePath(note.Title)

	if len(strings.Trim(note.Title, " ")) < 1 || !pkg.FileExists(notePath) {
		return nil, assets.NotExists(note.Title, "File")
	}

	if writingErr := pkg.WriteNote(notePath, []byte(note.Body)); writingErr != nil {
		return nil, writingErr
	}

	return &models.Note{Title: note.Title, Path: notePath, Body: note.Body}, nil
}

// Copy writes given notes' body, to machines main clipboard.
func (l *LocalService) Copy(note models.Note) error {
	notePath := l.GeneratePath(note.Title)
	if len(strings.Trim(note.Title, " ")) < 1 || !pkg.FileExists(notePath) {
		return assets.NotExists(note.Title, "File")
	}

	data, err := l.View(note)
	if err != nil {
		return err
	}

	return clipboard.WriteAll(data.Body)
}

// Mkdir creates a new working directory.
func (l *LocalService) Mkdir(dir models.Folder) (*models.Folder, error) {
	title := dir.Title
	folderPath := l.GeneratePath(dir.Title)

	if string(folderPath[len(folderPath)-1]) != "/" {
		folderPath += "/"
	}

	if string(title[len(title)-1]) != "/" {
		title += "/"
	}

	if pkg.FileExists(folderPath) {
		return nil, assets.AlreadyExists(folderPath, "directory")
	}

	if mkdirErr := pkg.NewFolder(folderPath); mkdirErr != nil {
		return nil, mkdirErr
	}

	return &models.Folder{Title: title, Path: folderPath}, nil
}

// GetAll gets all node [names], and returns it as array list.
func (l *LocalService) GetAll(additional string, ignore []string) ([]models.Node, []string, error) {
	path := l.GeneratePath(additional)

	// Generate array of all file names that are located in [path].
	files, pretty, err := pkg.ListDir(path, "", "", ignore, true)
	if err != nil {
		return nil, nil, err
	}

	if files == nil || len(files) == 0 {
		return nil, nil, assets.EmptyWorkingDirectory
	}

	// Generate node list via [files] array.
	nodes := []models.Node{}
	for i, title := range files {
		path := l.GeneratePath(title)
		nodes = append(nodes, models.Node{Title: title, Path: path, Pretty: pretty[i]})
	}

	return nodes, files, nil
}

// MoveNote moves all notes from "CURRENT" path to new path(given by settings parameter).
func (l *LocalService) MoveNotes(settings models.Settings) error {
	nodes, _, err := l.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		// Remove note appropriate by default settings.
		if err := l.Remove(node); err != nil {
			continue
		}

		// Create note appropriate by updated settings.
		node.Path = settings.LocalPath + node.Title
		if _, err := l.Create(node.ToNote()); err != nil {
			continue
		}
	}

	return nil
}

// TODO: add comment doc & functionality.
func (l *LocalService) Fetch(remote ServiceRepo) ([]models.Node, error) {
	return nil, nil
}

// TODO: add comment doc & functionality.
func (l *LocalService) Push(remote ServiceRepo) ([]models.Node, error) {
	return nil, nil
}
