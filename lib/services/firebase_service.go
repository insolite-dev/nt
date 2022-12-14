//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/atotto/clipboard"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
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

// notyaCollection generates the main firestore collection reference.
func (s *FirebaseService) NotyaCollection() firestore.CollectionRef {
	return *s.FireStore.Collection(s.Config.FirePath())
}

// GeneratePath generates a string valid path from provided "custom" base of collection reference
// and [models.Node] model.
// first returned value would be the full-valid path of provided node in collection,
// and second returned value would be the valid "base" collection of that path.
func (s *FirebaseService) GeneratePath(base *firestore.CollectionRef, n models.Node) (string, *firestore.CollectionRef) {
	collection := s.NotyaCollection()
	if base != nil {
		// if the base collection is provided, which is different than actual
		// provided main-base connection of notya, we have to set it to the value.
		collection = *base
	}

	// If the model's path is not empty, it will be the path
	// that function will generate document reference for it.
	// Otherwise, a combination of collection id and note title will be used as path.
	path := n.GetPath(s.Type())
	if len(path) == 0 {
		path = collection.ID + "/" + n.Title
	}

	return path, &collection
}

// GenerateDoc is [firebase.DocumentRef] generator, that used to generate concrete document reference by a string path.
//
// So, the: `<base-collection>/<document>/<sub-collection>/<sub-document>`
// string will be converted to:
// `Collection("<base-collection>").Doc("<document>").Collection("<sub-collection>").Doc("<sub-document>")`
//
// In case of being node [Folder] type. The collection reference return will be the "sub" collection
// of generated document. But for the [File] type, it'd return nil.
func (s *FirebaseService) GenerateDoc(base *firestore.CollectionRef, n models.Node) (*firestore.DocumentRef, *firestore.CollectionRef) {
	path, collection := s.GeneratePath(base, n)

	segments := strings.Split(path, "/")
	if segments[0] == collection.ID {
		segments = segments[1:]
	}

	doc := *collection.Doc(segments[0])
	for i := 1; i < len(segments); i++ {
		if len(segments[i]) != 0 {
			doc = *doc.Collection("sub").Doc(segments[i])
		}
	}

	if n.IsFile() {
		return &doc, nil
	}

	return &doc, doc.Collection("sub")
}

// GetDoc is a function that used to get document reference as [models.Node].
func (s *FirebaseService) GetDoc(n models.Node) (*models.Node, error) {
	path, _ := s.GeneratePath(nil, n)
	n.UpdatePath(s.Type(), path)

	nDoc, _ := s.GenerateDoc(nil, n)
	docSnapshot, err := nDoc.Get(s.Ctx)

	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return nil, assets.NotExists(path, "File")
		}

		return nil, err
	}

	var model models.Node
	model.FromJson(docSnapshot.Data())

	return &model, nil
}

// Type returns type of FirebaseService - FIRE.
func (s *FirebaseService) Type() string {
	return FIRE.ToStr()
}

// Path returns current service base working directory and name of working collection.
func (s *FirebaseService) Path() (string, string) {
	return s.Config.FirePath(), s.Config.FirebaseCollection
}

// Init creates notya working directory into current machine.
func (s *FirebaseService) Init(settings *models.Settings) error {
	if settings != nil {
		s.Config = *settings
	} else {
		localConfig, err := s.LS.Settings(nil)
		if err != nil {
			return err
		}

		s.Config = *localConfig // should be re-written later.
	}

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
		if err := s.WriteSettings(s.Config); err != nil {
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
func (s *FirebaseService) IsNodeExists(node models.Node) (bool, error) {
	doc, _ := s.GenerateDoc(nil, node)
	if _, err := doc.Get(s.Ctx); err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return false, nil
		}

		return true, err
	}

	return true, nil
}

// OpenSettigns, opens note remotely from firebase.
// caches it on local, makes able to modify after modifying overwrites on db.
func (s *FirebaseService) OpenSettings(settings models.Settings) error {
	prevSettings, err := s.Settings(nil)
	if err != nil {
		return err
	}

	title := time.Now().String() + models.SettingsName[1:]

	note := models.Note{
		Title: title,
		Body:  prevSettings.ToString(),
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

	if pkg.IsSettingsUpdated(*prevSettings, *updatedSettings) {
		return s.WriteSettings(*updatedSettings)
	}

	return nil
}

// Open, opens a remote note in local machine.
// clones it on local, makes able to modify, after modifying, overwrites on it db.
func (s *FirebaseService) Open(node models.Node) error {
	data, err := s.View(node.ToNote())
	if err != nil {
		return err
	}

	splitted := strings.Split(data.Title, "/")
	note := models.Note{Title: splitted[len(splitted)-1] + time.Now().String(), Body: data.Body}
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

// Remove deletes given node from [node.Path].
// If a node doesn't exists at provided note's path,
// it will return a already formatted error message.
func (s *FirebaseService) Remove(node models.Node) error {
	n := node

	path, _ := s.GeneratePath(nil, n)
	n.UpdatePath(s.Type(), path)

	if nodeExists, err := s.IsNodeExists(n); err != nil {
		return err
	} else if !nodeExists {
		return assets.NotExists(n.Title, "File or Directory")
	}

	noteDoc, _ := s.GenerateDoc(nil, n)
	if _, err := noteDoc.Delete(s.Ctx); err != nil {
		return err
	}

	return nil
}

// Rename changes reference ID of document.
func (s *FirebaseService) Rename(editNode models.EditNode) error {
	current, err := s.GetDoc(editNode.Current)
	if err != nil {
		return err
	}

	if editNode.Current.Title == editNode.New.Title {
		return assets.SameTitles
	}

	updated := editNode.New
	updated.Type = current.Type
	updated.Body = current.Body

	newPath, _ := s.GeneratePath(nil, updated)
	updated.Path = current.Path
	updated.UpdatePath(s.Type(), newPath)

	if nodeExists, err := s.IsNodeExists(updated); err != nil {
		return err
	} else if nodeExists {
		return assets.AlreadyExists(updated.Title, "file or folder")
	}

	if err := s.mv(models.EditNode{Current: *current, New: updated}); err != nil {
		return err
	}

	// Dive into sub collection of current folder.
	if current.IsFolder() || updated.IsFolder() {
		_, sub := s.GenerateDoc(nil, *current)

		nodes, _, err := s.ListDir(sub, "", []string{}, 0)
		if err != nil {
			// TODO: shouldn't cut the whole action for one error.
			return err
		}

		sort.Slice(
			nodes,
			func(i, j int) bool { return len(nodes[i].Title) > len(nodes[j].Title) },
		)

		for _, n := range nodes {
			newN := n
			newN = *newN.RebuildParent(*current, updated, s.Type(), s.Config)

			if err := s.Rename(models.EditNode{Current: n, New: newN}); err != nil {
				// TODO: shouldn't cut the whole action for one error.
				return err
			}
		}
	}

	return nil
}

// mv is a sub implementation of [Rename].
// Which used to move file or folder(without sub nodes)
// from current path to new path.
func (s *FirebaseService) mv(editNode models.EditNode) error {
	if editNode.Current.IsFolder() || editNode.New.IsFolder() {
		if _, err := s.Mkdir(editNode.New.ToFolder()); err != nil {
			return err
		}
	} else {
		if _, err := s.Create(editNode.New.ToNote()); err != nil {
			return err
		}
	}

	return s.Remove(editNode.Current)
}

// ClearNodes removes all nodes from collection.
// TODO: improve the speed of clearing
func (s *FirebaseService) ClearNodes() ([]models.Node, []error) {
	nodes, _, err := s.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil && err.Error() != assets.EmptyWorkingDirectory.Error() {
		return nil, []error{err}
	}

	// Sort nodes via title-len decreasing order.
	sort.Slice(
		nodes,
		func(i, j int) bool { return len(nodes[i].Title) > len(nodes[j].Title) },
	)

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

// GetAll fetches all documents and their sub documents (if they exist)
// from a database and returns them as an array of [models.Node]
//
// @param additional path, [typ] that is allowed to fetch, and ignore list.
// @returns an array of all nodes, titles of nodes and error if something went wrong.
func (s *FirebaseService) GetAll(additional, typ string, ignore []string) ([]models.Node, []string, error) {
	collection := s.NotyaCollection()
	if len(additional) > 0 {
		_, c := s.GenerateDoc(&collection, models.Node{Title: additional})
		if c != nil {
			collection = *c
		}
	}

	return s.ListDir(&collection, typ, ignore, 0)
}

// ListDir retrieves the documents and sub-collections from a specified Firebase CollectionRef.
//
// @param {firestore.CollectionRef} path - The Firebase CollectionRef to retrieve documents and sub-collections from.
// @param {string} typ - The type of documents to retrieve.
// @param {[]string} ignore - An array of sub-collection names to ignore.
// @param {int} level - The number of levels deep to retrieve sub-collections.
//
// @returns {[]models.Node, []string, error} A tuple containing an array of retrieved documents
// and sub-collections (models.Node), an array of ignored sub-collection names, and an error if one occurred.
func (s *FirebaseService) ListDir(path *firestore.CollectionRef, typ string, ignore []string, level int) ([]models.Node, []string, error) {
	var res []models.Node
	var titles []string

	iter := path.Documents(s.Ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return res, titles, err
		}

		// Ignore the current document, if it is ignorable.
		if pkg.IsIgnorable(doc.Ref.ID, ignore) {
			continue
		}

		// Decode data to node
		var node models.Node
		node.FromJson(doc.Data())

		if pkg.IsType(typ, node.IsFolder()) {
			node.Pretty = []string{strings.Repeat("  ", level) + node.GenPretty(), doc.Ref.ID}

			res = append(res, node)
			titles = append(titles, node.Title)
		}

		if node.IsFolder() {
			subPath := path.Doc(doc.Ref.ID).Collection("sub")
			sub, subTitles, err := s.ListDir(subPath, typ, ignore, level+1)
			if err != nil {
				// TODO: find a way of effective way of handling error
				continue
			}

			if len(sub) > 0 {
				res = append(res, sub...)
				titles = append(titles, subTitles...)
			}
		}
	}

	return res, titles, nil
}

// Create, creates a new file document at note's path.
// If a node(file or folder) already exists at provided note's path,
// it will return already formatted error message.
func (s *FirebaseService) Create(note models.Note) (*models.Note, error) {
	noteNode := note.ToNode()

	path, _ := s.GeneratePath(nil, noteNode)
	noteNode.UpdatePath(s.Type(), path)

	noteDoc, _ := s.GenerateDoc(nil, noteNode)
	if _, err := noteDoc.Create(s.Ctx, noteNode.ToJSON()); err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.AlreadyExists {
			return nil, assets.AlreadyExists(noteNode.Title, "file")
		}

		return nil, err
	}

	modifiedNote := noteNode.ToNote()
	return &modifiedNote, nil
}

// View, gets the note document from note's path.
// If a node doesn't exists at provided note's path,
// it will return a already formatted error message.
func (s *FirebaseService) View(note models.Note) (*models.Note, error) {
	noteNode := note.ToNode()

	path, _ := s.GeneratePath(nil, noteNode)
	noteNode.UpdatePath(s.Type(), path)

	noteDoc, _ := s.GenerateDoc(nil, noteNode)
	docSnapshot, err := noteDoc.Get(s.Ctx)

	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return nil, assets.NotExists(path, "File")
		}

		return nil, err
	}

	var model models.Note
	mapstructure.Decode(docSnapshot.Data(), &model)

	return &model, nil
}

// Edit, updates the already created note, with locally updated note data.
// If a node doesn't exists at provided note's path,
// it will return a already formatted error message.
func (s *FirebaseService) Edit(note models.Note) (*models.Note, error) {
	noteNode := note.ToNode()

	path, _ := s.GeneratePath(nil, noteNode)
	noteNode.UpdatePath(s.Type(), path)

	noteDoc, _ := s.GenerateDoc(nil, noteNode)
	if _, err := noteDoc.Set(s.Ctx, noteNode.ToJSON()); err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return nil, assets.NotExists(path, "File")
		}

		return nil, err
	}

	modifiedNote := noteNode.ToNote()
	return &modifiedNote, nil
}

// Copy fetches note from [note.Title], and copies its body to machine's clipboard.
func (s *FirebaseService) Copy(note models.Note) error {
	data, err := s.View(note)
	if err != nil {
		return err
	}

	return clipboard.WriteAll(data.Body)
}

// Cut, copies note data to machine's clipboard and removes it instantly.
func (s *FirebaseService) Cut(note models.Note) (*models.Note, error) {
	n, err := s.View(note)
	if err != nil {
		return nil, err
	}

	if err := clipboard.WriteAll(n.Body); err != nil {
		return nil, err
	}

	doc, _ := s.GenerateDoc(nil, note.ToNode())
	if _, err := doc.Delete(s.Ctx); err != nil {
		return nil, err
	}

	return n, nil
}

// Mkdir creates a document in provided folder path(from dir.Path)
// and plus that, creates a sub collection of current folder document.
// that sub collection gonna represent the files/folders that current
// directory includes.
func (s *FirebaseService) Mkdir(dir models.Folder) (*models.Folder, error) {
	dirNode := dir.ToNode()

	path, _ := s.GeneratePath(nil, dirNode)
	dirNode.UpdatePath(s.Type(), path)

	folderDoc, _ := s.GenerateDoc(nil, dirNode)
	if _, err := folderDoc.Create(s.Ctx, dirNode.ToJSON()); err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.AlreadyExists {
			return nil, assets.AlreadyExists(dirNode.Title, "folder")
		}

		return nil, err
	}

	modifiedDir := dirNode.ToFolder()
	return &modifiedDir, nil
}

// MoveNote moves all notes from "CURRENT" firebase collection
// to new collection(given by settings parameter).
func (s *FirebaseService) MoveNotes(settings models.Settings) error {
	nodes, _, err := s.GetAll("", "", models.NotyaIgnoreFiles)
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
	nodes, _, err := remote.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	fetched := []models.Node{}
	errors := []error{}

	for _, node := range nodes {
		exists, _ := s.IsNodeExists(node)
		if exists && !node.IsFolder() {
			local, err := s.View(node.ToNote())
			if err != nil {
				errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
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

		if node.IsFolder() && !exists {
			if _, err := s.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
			} else {
				fetched = append(fetched, node)
			}

			continue
		}

		if exists {
			continue
		}

		if _, err := s.Create(node.ToNote()); err != nil {
			errors = append(errors, assets.CannotDoSth("fetch", node.Title, err))
		} else {
			fetched = append(fetched, node)
		}
	}

	return fetched, errors
}

// Push uploads nodes(that doesn't exists on given remote) from [s](current) to given [remote].
func (s *FirebaseService) Push(remote ServiceRepo) ([]models.Node, []error) {
	nodes, _, err := s.GetAll("", "", models.NotyaIgnoreFiles)
	if err != nil {
		return nil, []error{err}
	}

	errors := []error{}
	pushed := []models.Node{}

	for _, node := range nodes {
		exists, _ := remote.IsNodeExists(node)

		if node.IsFolder() && !exists {
			if _, err := remote.Mkdir(node.ToFolder()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}

			continue
		}

		r, _ := remote.View(node.ToNote())
		if !exists {
			if _, err := remote.Create(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}

			continue
		}

		if r.Body != node.Body {
			if _, err := remote.Edit(node.ToNote()); err != nil {
				errors = append(errors, assets.CannotDoSth("push", node.Title, err))
			} else {
				pushed = append(pushed, node)
			}
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
