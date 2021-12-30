package assets_test

import (
	"testing"

	"github.com/anonistas/notya/assets"
)

func TestGenerateBanner(t *testing.T) {
	type args struct {
		slog   string
		banner string
	}

	tests := []struct {
		testname    string
		arguments   args
		expectedLen int
	}{
		{
			testname: "Default banner should be generated properly",
			arguments: args{
				assets.ShortSlog,
				assets.MinimalisticBanner,
			},
			expectedLen: 236,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.GenerateBanner(td.arguments.banner, td.arguments.slog)

			if len(got) != td.expectedLen {
				t.Errorf("Len of Sum was incorrect: Want: %v | Got: %v", td.expectedLen, len(got))
			}
		})
	}
}
