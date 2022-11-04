package requirements

import (
	"github.com/stretchr/testify/assert"
	"go.arcalot.io/log"
	log2 "log"
	"testing"
)

func TestContainerRequirements(t *testing.T) {
	testCases := map[string]struct {
		path           string
		expectedResult bool
	}{
		"good_dockerfile": {
			"../../fixtures/perfect",
			true,
		},
		"bad_dockerfile": {
			"../../fixtures/no_good",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
			act, err := ContainerfileRequirements(tc.path, logger)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestImageLanguage(t *testing.T) {
	python_file := []string{"plugin.py"}
	golang_file := []string{"plugin.go"}

	testCases := map[string]struct {
		filenames      []string
		expectedResult string
	}{
		"a": {
			python_file,
			"python",
		},
		"b": {
			golang_file,
			"go",
		},
		"c": {
			[]string{},
			"",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			act, err := ImageLanguage(tc.filenames)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}
