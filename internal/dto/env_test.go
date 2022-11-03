package dto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.arcalot.io/log"
	"os"
	"testing"
)

func TestLookupEnvVar(t *testing.T) {
	logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
	// these debug messages shouldn't be hard coded into this test
	envvar_key := "i_hope_this_isnt_used"
	envvar_val := ""
	type verbose struct {
		msg          string
		return_value string
	}

	v := LookupEnvVar(envvar_key, logger)
	assert.Equal(t, v.msg, fmt.Sprintf("%s not set", envvar_key))
	assert.Equal(t, v.return_value, "")

	os.Setenv(envvar_key, envvar_val)
	v = LookupEnvVar(envvar_key, logger)
	assert.Equal(t, v.msg, fmt.Sprintf("%s is empty", envvar_key))
	assert.Equal(t, v.return_value, "")

	envvar_val = "robot"
	os.Setenv(envvar_key, envvar_val)
	v = LookupEnvVar(envvar_key, logger)
	assert.Equal(t, v.msg, "")
	assert.Equal(t, v.return_value, envvar_val)

	os.Unsetenv(envvar_key)
}
