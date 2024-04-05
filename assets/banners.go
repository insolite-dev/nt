//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets

import "fmt"

// ShortSlog is a main slogan of nt.
var ShortSlog = "📝 Take notes quickly and expeditiously"

// MinimalisticBanner is a first banner of nt.
var MinimalisticBanner = `
  _   _   _____
 | \ | | |_   _|
 |  \| |   | |
 | |\  |   | |
 |_| \_|   |_| @ Insolite

`

// GenerateBanner merges slog and banner together, to get final result of application banner.
func GenerateBanner(banner string, slog string) string {
	template := `
  %v
 %v
   `

	return fmt.Sprintf(template, banner, slog)
}
