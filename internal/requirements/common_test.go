package requirements_test

import (
	log2 "log"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/requirements"
	"go.arcalot.io/assert"
	"go.arcalot.io/log"
)

func TestLanguageRequirements(t *testing.T) {
	logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
	act, err := requirements.LanguageRequirements(".", []string{"dummy_plugin.py"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log2.Fatal(err)
	}
	assert.Equals(t, act, false)

	_, err = requirements.LanguageRequirements(".", []string{"dummy_plugin.rs"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	assert.Error(t, err)

	act, err = requirements.LanguageRequirements(".", []string{"dummy_plugin.go"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log2.Fatal(err)
	}
	assert.Equals(t, act, false)
}

func TestPluginLanguage(t *testing.T) {
	plugin_dir := []string{"plugin"}
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
		"d": {
			plugin_dir,
			"",
		},
		"e": {
			[]string{"plugin", "pyproject.toml"},
			"python",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			act, err := requirements.PluginLanguage(tc.filenames)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equals(t, tc.expectedResult, act)
		})
	}
}
