package requirements_test

import (
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/requirements"
	"go.arcalot.io/assert"
	"go.arcalot.io/log"
	log2 "log"
	"testing"
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
