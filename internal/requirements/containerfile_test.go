package requirements_test

import (
	log2 "log"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/requirements"
	"go.arcalot.io/assert"
	"go.arcalot.io/log/v2"
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
			act, err := requirements.ContainerfileRequirements(tc.path, logger)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equals(t, tc.expectedResult, act)
		})
	}
}
