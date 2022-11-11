package requirements_test

import (
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/requirements"
	"go.arcalot.io/assert"
	"go.arcalot.io/log"
	log2 "log"
	"testing"
)

func TestGolangRequirements(t *testing.T) {
	min_correct := []string{"go.mod", "go.sum"}
	testCases := map[string]struct {
		filenames      []string
		expectedResult bool
	}{
		"a": {
			min_correct,
			true,
		},
		"b": {
			min_correct[1:],
			false,
		},
		"c": {
			min_correct[:1],
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
			act, err := requirements.GolangRequirements(tc.filenames, logger)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equals(t, tc.expectedResult, act)
		})
	}
}
