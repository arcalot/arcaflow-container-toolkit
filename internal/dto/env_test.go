package dto_test

import (
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
	envvar_key := "foo"
	envvar_val := "bar"

	_, err := dto.LookupEnvVar(registries, envvar_key, logger)
	assert.Error(t, err)

	err = os.Setenv(envvar_key, envvar_val)
	if err != nil {
		log.Fatal(err)
	}
	v, err := dto.LookupEnvVar(registries, envvar_key, logger)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equals(t, v, envvar_val)

	envvar_val = "robot"
	err = os.Setenv(envvar_key, envvar_val)
	if err != nil {
		log.Fatal(err)
	}
	v, err = dto.LookupEnvVar(registries, envvar_key, logger)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equals(t, v, "robot")

	err = os.Unsetenv(envvar_key)
	if err != nil {
		log.Fatal(err)
	}
}
