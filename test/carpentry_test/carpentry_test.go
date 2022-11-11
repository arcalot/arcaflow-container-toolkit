package carpentry_test

import (
	"bytes"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/carpentry"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/dto"
	mock_ces "github.com/arcalot/arcaflow-plugin-image-builder/mocks/ce_service"
	"github.com/golang/mock/gomock"

	"go.arcalot.io/assert"
	arcalog "go.arcalot.io/log"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func emptyPythonCodeStyle(abspath string, stdout *bytes.Buffer, logger arcalog.Logger) error {
	return nil
}

func TestBuildCmdMain(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mock_ces.NewMockContainerEngineService(ctrl)
	rg1 := dto.Registry{
		Url:      "reg1.io",
		Username: "user1",
		Password: "secret1",
	}
	rg2 := dto.Registry{
		Url:      "reg2.io",
		Username: "user2",
		Password: "secret2",
	}
	conf := dto.Carpenter{
		Revision:         "20220928",
		Image_Name:       "dummy",
		Image_Tag:        "latest",
		Project_Filepath: ".",
		Registries:       []dto.Registry{rg1, rg2},
	}
	python_filenames := []string{
		"plugin.py",
		"test_plugin.py",
		"Dockerfile",
		"requirements.txt",
		"pyproject.toml"}
	passed, err := carpentry.Carpentry(
		true, true, cec, conf, ".",
		python_filenames, logger, emptyPythonCodeStyle)
	assert.Equals(t, passed, false)
	assert.NoError(t, err)
}

func TestAllTrue(t *testing.T) {
	a := make([]bool, 3)
	a[0] = true
	a[1] = false
	a[2] = true
	assert.Equals(t, carpentry.AllTrue(a), false)
	a[1] = true
	assert.Equals(t, carpentry.AllTrue(a), true)
}

func TestCliCarpentry(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	assert.Error(t, carpentry.CliCarpentry(true, true, logger, "podman"))
}

func TestFlake8(t *testing.T) {
	stdout := &bytes.Buffer{}
	logger := arcalog.New(arcalog.Config{
		Level:       arcalog.LevelInfo,
		Destination: arcalog.DestinationStdout,
		Stdout:      os.Stdout,
	})
	err := carpentry.Flake8PythonCodeStyle("/githug/workplace", stdout, logger)
	assert.Error(t, err)

	afp, patherr := filepath.Abs("../../fixtures/pep8_compliant")
	if patherr != nil {
		log.Fatal(patherr)
	}
	assert.Nil(t, carpentry.Flake8PythonCodeStyle(afp, &bytes.Buffer{}, logger))

	afp, patherr = filepath.Abs("../../fixtures/pep8_non_compliant")
	if patherr != nil {
		log.Fatal(patherr)
	}
	assert.Error(t, carpentry.Flake8PythonCodeStyle(afp, &bytes.Buffer{}, logger))
}
