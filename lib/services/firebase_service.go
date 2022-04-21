// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/anonistas/notya/lib/models"
)

// FirebaseService is a class implementation of service repo.
// Which's methods are based on Firebase client.
// ...
type FirebaseService struct {
	// Notya related.
	Stdargs models.StdArgs
	Config  models.Settings

	// Firebase related.
	FireApp   *firebase.App
	FireAuth  *auth.Client
	FireStore *firestore.Client
}

// Mark [FirebaseService] as [ServiceRepo].
var _ ServiceRepo = &FirebaseService{}

// NewFirebaseService creates new firebase service by given arguments.
func NewFirebaseService(stdargs models.StdArgs) *FirebaseService {
	return &FirebaseService{Stdargs: stdargs}
}

// TODO: add documentation with feature.
func (l *FirebaseService) Path() string {
	return ""
}

// TODO: add documentation with feature.
func (l *FirebaseService) Init() error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Settings() (*models.Settings, error) {
	return nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) WriteSettings(settings models.Settings) error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Open(node models.Node) error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Remove(node models.Node) error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Rename(node models.EditNode) error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) GetAll(additional string) ([]models.Node, []string, error) {
	return nil, nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Create(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) View(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Edit(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Copy(note models.Note) error {
	return nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) Mkdir(dir models.Folder) (*models.Folder, error) {
	return nil, nil
}

// TODO: add documentation with feature.
func (l *FirebaseService) MoveNotes(settings models.Settings) error {
	return nil
}
