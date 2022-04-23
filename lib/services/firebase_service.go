// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirebaseService is a class implementation of service repo.
// Which's methods are based on Firebase client.
// ...
type FirebaseService struct {
	LS      ServiceRepo // embedded local service.
	Stdargs models.StdArgs
	Config  models.Settings

	// Firebase related.
	Ctx       context.Context
	FireApp   *firebase.App
	FireAuth  *auth.Client
	FireStore *firestore.Client
}

// Mark [FirebaseService] as [ServiceRepo].
var _ ServiceRepo = &FirebaseService{}

// NewFirebaseService creates new firebase service by given arguments.
func NewFirebaseService(stdargs models.StdArgs, ls ServiceRepo) *FirebaseService {
	return &FirebaseService{
		LS:      ls,
		Stdargs: stdargs,
		Ctx:     context.Background(),
	}
}

// Path returns current service'base working directory.
func (s *FirebaseService) Path() string {
	return s.LS.Path()
}

// notyaCollection generates the main firestore collection refrence.
func (s *FirebaseService) NotyaCollection(sub *string) firestore.CollectionRef {
	collection := *s.FireStore.Collection(s.Config.Name)

	if sub != nil && sub != &s.Config.Name {
		collection = *collection.Parent.Collection(*sub)
	}

	return collection
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

// Init creates notya working directory into current machine.
func (s *FirebaseService) Init() error {
	if s.FireApp != nil && s.FireStore != nil && s.FireAuth != nil {
		return nil
	}

	localConfig, err := s.LS.Settings()
	if err != nil {
		return err
	}
	s.Config = *localConfig // should be re-written later.

	if err := s.InitFirebase(); err != nil {
		return err
	}

	config, err := s.Settings()
	if status.Code(err) == codes.NotFound {
		if err := s.WriteSettings(*localConfig); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	s.Config = *config // set remote settigns data instead of local.

	return err
}

// Initializes firebase services as [s.FireApp], [s.FireAuth], and [s.FireStore].
func (s *FirebaseService) InitFirebase() error {
	opts := option.WithCredentialsFile(s.Config.FirebaseAccountKey)
	config := &firebase.Config{ProjectID: s.Config.FirebaseProjectID}

	app, err := firebase.NewApp(s.Ctx, config, opts)
	if err != nil {
		return err
	}
	s.FireApp = app

	authClient, err := s.FireApp.Auth(s.Ctx)
	if err != nil {
		return err
	}
	s.FireAuth = authClient

	firestore, err := s.FireApp.Firestore(s.Ctx)
	if err != nil {
		return err
	}
	s.FireStore = firestore

	return nil
}

// Settings gets and returns current settings state data.
func (s *FirebaseService) Settings() (*models.Settings, error) {
	data, err := s.getFireDoc(s.NotyaCollection(nil), models.SettingsName)
	if err != nil {
		return nil, err
	}

	var settings models.Settings
	mapstructure.Decode(data, &settings)

	return &settings, nil
}

// WriteSettings overwrites settings data by given settings model.
func (s *FirebaseService) WriteSettings(settings models.Settings) error {
	if !settings.IsValid() {
		return assets.InvalidSettingsData
	}

	collection := s.NotyaCollection(nil)
	if _, err := collection.Doc(models.SettingsName).Set(s.Ctx, settings.ToJSON()); err != nil {
		return err
	}

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
