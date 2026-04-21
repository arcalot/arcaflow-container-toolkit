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
		path            string
		labelValidation string
		expectedResult  bool
	}{
		"good_dockerfile_strict": {
			"../../fixtures/perfect",
			"strict",
			true,
		},
		"bad_dockerfile_strict": {
			"../../fixtures/no_good",
			"strict",
			false,
		},
		"good_dockerfile_lenient": {
			"../../fixtures/perfect",
			"lenient",
			true,
		},
		"bad_dockerfile_lenient": {
			"../../fixtures/no_good",
			"lenient",
			false,
		},
		"no_labels_lenient": {
			"../../fixtures/no_labels",
			"lenient",
			true,
		},
		"no_labels_strict": {
			"../../fixtures/no_labels",
			"strict",
			false,
		},
		"wrong_label_values_strict": {
			"../../fixtures/wrong_label_values",
			"strict",
			false,
		},
		"wrong_label_values_lenient": {
			"../../fixtures/wrong_label_values",
			"lenient",
			true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
			act, err := requirements.ContainerfileRequirements(tc.path, tc.labelValidation, logger)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equals(t, tc.expectedResult, act)
		})
	}
}
