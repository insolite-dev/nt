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
func (l *LocalService) GeneratePath(n models.Node) string {
	if strings.Trim(n.Path, " ") != "" {
		return n.Path
	}

	local := l.Config.LocalPath

	if string(local[len(local)-1]) != "/" {
		local += "/"
	}

	return local + n.Title
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
		settingsPath = l.GeneratePath(models.Node{Title: *p})
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

// IsNodeExists checks for a file or folder at [node.Path]
// or at generated path from [node.Title].
// Note: rather than remote services, error checking is not required.
func (l *LocalService) IsNodeExists(node models.Node) (bool, error) {
	exists := pkg.FileExists(l.GeneratePath(node)) || len(strings.Trim(node.Title, " ")) < 1
	return exists, nil
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
	if nodeExists, _ := l.IsNodeExists(node); !nodeExists {
		return assets.NotExists(node.Title, "File or Directory")
	}

	if err := pkg.OpenViaEditor(
		l.GeneratePath(node),
		l.Stdargs, l.Config,
	); err != nil {
		return err
	}

	return nil
}

// Remove deletes given node.
func (l *LocalService) Remove(node models.Node) error {
	if nodeExists, _ := l.IsNodeExists(node); !nodeExists {
		return assets.NotExists(node.Title, "File or Directory")
	}

	nodePath := l.GeneratePath(node)

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

	if currentExists, _ := l.IsNodeExists(editNode.Current); !currentExists {
		return assets.NotExists(editNode.Current.Title, "File or Directory")
	}

	if editNode.Current.Title == editNode.New.Title {
		return assets.SameTitles
	}

	if newExists, _ := l.IsNodeExists(editNode.New); newExists {
		return assets.AlreadyExists(editNode.New.Title, "File or Directory")
	}

	if err := os.Rename(editNode.Current.Path, editNode.New.Path); err != nil {
		return err
	}

	return nil
}

// ClearNodes removes all nodes from local (including folders).
func (l *LocalService) ClearNodes() ([]models.Node, []error) {
	nodes, _, err := l.GetAll("", models.NotyaIgnoreFiles)
	if err != nil && err.Error() != assets.EmptyWorkingDirectory.Error() {
		return nil, []error{err}
	}

	// Sort nodes via title-len decreasing order.
	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) > len(nodes[j].Title) },
	)

	var res []models.Node
	var errs []error

	for _, n := range nodes {
		if err := l.Remove(n); err != nil {
			errs = append(errs, assets.CannotDoSth("remove", n.Title, err))
			continue
		}

		res = append(res, n)
	}

	return res, errs
}

// Create creates new note file.
// and fills it's data by given note model.
func (l *LocalService) Create(note models.Note) (*models.Note, error) {
	notePath := l.GeneratePath(note.ToNode())

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); nodeExists {
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
	notePath := l.GeneratePath(note.ToNode())

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); !nodeExists {
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
	notePath := l.GeneratePath(note.ToNode())

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); !nodeExists {
		return nil, assets.NotExists(note.Title, "File")
	}

	if writingErr := pkg.WriteNote(notePath, []byte(note.Body)); writingErr != nil {
		return nil, writingErr
	}

	return &models.Note{Title: note.Title, Path: notePath, Body: note.Body}, nil
}

// Copy writes given notes' body, to machines main clipboard.
func (l *LocalService) Copy(note models.Note) error {
	if nodeExists, _ := l.IsNodeExists(note.ToNode()); !nodeExists {
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
	folderPath := l.GeneratePath(dir.ToNode())

	if string(folderPath[len(folderPath)-1]) != "/" {
		folderPath += "/"
	}

	if string(title[len(title)-1]) != "/" {
		title += "/"
	}

	if dirExists, _ := l.IsNodeExists(dir.ToNode()); dirExists {
		return nil, assets.AlreadyExists(folderPath, "directory")
	}

	if mkdirErr := pkg.NewFolder(folderPath); mkdirErr != nil {
		return nil, mkdirErr
	}

	return &models.Folder{Title: title, Path: folderPath}, nil
}

// GetAll gets all node [names], and returns it as array list.
func (l *LocalService) GetAll(additional string, ignore []string) ([]models.Node, []string, error) {
	path := l.GeneratePath(models.Node{Title: additional})

	// Generate array of all file names that are located in [path].
	files, pretty, err := pkg.ListDir(path, "", "", ignore, true)
	if err != nil {
		return nil, nil, err
	}

	if len(files) == 0 {
		return nil, nil, assets.EmptyWorkingDirectory
	}

	// Generate node list via [files] array.
	nodes := []models.Node{}
	for i, title := range files {
		path := l.GeneratePath(models.Node{Title: title})
		node := models.Node{Title: title, Path: path, Pretty: pretty[i]}

		if !pkg.IsDir(path) {
			data, err := l.View(node.ToNote())
			if err == nil {
				node = models.Node{Title: title, Path: path, Pretty: pretty[i], Body: data.Body}
			}
		}

		nodes = append(nodes, node)
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

// Fetch creates a clone of nodes(that doesn't exists on [l](local-service)) from given [remote] service.
func (l *LocalService) Fetch(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := remote.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	// Sort nodes via title-len decreasing order.
	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) > len(nodes[j].Title) },
	)

	fetched := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		isDir := (len(node.Pretty) > 0 && node.Pretty[0] == models.FolderPretty) || string(node.Title[len(node.Title)-1]) == "/"

		if exists, _ := l.IsNodeExists(node); exists && !isDir {
			local, err := l.View(node.ToNote())
			if err != nil {
				errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
				continue
			}

			if local.Body != node.Body {
				local.Body = node.Body
				if _, err := l.Edit(*local); err != nil {
					errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
					continue
				}

				fetched = append(fetched, node)
			}

			continue
		}

		if isDir {
			if _, err := l.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
			} else {
				fetched = append(fetched, node)
			}
			continue
		}

		if _, err := l.Create(node.ToNote()); err != nil {
			errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
		} else {
			fetched = append(fetched, node)
		}
	}

	return fetched, errors
}

// Push uploads nodes(that doesn't exists on given remote) from [l](current) to given [remote].
func (l *LocalService) Push(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := l.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	// Sort nodes via title-len decreasing order.
	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) > len(nodes[j].Title) },
	)

	fetched := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		if pkg.IsDir(l.GeneratePath(node)) {
			if _, err := remote.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				fetched = append(fetched, node)
			}

			continue
		}

		r, err := remote.View(node.ToNote())
		if err != nil && err.Error() != assets.NotExists("", node.Title).Error() {
			errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			continue
		} else if err != nil {
			if _, err := remote.Create(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				fetched = append(fetched, node)
			}

			continue
		}

		if r.Body != node.Body {
			if _, err := remote.Edit(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
				continue
			}

			fetched = append(fetched, node)
		}
	}

	return fetched, errors
}

// Migrate overwrites all notes of given [remote] service with [l](current-service).
func (l *LocalService) Migrate(remote ServiceRepo) ([]models.Node, []error) {
	if _, err := remote.ClearNodes(); err != nil {
		return nil, err
	}

	return l.Push(remote)
}
