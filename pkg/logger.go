// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"fmt"
)

// Color and Icon of "CURRENT" output message.
var (
	Color string
	Icon  string
)

// level is a custom type of `string-level`.
// used define level for [setOutputStyle] function.
type level string

// Define constant app levels.
const (
	errorL   level = "error"
	successL level = "success"
	infoL    level = "info"
)

// Define constant color codes.
const (
	RED     string = "\033[0;31m"
	GREEN   string = "\033[0;32m"
	YELLOW  string = "\033[1;33m"
	NOCOLOR string = "\033[0m"
)

// outputLevel sets [Color] and [Icon] by given `level`,
// and then, returns final printable level title.
//
// Result cases:
// [ERROR] - (powered with red color)
// [SUCCESS] - (powered with green color)
// [INFO] - (powered with yellow color)
func outputLevel(l level) string {
	switch l {
	case errorL:
		Color = RED
		Icon = "[ERROR]"
	case successL:
		Color = GREEN
		Icon = "[OK]"
		break
	case infoL:
		Color = YELLOW
		Icon = "[INFO]"
	default:
		Color = NOCOLOR
	}

	return fmt.Sprintf("%s%s%s", Color, Icon, NOCOLOR)
}
