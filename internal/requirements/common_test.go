package requirements

import (
	"github.com/stretchr/testify/assert"
	"go.arcalot.io/log"
	log2 "log"
	"testing"
)

func TestLanguageRequirements(t *testing.T) {
	logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
	act, err := LanguageRequirements(".", []string{"dummy_plugin.py"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log2.Fatal(err)
	}
	assert.False(t, act)

	act, err = LanguageRequirements(".", []string{"dummy_plugin.rs"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log2.Fatal(err)
	}
	assert.False(t, act)

	act, err = LanguageRequirements(".", []string{"dummy_plugin.go"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log2.Fatal(err)
	}
	assert.False(t, act)
}
