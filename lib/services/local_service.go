//local_ser
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package services

import (
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
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
func (l *LocalService) GeneratePath(base string, n models.Node) (string, error) {
	path := n.GetPath(l.Type())

	if strings.Trim(path, " ") != "" {
		return path, nil
	}

	if string(base[len(base)-1]) != "/" {
		base += "/"
	}

	// If the title doesn't exists, we have to break up adding up empty string
	// to home path and returning it.
	if len(strings.Trim(n.Title, " ")) == 0 {
		return base + n.Title, errors.New("returned home path")
	}

	return base + n.Title, nil
}

// Type returns type of LocalService - LOCAL
func (l *LocalService) Type() string {
	return LOCAL.ToStr()
}

// Path returns current service's base working directory.
func (l *LocalService) Path() (string, string) {
	return l.NotyaPath, l.Config.NotesPath
}

// StateConfig returns current configuration of state i.e [l.Config].
func (l *LocalService) StateConfig() models.Settings {
	return l.Config
}

// Init creates notya working directory into current machine.
func (l *LocalService) Init(settings *models.Settings) error {
	notyaPath, err := pkg.NotyaPWD(l.Config)
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return err
	}

	l.NotyaPath = *notyaPath + "/"
	settingsPath := l.NotyaPath + models.SettingsName

	notyaDirSetted := pkg.FileExists(l.NotyaPath)
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
		settingsPath, _ = l.GeneratePath(l.NotyaPath, models.Node{Title: *p})
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

	if writeErr := pkg.WriteNote(settingsPath, settings.ToString()); writeErr != nil {
		return writeErr
	}

	return nil
}

// IsNodeExists checks for a file or folder at [node.Path]
// or at generated path from [node.Title].
// Note: rather than remote services, error checking is not required.
func (l *LocalService) IsNodeExists(node models.Node) (bool, error) {
	path, err := l.GeneratePath(l.Config.NotesPath, node)
	if err != nil {
		return false, err
	}

	exists := pkg.FileExists(path) || len(strings.Trim(node.Title, " ")) < 1
	return exists, nil
}

// OpenSettings opens given settings via editor.
func (l *LocalService) OpenSettings(settings models.Settings) error {
	path := l.NotyaPath + models.SettingsName
	if len(settings.ID) > 0 {
		path = l.NotyaPath + settings.ID
	}

	settingsNode := models.Node{Path: map[string]string{l.Type(): path}}
	if nodeExists, _ := l.IsNodeExists(settingsNode); !nodeExists {
		return assets.NotExists(path, "A configuration file")
	}

	return pkg.OpenViaEditor(path, l.Stdargs, l.Config)
}

// Open opens given node(file or folder) via editor.
func (l *LocalService) Open(node models.Node) error {
	if nodeExists, _ := l.IsNodeExists(node); !nodeExists {
		return assets.NotExists(node.Title, "File")
	}

	path, err := l.GeneratePath(l.Config.NotesPath, node)
	if err != nil {
		return nil
	}

	return pkg.OpenViaEditor(path, l.Stdargs, l.Config)
}

// Remove deletes given node.
func (l *LocalService) Remove(node models.Node) error {
	if nodeExists, _ := l.IsNodeExists(node); !nodeExists {
		return assets.NotExists(node.Title, "File or Directory")
	}

	nodePath, err := l.GeneratePath(l.Config.NotesPath, node)
	if err != nil {
		return nil
	}

	// Check for directory, to remove sub nodes of it.
	if pkg.IsDir(nodePath) {
		subNodes, _, err := l.GetAll(pkg.NormalizePath(node.Title), "", []string{})
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
			title := node.ToFolder().Title + subNode.ToNote().Title
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

// Rename changes given file's or folder's name.
func (l *LocalService) Rename(editNode models.EditNode) error {
	editNode.Current.Path = map[string]string{l.Type(): l.Config.NotesPath + editNode.Current.Title}
	editNode.New.Path = map[string]string{l.Type(): l.Config.NotesPath + editNode.New.Title}

	if currentExists, _ := l.IsNodeExists(editNode.Current); !currentExists {
		return assets.NotExists(editNode.Current.Title, "File or Directory")
	}

	if editNode.Current.Title == editNode.New.Title {
		return assets.SameTitles
	}

	if newExists, _ := l.IsNodeExists(editNode.New); newExists {
		return assets.AlreadyExists(editNode.New.Title, "File or Directory")
	}

	current := editNode.Current.GetPath(l.Type())
	edited := editNode.New.GetPath(l.Type())

	if err := os.Rename(current, edited); err != nil {
		return err
	}

	return nil
}

// ClearNodes removes all nodes from local (including folders).
func (l *LocalService) ClearNodes() ([]models.Node, []error) {
	nodes, _, err := l.GetAll("", "", models.NotyaIgnoreFiles)
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
	notePath, err := l.GeneratePath(l.Config.NotesPath, note.ToNode())
	if err != nil {
		return nil, assets.InvalidPathForAct
	}

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); nodeExists {
		return nil, assets.AlreadyExists(note.Title, "file")
	}

	if creatingErr := pkg.WriteNote(notePath, note.Body); creatingErr != nil {
		return nil, creatingErr
	}

	return &models.Note{Title: note.Title, Path: map[string]string{l.Type(): notePath}}, nil
}

// View opens note-file from given [note.Name], then takes it body,
// and returns new fully-filled note.
func (l *LocalService) View(note models.Note) (*models.Note, error) {
	notePath, err := l.GeneratePath(l.Config.NotesPath, note.ToNode())
	if err != nil {
		return nil, assets.InvalidPathForAct
	}

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); !nodeExists {
		return nil, assets.NotExists(note.Title, "File")
	}

	res, err := pkg.ReadBody(notePath)
	if err != nil {
		return nil, err
	}

	// Re-generate note with full body.
	modifiedNote := models.Note{Title: note.Title, Path: map[string]string{l.Type(): notePath}, Body: *res}

	return &modifiedNote, nil
}

// Edit overwrites exiting file's content-body.
func (l *LocalService) Edit(note models.Note) (*models.Note, error) {
	notePath, err := l.GeneratePath(l.Config.NotesPath, note.ToNode())
	if err != nil {
		return nil, assets.InvalidPathForAct
	}

	if nodeExists, _ := l.IsNodeExists(note.ToNode()); !nodeExists {
		return nil, assets.NotExists(note.Title, "File")
	}

	if writingErr := pkg.WriteNote(notePath, note.Body); writingErr != nil {
		return nil, writingErr
	}

	return &models.Note{Title: note.Title, Path: map[string]string{l.Type(): notePath}, Body: note.Body}, nil
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

// Cut, copies note data to machine's clipboard and removes it instantly.
func (l *LocalService) Cut(note models.Note) (*models.Note, error) {
	if err := l.Copy(note); err != nil {
		return nil, err
	}

	n, err := l.View(note)
	if err != nil {
		return nil, err
	}

	if err := l.Remove(note.ToNode()); err != nil {
		return nil, err
	}

	return n, nil
}

// Mkdir creates a new working directory.
func (l *LocalService) Mkdir(dir models.Folder) (*models.Folder, error) {
	title := dir.Title

	folderPath, err := l.GeneratePath(l.Config.NotesPath, dir.ToNode())
	if err != nil {
		return nil, assets.InvalidPathForAct
	}

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

	return &models.Folder{Title: title, Path: map[string]string{l.Type(): folderPath}}, nil
}

// GetAll fetches all nodes(files and folders) from current active local directory.
func (l *LocalService) GetAll(additional, typ string, ignore []string) ([]models.Node, []string, error) {
	path, _ := l.GeneratePath(l.Config.NotesPath, models.Node{Title: additional})

	// Generate array of all file names that are located in [path].
	files, pretty, err := pkg.ListDir(path, path, typ, ignore, 0)
	if err != nil {
		return nil, nil, err
	}

	if len(files) == 0 {
		return nil, nil, assets.EmptyWorkingDirectory
	}

	// Generate node list via [files] array.
	nodes := []models.Node{}
	for i, title := range files {
		p, err := l.GeneratePath(l.Config.NotesPath, models.Node{Title: title})
		if err != nil {
			continue
		}

		path := map[string]string{l.Type(): p}
		node := models.Node{Type: models.FOLDER, Title: title, Path: path, Pretty: pretty[i]}

		if !pkg.IsDir(p) {
			data, err := l.View(node.ToNote())
			if err == nil {
				node = models.Node{Type: models.FILE, Title: title, Path: path, Body: data.Body, Pretty: pretty[i]}
			}
		}

		nodes = append(nodes, node)
	}

	return nodes, files, nil
}

// MoveNotes moves all notes from "CURRENT" path to new path(given by settings parameter).
func (l *LocalService) MoveNotes(settings models.Settings) error {
	nodes, _, err := l.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil {
		return err
	}

	couldntMoved := []models.Node{}

	// First iteration for cloning notes from current settings to provided [settings].
	for _, node := range nodes {
		p := node.Path // a original path holder for any error case.

		node.UpdatePath(l.Type(), pkg.NormalizePath(settings.NotesPath)+node.Title)

		if node.IsFolder() {
			if _, err := l.Mkdir(node.ToFolder()); err != nil {
				node.Path = p
				couldntMoved = append(couldntMoved, node)
			}
			continue
		}

		if _, err := l.Create(node.ToNote()); err != nil {
			node.Path = p
			couldntMoved = append(couldntMoved, node)
		}
	}

	// Second iteration for cleaning up current settings.
	for _, node := range nodes {
		// If node couldn't moved to new settings appropriate place,
		// we shouldn't remote it from old settings appropriate place.
		cm := func(n models.Node, couldntMoved []models.Node) bool {
			for _, c := range couldntMoved {
				if c.Title == n.Title && c.GetPath(l.Type()) == n.GetPath(l.Type()) {
					return true
				}
			}
			return false
		}(node, couldntMoved)

		if !cm {
			l.Remove(node)
		}
	}

	return nil
}

// Fetch creates a clone of nodes(that doesn't exists on [l](local-service)) from given [remote] service.
func (l *LocalService) Fetch(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := remote.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	// Sort nodes via title-len ascending order.
	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) < len(nodes[j].Title) },
	)

	fetched := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		exists, _ := l.IsNodeExists(node)
		if exists && !node.IsFolder() {
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

		if node.IsFolder() && !exists {
			if _, err := l.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
			} else {
				fetched = append(fetched, node)
			}

			continue
		}

		if exists {
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
	nodes, _, err := l.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) < len(nodes[j].Title) },
	)

	pushed := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		path, err := l.GeneratePath(l.Config.NotesPath, node)
		if err != nil {
			continue
		}

		exists, _ := remote.IsNodeExists(node)

		if pkg.IsDir(path) && !exists {
			if _, err := remote.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}

			continue
		}

		r, _ := remote.View(node.ToNote())
		if !exists {
			if _, err := remote.Create(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}

			continue
		}

		if r.Body != node.Body {
			if _, err := remote.Edit(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}
		}
	}

	return pushed, errors
}

// Migrate overwrites all notes of given [remote] service with [l](current-service).
func (l *LocalService) Migrate(remote ServiceRepo) ([]models.Node, []error) {
	if _, err := remote.ClearNodes(); err != nil {
		return nil, err
	}

	return l.Push(remote)
}
