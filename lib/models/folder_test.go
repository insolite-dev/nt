//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models_test

import (
	"testing"

	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
)

func TestFolderToNode(t *testing.T) {
	tests := []struct {
		dir      models.Folder
		expected models.Node
	}{
		{
			dir:      models.Folder{},
			expected: models.Node{},
		},
		{
			dir:      models.Folder{Title: "folder/", Path: map[string]string{services.LOCAL.ToStr(): "~/folder"}},
			expected: models.Node{Title: "folder/", Path: map[string]string{services.LOCAL.ToStr(): "~/folder"}},
		},
	}

	for _, td := range tests {
		got := td.dir.ToNode()
		path := got.GetPath(services.LOCAL.ToStr())
		if got.Title != td.expected.Title || path != td.expected.GetPath(services.LOCAL.ToStr()) {
			t.Errorf("Sum was different of [Folder-To-Node] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
