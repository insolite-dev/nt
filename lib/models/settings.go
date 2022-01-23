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
	SettingsName     = ".settings.json"
	DefaultEditor    = "vi"
	DefaultLocalPath = "notya"
)

// Settings is a main structure model of application settings.
type Settings struct {
	Editor    string `json:"editor" default:"vi"`
	LocalPath string `json:"local_path" mapstructure:"local_path" survey:"local_path" default:"notya"`
}

// InitSettings returns default variant of settings structure model.
func InitSettings(localPath string) Settings {
	return Settings{
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
func FromJSON(value string) Settings {
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(value), &m)

	var s Settings
	mapstructure.Decode(m, &s)

	return s
}
