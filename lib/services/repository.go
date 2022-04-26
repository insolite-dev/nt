// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import "github.com/anonistas/notya/lib/models"

var (
	LOCAL ServiceType = "LOCAL"
	FIRE  ServiceType = "FIREBASE"
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
	// Type returns its type.
	// - LOCAL, if it's local service implementation.
	// - FIRE, if it's firebase service implementation.
	// and etc ...
	Type() string

	// Path returns the path of current base service.
	// In case of local storage implementation, path would be the folder path of the notes.
	Path() string

	// Current config data of service implementation.
	StateConfig() models.Settings

	// Init setups all kinda minimal services for application.
	Init() error
	Settings() (*models.Settings, error)
	WriteSettings(settings models.Settings) error

	// General functions that used for both [Note]s and [Folder]s
	Open(node models.Node) error
	Remove(node models.Node) error
	Rename(editNode models.EditNode) error

	//
	// TODO: Add functionality to provide ignorable files
	// Like: ignore folders, files etc.
	//
	GetAll(additional string) ([]models.Node, []string, error)

	// Note(file) related functions.
	Create(note models.Note) (*models.Note, error)
	View(note models.Note) (*models.Note, error)
	Edit(note models.Note) (*models.Note, error)
	Copy(note models.Note) error

	// Folder(directory) related functions.
	Mkdir(dir models.Folder) (*models.Folder, error)

	// MoveNotes moves all exiting notes from CURRENT directory
	// to new one, appropriate by settings which comes from arguments.
	MoveNotes(settings models.Settings) error
}
