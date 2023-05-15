package dto_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/dto"
	"go.arcalot.io/assert"
	arcalog "go.arcalot.io/log"
)

func TestLookupEnvVar(t *testing.T) {
	registries := "test"
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	envvar_key := "i_hope_this_isnt_used"
	envvar_val := ""

	v := dto.LookupEnvVar(registries, envvar_key, logger)
	assert.Equals(t, v, fmt.Sprintf("%s not set", envvar_key))

	err := os.Setenv(envvar_key, envvar_val)
	if err != nil {
		log.Fatal(err)
	}
	v = dto.LookupEnvVar(registries, envvar_key, logger)
	assert.Equals(t, v, fmt.Sprintf("%s is empty", envvar_key))

	envvar_val = "robot"
	err = os.Setenv(envvar_key, envvar_val)
	if err != nil {
		log.Fatal(err)
	}
	v = dto.LookupEnvVar(registries, envvar_key, logger)
	assert.Equals(t, v, "")

	err = os.Unsetenv(envvar_key)
	if err != nil {
		log.Fatal(err)
	}
}
