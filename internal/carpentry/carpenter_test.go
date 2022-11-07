package carpentry

import (
	"bytes"
	"fmt"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/dto"
	"github.com/arcalot/arcaflow-plugin-image-builder/mocks/mock_ce_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	arcalog "go.arcalot.io/log"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func emptyPythonCodeStyle(abspath string, stdout *bytes.Buffer, logger arcalog.Logger) error {
	return nil
}

func textPythonCodeStyle(abspath string, stdout *bytes.Buffer, logger arcalog.Logger) error {
	_, err := stdout.WriteString("bad code")
	if err != nil {
		return err
	}
	return fmt.Errorf("code style error")
}

func TestBuildCmdMain(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineClient(ctrl)
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
	Carpentry(
		true, true, cec, conf, ".",
		python_filenames, logger, emptyPythonCodeStyle)
}

func TestAllTrue(t *testing.T) {
	a := make([]bool, 3)
	a[0] = true
	a[1] = false
	a[2] = true
	assert.False(t, AllTrue(a))

	a[1] = true
	assert.True(t, AllTrue(a))
}

func TestCliCarpentry(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	err := CliCarpentry(true, true, logger, "podman")
	var ecec *ErrorCEC
	assert.IsType(t, ecec, err)
}

func TestFlake8(t *testing.T) {
	stdout := &bytes.Buffer{}
	//logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	logger := arcalog.New(arcalog.Config{
		Level:       arcalog.LevelInfo,
		Destination: arcalog.DestinationStdout,
		Stdout:      os.Stdout,
	})
	err := flake8PythonCodeStyle("/githug/workplace", stdout, logger)
	var path_e_t *fs.PathError
	assert.IsType(t, path_e_t, err)

	afp, patherr := filepath.Abs("../../fixtures/pep8_compliant")
	if patherr != nil {
		log.Fatal(patherr)
	}
	fmt.Println(afp)
	err = flake8PythonCodeStyle(afp, stdout, logger)
	assert.Nil(t, err)

	afp, patherr = filepath.Abs("../../fixtures/pep8_non_compliant")
	if patherr != nil {
		log.Fatal(patherr)
	}
	stdout = &bytes.Buffer{}
	err = flake8PythonCodeStyle(afp, stdout, logger)
	var exec_exit_e_t *exec.ExitError
	assert.IsType(t, err, exec_exit_e_t)
}
