// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import "github.com/anonistas/notya/lib/models"

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
	// Path returns the path of current base service.
	// In case of local storage implementation, path would be the folder path of the notes.
	Path() string

	// Init setups all kinda minimal services for application.
	Init() error
	Settings() (*models.Settings, error)
	WriteSettings(settings models.Settings) error

	// General functions that used for both [Note]s and [Folder]s
	Open(node models.Node) error
	Remove(node models.Node) error
	Rename(editNode models.EditNode) error
	GetAll() ([]models.Node, error)

	// Note(file) related functions.
	Create(note models.Note) (*models.Note, error)
	View(note models.Note) (*models.Note, error)
	Edit(note models.Note) (*models.Note, error)

	// Folder(directory) related functions.
	Mkdir(dir models.Folder) (*models.Folder, error)

	// MoveNotes moves all exiting notes from CURRENT directory
	// to new one, appropriate by settings which comes from arguments.
	MoveNotes(settings models.Settings) error
}
