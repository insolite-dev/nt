//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package services

import "github.com/insolite-dev/notya/lib/models"

var (
	LOCAL ServiceType = "LOCAL"
	FIRE  ServiceType = "FIREBASE"

	Services []string = []string{
		LOCAL.ToStr(),
		FIRE.ToStr(),
	}
)

// Custom string struct to define type of services
type ServiceType string

// ToStr returns exact key value of ServiceType.
func (s *ServiceType) ToStr() string {
	switch s {
	case &LOCAL:
		return "LOCAL"
	case &FIRE:
		return "FIREBASE"
	}

	return "undefined"
}

// ServiceRepo is a abstract class for all service implementations.
//     ╭──────╮     ╭────────────────────╮
// ... │ User │ ──▶ │ Interface Commands │
//     ╰──────╯     ╰────────────────────╯
//                            │
//                ╭───────────────────────╮
//                ▼                       ▼
//        ╭───────────────╮       ╭────────────────╮
//        │ Local Service │       │ Remote Service │
//        ╰───────────────╯       ╰────────────────╯
//        Connected to local       Connected to user defined
//        storage, and uses        key-store remote database, and uses
//        ~notya/ as main root     notya/ as base root key map.
//        folder for notes.
//
type ServiceRepo interface {
	// Type returns the current implementation's type.
	// - LOCAL, if it's local service implementation.
	// - FIRE, if it's firebase service implementation.
	// and etc ...
	Type() string

	// Path returns the path of current base service. And base service's notes path.
	// ...   main | notes
	Path() (string, string)

	// Current config data of service implementation.
	StateConfig() models.Settings

	// Init setups all kinda minimal services for application.
	Init() error

	// Settings reads and parses current configuration file and returns
	// it as settings model pointer. In case of a error, setting model will be
	// [nil] and [error] will be provided.
	Settings(p *string) (*models.Settings, error)

	// WriteSettings overwrites current configuration data,
	// with provided [settings] model.
	WriteSettings(settings models.Settings) error

	// OpenSettings opens provided settings with [current] editor
	// that we take it from provided settings.
	OpenSettings(settings models.Settings) error

	// General functions that used for both [Note]s and [Folder]s
	IsNodeExists(node models.Node) (bool, error)
	Open(node models.Node) error
	Remove(node models.Node) error
	Rename(editNode models.EditNode) error
	ClearNodes() ([]models.Node, []error)

	// Note(file) related functions.
	GetAll(additional string, ignore []string) ([]models.Node, []string, error)
	Create(note models.Note) (*models.Note, error)
	View(note models.Note) (*models.Note, error)
	Edit(note models.Note) (*models.Note, error)
	Copy(note models.Note) error
	Cut(note models.Note) (*models.Note, error)

	// Folder(directory) related functions.
	Mkdir(dir models.Folder) (*models.Folder, error)

	// MoveNotes moves all exiting notes from CURRENT directory
	// to new one, appropriate by settings which comes from arguments.
	MoveNotes(settings models.Settings) error

	// Fetch fetches nodes(that doesn't exists
	// on current service) from remote service to local service.
	Fetch(remote ServiceRepo) ([]models.Node, []error)

	// Push uploads all notes from local service to provided remote.
	Push(remote ServiceRepo) ([]models.Node, []error)

	// Migrate clones current service data to [remote] service data.
	// [remote] service data would be cleared and replaced with current service data.
	Migrate(remote ServiceRepo) ([]models.Node, []error)
}
