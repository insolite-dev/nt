// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/anonistas/notya/lib/models"
)

// NotyaPWD, generates path of notya's notes directory.
// ╭───────────────────────╮   ╭────────╮   ╭────────────╮
// │ ~/user-home-directory │ + │ /notya │ = │ local path │
// ╰───────────────────────╯   ╰────────╯   ╰────────────╯
func NotyaPWD(settings models.Settings) (*string, error) {
	path := settings.LocalPath

	// Initialize default notya path.
	if len(path) == 0 || path == models.DefaultLocalPath {
		uhd, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		path = uhd + "/" + models.DefaultLocalPath
	}

	return &path, nil
}

// FileExists, checks if any type of file exists at given path.
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// WriteNote, creates new file and writes to its data.
// If file already exists at given path with same name, then it updates it's body.
//
// Could be used for create and edit.
func WriteNote(path string, body []byte) error {
	err := os.WriteFile(path, body, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// NewFolder, creates new empty working directory at given path(name).
func NewFolder(name string) error {
	if err := os.Mkdir(name, 0o750); err != nil {
		return err
	}

	return nil
}

// Delete, removes file or folder, from given path.
func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// ReadBody, opens file from given path, and takes its body to return.
func ReadBody(path string) (*string, error) {
	resbyte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res := string(resbyte)
	return &res, nil
}

// ListDir, reads all files from given-path directory.
func ListDir(path, prevPath, space string, ignore []string, tree bool) ([]string, []string, error) {
	// Read directory's files.
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	// Get path info.
	pathFI, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	// Convert list to string list.
	var res, pretty []string
	for _, d := range list {
		iIgn := false

		// Set i[index]Ignore variable
		for _, ig := range ignore {
			if ig == d.Name() {
				iIgn = true
			}
		}
		if iIgn {
			iIgn = false // Ignore file, set i[index]Ignore to false.
			continue
		}

		name := d.Name()
		var prettyName, localPrevPath string

		// Ignore base path's name.
		if len(space) > 2 {
			localPrevPath = prevPath
		}
		// Get full name via including ~/notya
		if len(space) > 0 {
			name = localPrevPath + pathFI.Name() + "/" + d.Name()
		}

		var subnames, subpretty []string
		if d.IsDir() {
			name += "/"
			prettyName = space + ""

			// Split sub-nodes as tree, eg:
			//
			//   My-Folder
			//     note.txt
			//     second_note.txt
			//     Sub-folder
			//       sub_note.txt
			//      ...
			if tree {
				n := models.Node{Title: pathFI.Name(), Path: path}
				sn, sp, err := ListDir(
					n.StructAsFolder().Path+d.Name(), n.StructAsFolder().Title,
					// Push to right by two empty charachters each sub node.
					space+"  ", ignore, tree,
				)
				if err != nil {
					return res, pretty, err
				}
				subnames = sn
				subpretty = sp
			}
		} else {
			prettyName = space + ""
		}

		res = append(res, name)
		pretty = append(pretty, prettyName)

		if len(subnames) > 0 && len(subpretty) > 0 {
			res = append(res, subnames...)
			pretty = append(pretty, subpretty...)
		}
	}

	return res, pretty, nil
}

// OpenViaEditor opens file in custom(appropriate from settings) from given path.
func OpenViaEditor(filepath string, stdargs models.StdArgs, settings models.Settings) error {
	// Look editor's execution path from current running machine.
	editor, pathErr := exec.LookPath(settings.Editor)
	if pathErr != nil {
		return pathErr
	}

	// Generate vi command to open file.
	editorCmd := &exec.Cmd{
		Path:   editor,
		Args:   []string{editor, filepath},
		Stdin:  stdargs.Stdin,
		Stdout: stdargs.Stdout,
		Stderr: stdargs.Stderr,
	}

	if err := editorCmd.Run(); err != nil {
		return err
	}

	return nil
}

// IsDir checks if the file (at provided [path]) is directory or not.
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err != nil {
		return false
	}

	return i.IsDir()
}
