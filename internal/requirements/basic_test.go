package requirements

import (
	"github.com/stretchr/testify/assert"
	"go.arcalot.io/log"
	log2 "log"
	"testing"
)

func TestBasicRequirements(t *testing.T) {
	min_correct := []string{"README.md", "Dockerfile", "plugin_test.py"}
	no_dockerfile := []string{"README.md", "plugin_test.py"}

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
			min_correct[:2],
			false,
		},
		"d": {
			no_dockerfile,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
			act, err := BasicRequirements(tc.filenames, logger)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}
