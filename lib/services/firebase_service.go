// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/anonistas/notya/lib/models"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
)

// FirebaseService is a class implementation of service repo.
// Which's methods are based on Firebase client.
// ...
type FirebaseService struct {
	LS      ServiceRepo // embedded local service.
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
func NewFirebaseService(stdargs models.StdArgs, ls ServiceRepo) *FirebaseService {
	return &FirebaseService{Stdargs: stdargs, LS: ls}
}

// Path returns current service'base working directory.
func (s *FirebaseService) Path() string {
	return s.LS.Path()
}

// TODO: add documentation with feature.
func (s *FirebaseService) Init() error {
	if s.FireApp != nil && s.FireStore != nil && s.FireAuth != nil {
		return nil
	}

	config, err := s.LS.Settings()
	if err != nil {
		return err
	}
	s.Config = *config

	return s.InitFirebase()
}

// Initializes firebase services as [s.FireApp], [s.FireAuth], and [s.FireStore].
func (s *FirebaseService) InitFirebase() error {
	ctx := context.Background()

	opts := option.WithCredentialsFile(s.Config.FirebaseAccountKey)
	config := &firebase.Config{ProjectID: s.Config.Name}

	app, err := firebase.NewApp(ctx, config, opts)
	if err != nil {
		return err
	}
	s.FireApp = app

	authClient, err := s.FireApp.Auth(ctx)
	if err != nil {
		return err
	}
	s.FireAuth = authClient

	firestore, err := s.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	s.FireStore = firestore

	return nil
}

// notyaCollection generates the main firestore collection refrence.
func (s *FirebaseService) NotyaCollection() firestore.CollectionRef {
	return *s.FireStore.Collection(s.Config.Name)
}

// getFireDoc gets concrete collection's concrete data (as map).
func (s *FirebaseService) getFireDoc(collection firestore.CollectionRef, doc string) (res map[string]interface{}, err error) {
	ctx := context.Background()
	docSnap, err := collection.Doc(doc).Get(ctx)

	if err != nil {
		return nil, err
	}

	return docSnap.Data(), nil
}

// Settings gets and returns current settings state data.
func (s *FirebaseService) Settings() (*models.Settings, error) {
	data, err := s.getFireDoc(s.NotyaCollection(), models.SettingsName)
	if err != nil {
		return nil, err
	}

	var settings models.Settings
	mapstructure.Decode(data, &settings)

	return &settings, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) WriteSettings(settings models.Settings) error {
	return nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Open(node models.Node) error {
	return s.LS.Open(node)
}

// TODO: add documentation & feature.
func (s *FirebaseService) Remove(node models.Node) error {
	return nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Rename(node models.EditNode) error {
	return nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) GetAll(additional string) ([]models.Node, []string, error) {
	return nil, nil, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Create(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) View(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Edit(note models.Note) (*models.Note, error) {
	return nil, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Copy(note models.Note) error {
	return nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Mkdir(dir models.Folder) (*models.Folder, error) {
	return nil, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) MoveNotes(settings models.Settings) error {
	return nil
}
