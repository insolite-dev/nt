// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/atotto/clipboard"
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

// StateConfig returns current configuration of state i.e [s.Config].
func (s *FirebaseService) StateConfig() models.Settings {
	return s.Config
}

// notyaCollection generates the main firestore collection refrence.
func (s *FirebaseService) NotyaCollection() firestore.CollectionRef {
	return *s.FireStore.Collection(s.Config.FirePath())
}

// GetFireDoc gets concrete collection's concrete data (as map).
func (s *FirebaseService) GetFireDoc(collection firestore.CollectionRef, doc string) (res map[string]interface{}, err error) {
	docSnap, err := collection.Doc(doc).Get(s.Ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, assets.NotExists("", doc)
		}

		return nil, err
	}

	return docSnap.Data(), nil
}

// Type returns type of FirebaseService - FIRE.
func (s *FirebaseService) Type() string {
	return FIRE.ToStr()
}

// Path returns current service'base working directory.
func (s *FirebaseService) Path() string {
	return s.Config.FirePath()
}

// Init creates notya working directory into current machine.
func (s *FirebaseService) Init() error {
	localConfig, err := s.LS.Settings(nil)
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

	config, err := s.Settings(nil)
	if status.Code(err) == codes.NotFound {
		if err := s.WriteSettings(*localConfig); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	s.Config = *config // set remote settigns data instead of local.

	return nil
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
func (s *FirebaseService) Settings(p *string) (*models.Settings, error) {
	sp := models.SettingsName
	if p != nil && len(*p) != 0 {
		sp = *p
	}

	collection := s.FireStore.Collection(s.Config.Name)
	docSnap, err := collection.Doc(sp).Get(s.Ctx)
	if err != nil {
		return nil, err
	}

	var settings models.Settings
	mapstructure.Decode(docSnap.Data(), &settings)

	return &settings, nil
}

// WriteSettings overwrites settings data by given settings model.
func (s *FirebaseService) WriteSettings(settings models.Settings) error {
	if !settings.IsValid() {
		return assets.InvalidSettingsData
	}

	collection := s.FireStore.Collection(s.Config.Name)
	if _, err := collection.Doc(models.SettingsName).Set(s.Ctx, settings.ToJSON()); err != nil {
		return err
	}

	return nil
}

// IsNodeExists checks if an element(given node) exists at notya collection or not.
// Note: rather than local-service error checking is required.
func (s *FirebaseService) IsNodeExists(node models.Node) (bool, error) {
	if len(strings.Trim(node.Title, " ")) < 1 {
		return true, nil
	}

	collection := s.NotyaCollection()
	_, err := collection.Doc(node.Title).Get(s.Ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return false, nil
	}

	return true, err
}

// OpenSettigns, opens note remotly from firebase.
// caches it on local, makes able to modify after modifing overwrites on db.
func (s *FirebaseService) OpenSettings(settings models.Settings) error {
	prevSettings, err := s.Settings(nil)
	if err != nil {
		return err
	}

	title := time.Now().String() + models.SettingsName[1:]

	note := models.Note{
		Title: title,
		Body:  string(prevSettings.ToByte()),
	}
	if _, err := s.LS.Create(note); err != nil {
		return err
	}

	// Open cloned settings data via editor.
	prevSettings.ID = note.Title
	if err := s.LS.OpenSettings(*prevSettings); err != nil {
		return err
	}

	updatedSettings, err := s.LS.Settings(&prevSettings.ID)
	if err != nil {
		return err
	}

	// Clear cache, and skip error.
	_ = s.LS.Remove(note.ToNode())

	if models.IsUpdated(*prevSettings, *updatedSettings) {
		return s.WriteSettings(*updatedSettings)
	}

	return nil
}

// Open, opens note remotly from firebase.
// caches it on local, makes able to modify after modifing overwrites on db.
func (s *FirebaseService) Open(node models.Node) error {
	data, err := s.View(node.ToNote())
	if err != nil {
		return err
	}

	note := models.Note{Title: data.Title + time.Now().String(), Body: data.Body}
	if _, err := s.LS.Create(note); err != nil {
		return err
	}

	// Open via editor to edit.
	openErr := s.LS.Open(note.ToNode())
	if openErr != nil {
		return openErr
	}

	// Get updated note.
	updatedNote, err := s.LS.View(note)
	if err != nil {
		return err
	}

	// Clear cache, and skip error.
	_ = s.LS.Remove(updatedNote.ToNode())

	note = models.Note{Title: data.Title, Path: data.Path, Body: updatedNote.Body}
	if _, err := s.Edit(note); err != nil {
		return err
	}

	return nil
}

// Remove deletes given node.
func (s *FirebaseService) Remove(node models.Node) error {
	collection := s.NotyaCollection()

	if nodeExists, err := s.IsNodeExists(node); err != nil {
		return err
	} else if !nodeExists {
		return assets.NotExists("", node.Title)
	}

	_, err := collection.Doc(node.Title).Delete(s.Ctx)
	return err
}

// Rename changes reference ID of document.
func (s *FirebaseService) Rename(editNode models.EditNode) error {
	data, err := s.View(editNode.Current.ToNote())
	if err != nil {
		return err
	}

	if editNode.Current.Title == editNode.New.Title {
		return assets.SameTitles
	}

	if nodeExists, err := s.IsNodeExists(editNode.New); err != nil {
		return err
	} else if nodeExists {
		return assets.AlreadyExists(editNode.New.Title, "doc")
	}

	if err := s.Remove(editNode.Current); err != nil {
		return err
	}

	_, createErr := s.Create(models.Note{
		Title: editNode.New.Title,
		Body:  data.Body,
	})

	return createErr
}

// ClearNodes removes all nodes from collection.
func (s *FirebaseService) ClearNodes() ([]models.Node, []error) {
	nodes, _, err := s.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	var res []models.Node
	var errs []error

	for _, n := range nodes {
		if err := s.Remove(n); err != nil {
			errs = append(errs, assets.CannotDoSth("remove", n.Title, err))
			continue
		}

		res = append(res, n)
	}

	return res, errs
}

// GetAll returns all elements from notya collection.
func (s *FirebaseService) GetAll(additional string, ignore []string) ([]models.Node, []string, error) {
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

		for _, ig := range ignore {
			if doc.Ref.ID == ig {
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

		// Since each doc is file, we've not to care about folder pretties.
		node.Pretty = []string{models.NotePretty, doc.Ref.ID}

		nodes = append(nodes, node)
		titles = append(titles, doc.Ref.ID)
	}

	return nodes, titles, nil
}

// Create, creates a new note element at [note.Title] and sets element-body as json.
func (s *FirebaseService) Create(note models.Note) (*models.Note, error) {
	collection := s.NotyaCollection()

	if nodeExists, err := s.IsNodeExists(note.ToNode()); err != nil {
		return nil, err
	} else if nodeExists {
		return nil, assets.AlreadyExists(note.Title, "doc")
	}

	if _, err := collection.Doc(note.Title).Set(s.Ctx, note.ToJSON()); err != nil {
		return nil, err
	}

	return &note, nil
}

// View fetches note from [note.Title].
func (s *FirebaseService) View(note models.Note) (*models.Note, error) {
	collection := s.NotyaCollection()

	data, err := s.GetFireDoc(collection, note.Title)
	if err != nil {
		return nil, err
	}

	var model models.Note
	mapstructure.Decode(data, &model)

	return &model, nil
}

// TODO: add documentation & feature.
func (s *FirebaseService) Edit(note models.Note) (*models.Note, error) {
	collection := s.NotyaCollection()

	if nodeExists, err := s.IsNodeExists(note.ToNode()); err != nil {
		return nil, err
	} else if !nodeExists {
		return nil, assets.NotExists("", note.Title)
	}

	if _, err := collection.Doc(note.Title).Set(s.Ctx, note.ToJSON()); err != nil {
		return nil, err
	}

	return &note, nil
}

// Copy fetches note from [note.Title], and copies its body to machine's clipboard.
func (s *FirebaseService) Copy(note models.Note) error {
	data, err := s.View(note)
	if err != nil {
		return err
	}

	return clipboard.WriteAll(data.Body)
}

// Mkdir does nothing 'cause of firebase document structure.
// Have to returns [assets.FolderingInFirebase].
func (s *FirebaseService) Mkdir(dir models.Folder) (*models.Folder, error) {
	return nil, assets.NotAvailableForFirebase
}

// MoveNote moves all notes from "CURRENT" firebase collection
// to new collection(given by settings parameter).
func (s *FirebaseService) MoveNotes(settings models.Settings) error {
	nodes, _, err := s.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return err
	}

	prevSettings := s.Config
	for _, node := range nodes {
		// Remove note appropriate by default settings
		s.Config.FirebaseCollection = prevSettings.FirebaseCollection
		if err := s.Remove(node); err != nil {
			continue
		}

		// Create note appropriate by updated settings
		s.Config.FirebaseCollection = settings.FirebaseCollection
		if _, err := s.Create(node.ToNote()); err != nil {
			continue
		}
	}

	return nil
}

// Fetch creates a clone of nodes(that doesn't exists on
// [s](firebase-service)) from given [remote] service.
func (s *FirebaseService) Fetch(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := remote.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	fetched := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		isDir := (len(node.Pretty) > 0 && node.Pretty[0] == models.FolderPretty) || string(node.Title[len(node.Title)-1]) == "/"
		if isDir {
			errors = append(errors, assets.CannotDoSth("fetch", node.Title, assets.NotAvailableForFirebase))
			continue
		}

		if exists, err := s.IsNodeExists(node); err != nil {
			errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
			continue
		} else if exists {
			local, err := s.View(node.ToNote())
			if err != nil {
				errors = append(errors, err)
				continue
			}

			if local.Body != node.Body {
				local.Body = node.Body
				if _, err := s.Edit(*local); err != nil {
					errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
					continue
				}

				fetched = append(fetched, node)
			}

			continue
		}

		if _, err := s.Create(node.ToNote()); err != nil {
			errors = append(errors, err)
		} else {
			fetched = append(fetched, node)
		}
	}

	return fetched, errors
}

// Push uploads nodes(that doens't exists on given remote) from [s](current) to given [remote].
func (s *FirebaseService) Push(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := s.GetAll("", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	errors := []error{}
	pushed := []models.Node{}

	for _, node := range nodes {
		exists, err := remote.IsNodeExists(node)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if !exists {
			if _, err := remote.Create(node.ToNote()); err != nil {
				errors = append(errors, err)
			} else {
				pushed = append(pushed, node)
			}

			continue
		}

		r, err := remote.View(node.ToNote())
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if r.Body == node.Body {
			continue
		}

		if _, err := remote.Edit(node.ToNote()); err != nil {
			errors = append(errors, err)
		} else {
			pushed = append(pushed, node)
		}

	}

	return pushed, errors
}

// Migrate overwrites all notes of given [remote] service with [s](firebase-service).
func (s *FirebaseService) Migrate(remote ServiceRepo) ([]models.Node, []error) {
	if _, err := remote.ClearNodes(); err != nil {
		return nil, err
	}

	return s.Push(remote)
}
