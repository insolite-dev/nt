//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// Constant values of settings.
const (
	DefaultAppName   = "notya"
	SettingsName     = ".settings.json"
	DefaultEditor    = "vi"
	DefaultLocalPath = "notya"
)

// NotyaIgnoreFiles are those files that shouldn't
// be represented as note files.
var NotyaIgnoreFiles []string = []string{
	SettingsName,
	".DS_Store", // Darwin related.
	".git",
}

// Settings is a main structure model of application settings.
//
//	Example:
//
// ╭────────────────────────────────────────────────────╮
// │ Name: notya                                        │
// │ Editor: vi                                         │
// │ Notes Path: /User/random-user/notya/notes          │
// │ Firebase Project ID: notya-98tf3                   │
// │ Firebase Account Key: /User/.../notya/key.json     │
// │ Firebase Collection: notya-notes                   │
// ╰────────────────────────────────────────────────────╯
type Settings struct {
	// Alert: development related field, shouldn't be used in production.
	ID string `json:",omitempty"`

	// The custom name of your notya application.
	Name string `json:"name" default:"notya"`

	// Editor app of application.
	// Could be:
	//   - vi
	//   - vim
	//   - nvim
	//   - code
	//   - code-insiders
	//  and etc. shortly each code editor that could be opened by its command.
	//  like: `code .` or `nvim .`.
	Editor string `json:"editor" default:"vi"`

	// Local "notes" folder path for notes, independently from [~/notya/] folder.
	// Must be given full path, like: "./User/john-doe/.../my-notya-notes/"
	//
	// Does same job as [FirebaseCollection] for local env.
	NotesPath string `json:"notes_path" mapstructure:"notes_path" survey:"notes_path"`

	// The project id of your firebase project.
	//
	// It is required for firebase remote connection.
	FirebaseProjectID string `json:"fire_project_id,omitempty" mapstructure:"fire_project_id,omitempty" survey:"fire_project_id"`

	// The path of key of "firebase-service" account file.
	// Must be given full path, like: "./User/john-doe/.../..."
	//
	// It is required for firebase remote connection.
	FirebaseAccountKey string `json:"fire_account_key,omitempty" mapstructure:"fire_account_key,omitempty" survey:"fire_account_key"`

	// The concrete collection of nodes.
	// Does same job as [NotesPath] but has to take just name of collection.
	FirebaseCollection string `json:"fire_collection,omitempty" mapstructure:"fire_collection,omitempty" survey:"fire_collection"`
}

// CopyWith updates pointed settings with a new data.
// if given argument is not nil, it will be overwritten
// inside pointed settings model.
func (s *Settings) CopyWith(
	ID *string,
	Name *string,
	Editor *string,
	NotesPath *string,
	FirebaseProjectID *string,
	FirebaseAccountKey *string,
	FirebaseCollection *string,
) Settings {
	ss := *s

	if ID != nil {
		ss.ID = *ID
	}
	if Name != nil {
		ss.Name = *Name
	}
	if Editor != nil {
		ss.Editor = *Editor
	}
	if NotesPath != nil {
		ss.NotesPath = *NotesPath
	}
	if FirebaseProjectID != nil {
		ss.FirebaseProjectID = *FirebaseProjectID
	}
	if FirebaseAccountKey != nil {
		ss.FirebaseAccountKey = *FirebaseAccountKey
	}
	if FirebaseCollection != nil {
		ss.FirebaseCollection = *FirebaseCollection
	}

	return ss
}

// InitSettings returns default variant of settings structure model.
func InitSettings(notesPath string) Settings {
	return Settings{
		Name:      DefaultAppName,
		Editor:    DefaultEditor,
		NotesPath: notesPath,
	}
}

// ToString converts settings model to a formatted JSON string,
func (s *Settings) ToString() string {
	jsonBytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// ToJSON converts string structure model to map value.
func (s *Settings) ToJSON() map[string]interface{} {
	b, _ := json.Marshal(&s)

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	return m
}

// DecodeSettings converts string(map) value to Settings structure.
func DecodeSettings(value string) Settings {
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(value), &m)

	var s Settings
	mapstructure.Decode(m, &s)

	return s
}

// FirePath returns valid firebase collection name.
func (s *Settings) FirePath() string {
	if len(s.FirebaseCollection) > 0 {
		return s.FirebaseCollection
	} else if len(s.Name) > 0 {
		return s.Name
	}

	return DefaultAppName
}

// IsValid checks validness of settings structure.
func (s *Settings) IsValid() bool {
	return len(s.Name) > 0 && len(s.Editor) > 0 && len(s.NotesPath) > 0
}
