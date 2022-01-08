// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import "github.com/anonistas/notya/lib/models"

// ServiceRepo is a repository template for all services.
//
// So, local service is just a ServiceRepo implementation which is connected to local device storage.
// Or we could have remote service, which would be also a ServiceRepo implementation which that is connected to remote DB.
type ServiceRepo interface {
	Init() error

	Open(note models.Note) error
	Remove(note models.Note) error
	Create(note models.Note) (*models.Note, error)
	View(note models.Note) (*models.Note, error)
	Edit(note models.Note) (*models.Note, error)
	Rename(editnote models.EditNote) (*models.Note, error)
	Copy(note models.Note) (*string, error)

	GetAll() ([]string, error)
}
