// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// All editors listed.
var (
	// VI related.
	VI     string = "vi"
	Vim    string = "vim"
	NeoVim string = "nvim"
	MacVim string = "mvim"
	GUIVim string = "gvim"

	// VS-Code related.
	VSCode         string = "code"
	VSCodeInsiders string = "code-insiders"
)

// Constant values of settings.
const (
	SettingsName = ".settings.json"
)

// Settings is a main structure model of application settings.
type Settings struct {
	Editor string `json:"editor" default:"vi"`
}

// InitSettings returns default variant of settings structure model.
func InitSettings() Settings {
	return Settings{Editor: VI}
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
