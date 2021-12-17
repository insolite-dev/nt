// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import "github.com/anonistas/notya/lib/models"

// ServiceRepo is a repository template for all services.
//
// So, local service is just a ServiceRepo implementation which is connected to local device storage.
type ServiceRepo interface {
	Init() error
	CreateNote(note models.Note) error
	ViewNote(note models.Note) (*models.Note, error)
	EditNote(note models.Note) error
	Rename(editnote models.EditNote) error
	Remove(note models.Note) error
}
