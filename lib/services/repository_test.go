//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package services_test

import (
	"testing"

	"github.com/insolite-dev/notya/lib/services"
)

func TestServiceTypeToStr(t *testing.T) {
	tests := []struct {
		t        *services.ServiceType
		expected string
	}{
		{t: &services.LOCAL, expected: "LOCAL"},
		{t: &services.FIRE, expected: "FIREBASE"},
		{t: nil, expected: "undefined"},
	}

	for _, td := range tests {
		got := td.t.ToStr()
		if got != td.expected {
			t.Errorf("ServiceTypeToStr sum is different, Want: %v | Got: %v", td.expected, got)
		}
	}
}
