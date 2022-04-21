// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

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
}

// Settings is a main structure model of application settings.
//
//  Example
// ╭────────────────────────────────────────────────────╮
// │ Editor: vi                                         │
// │ Local Path: /User/random-user/notya/.settings.json │
// ╰────────────────────────────────────────────────────╯
type Settings struct {
	Name               string `json:"name" default:"notya"`
	Editor             string `json:"editor" default:"vi"`
	LocalPath          string `json:"local_path" mapstructure:"local_path" survey:"local_path"`
	FirebaseAccountKey string `json:"firebase" mapstructure:"firebase"`
}

// InitSettings returns default variant of settings structure model.
func InitSettings(localPath string) Settings {
	return Settings{
		Name:      DefaultAppName,
		Editor:    DefaultEditor,
		LocalPath: localPath,
	}
}

// ToByte converts settings model to JSON map,
// but returns that value as byte array.
func (s *Settings) ToByte() []byte {
	b, _ := json.Marshal(&s)

	var j map[string]interface{}
	_ = json.Unmarshal(b, &j)

	res, _ := json.Marshal(&j)

	return res
}

// FromJSON converts string(map) value to Settings structure.
func DecodeSettings(value string) Settings {
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(value), &m)

	var s Settings
	mapstructure.Decode(m, &s)

	return s
}

// IsValid checks validness of settings structure.
func (s *Settings) IsValid() bool {
	return len(s.Name) > 0 && len(s.Editor) > 0 && len(s.LocalPath) > 0
}

func IsUpdated(old, current Settings) bool {
	return old.Editor != current.Editor || old.LocalPath != current.LocalPath
}

func IsPathUpdated(old, current Settings) bool {
	return old.LocalPath != current.LocalPath
}
