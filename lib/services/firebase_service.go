// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
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
	return s.Config.FirebaseCollection
}

// StateConfig returns current configuration of state i.e [s.Config].
func (s *FirebaseService) StateConfig() models.Settings {
	return s.Config
}

// notyaCollection generates the main firestore collection refrence.
func (s *FirebaseService) NotyaCollection() firestore.CollectionRef {
	return *s.FireStore.Collection(s.Config.FirebaseCollection)
}

// IsDocumentExists checks if element at given title exists or not.
func (s *FirebaseService) IsNotDocumentExists(title string) bool {
	if len(strings.Trim(title, " ")) < 1 {
		return true
	}

	collection := s.NotyaCollection()
	_, err := collection.Doc(title).Get(s.Ctx)

	return status.Code(err) == codes.NotFound
}

// GetFireDoc gets concrete collection's concrete data (as map).
func (s *FirebaseService) GetFireDoc(collection firestore.CollectionRef, doc string) (res map[string]interface{}, err error) {
	docSnap, err := collection.Doc(doc).Get(s.Ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, assets.NotExists(fmt.Sprintf("%v collection", s.Config.FirebaseCollection), doc)
		}

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

	if len(s.Config.FirebaseProjectID) == 0 {
		return assets.InvalidFirebaseProjectID
	}

	// Check validness of firebase account key.
	if !pkg.FileExists(s.Config.FirebaseAccountKey) || len(s.Config.FirebaseAccountKey) == 0 {
		return assets.FirebaseServiceKeyNotExists
	}

	if len(s.Config.FirebaseCollection) == 0 {
		s.Config.FirebaseCollection = s.Config.Name
	}

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
	data, err := s.GetFireDoc(s.NotyaCollection(), models.SettingsName)
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

	collection := s.NotyaCollection()
	if _, err := collection.Doc(models.SettingsName).Set(s.Ctx, settings.ToJSON()); err != nil {
		return err
	}

	return nil
}

// TODO: add documentation & feature.
// TODO: impl after [s.View].
func (s *FirebaseService) Open(node models.Node) error {
	return s.LS.Open(node)
}

// Remove deletes given node.
func (s *FirebaseService) Remove(node models.Node) error {
	collection := s.NotyaCollection()

	if s.IsNotDocumentExists(node.Title) {
		return assets.NotExists("", node.Title)
	}

	_, err := collection.Doc(node.Title).Delete(s.Ctx)
	return err
}

// TODO: add documentation & feature.
func (s *FirebaseService) Rename(node models.EditNode) error {
	return nil
}

// GetAll returns all elements from notya collection.
func (s *FirebaseService) GetAll(additional string) ([]models.Node, []string, error) {
	var nodes []models.Node
	var titles []string

	collection := s.NotyaCollection()
	iter := collection.Documents(s.Ctx)
	defer iter.Stop()

	for {
		ignoreCurrent := false

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nodes, titles, err
		}

		for _, ignore := range models.NotyaIgnoreFiles {
			if doc.Ref.ID == ignore {
				ignoreCurrent = true // mark current loop as ignorable.
			}
		}

		if ignoreCurrent {
			ignoreCurrent = false // reset ignorable for next item.
			continue
		}

		// Decode data to node
		var node models.Node
		var _ = mapstructure.Decode(doc.Data(), &node)

		nodes = append(nodes, node)
		titles = append(titles, doc.Ref.ID)
	}

	return nodes, titles, nil
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
