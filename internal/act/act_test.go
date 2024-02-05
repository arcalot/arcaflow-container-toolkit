package act_test

import (
	"bytes"
	"testing"

	act "go.arcalot.io/arcaflow-container-toolkit/internal/act"
	"go.arcalot.io/arcaflow-container-toolkit/internal/dto"
	mock_ces "go.arcalot.io/arcaflow-container-toolkit/mocks/ce_service"
	"go.arcalot.io/assert"
	arcalog "go.arcalot.io/log/v2"
	"go.uber.org/mock/gomock"
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
	conf := dto.ACT{
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
	passed, err := act.ACT(
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
	assert.Equals(t, act.AllTrue(a), false)
	a[1] = true
	assert.Equals(t, act.AllTrue(a), true)
}

func TestCliAct(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	assert.Error(t, act.CliACT(true, true, logger, "podman"))
}
