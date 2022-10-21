package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"testing"

	mocks "github.com/arcalot/arcaflow-plugin-image-builder/mocks/mock_ce_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	arcalog "go.arcalot.io/log"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestIntMinBasic(t *testing.T) {
	ans := IntMin(2, -2)
	if ans != -2 {
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
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

func TestUserIsQuayRobot(t *testing.T) {
	testCases := map[string]struct {
		username       string
		expectedResult bool
	}{
		"a": {
			"river+robot",
			true,
		},
		"b": {
			"river+",
			false,
		},
		"c": {
			"river",
			false,
		},
		"d": {
			"+robot",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			act, err := UserIsQuayRobot(tc.username)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestImageLanguage(t *testing.T) {
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
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			act, err := ImageLanguage(tc.filenames)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestPythonFileRequirements(t *testing.T) {
	min_correct := []string{"requirements.txt", "app.py", "pyproject.toml"}
	testCases := map[string]struct {
		filenames      []string
		expectedResult bool
	}{
		"a": {
			min_correct,
			true,
		},
		"b": {
			min_correct[:1],
			false,
		},
		"c": {
			min_correct[2:],
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
			act, err := PythonFileRequirements(tc.filenames, logger)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})

	}

}

func TestBasicRequirements(t *testing.T) {
	min_correct := []string{"README.md", "Dockerfile", "plugin_test.py"}
	no_dockerfile := []string{"README.md", "plugin_test.py"}

	testCases := map[string]struct {
		filenames      []string
		expectedResult bool
	}{
		"a": {
			min_correct,
			true,
		},
		"b": {
			min_correct[1:],
			false,
		},
		"c": {
			min_correct[:2],
			false,
		},
		"d": {
			no_dockerfile,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
			act, err := BasicRequirements(tc.filenames, logger)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestFilterByIndex(t *testing.T) {
	a := Registry{Url: "a"}
	b := Registry{Url: "b"}
	c := Registry{Url: "c"}
	d := Registry{Url: "d"}
	e := Registry{Url: "e"}
	var PlaceHolder struct{}
	list := []Registry{a, b, c, d, e}
	remove := map[string]Empty{
		"1": PlaceHolder,
		"3": PlaceHolder,
	}
	actualList := FilterByIndex(list, remove)
	assert.Equal(t, actualList[0], a)
	assert.Equal(t, actualList[1], c)
	assert.Equal(t, actualList[2], e)
}

func TestContainerRequirements(t *testing.T) {
	testCases := map[string]struct {
		path           string
		expectedResult bool
	}{
		"good_dockerfile": {
			"../fixtures/perfect",
			true,
		},
		"bad_dockerfile": {
			"../fixtures/no_good",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
			act, err := ContainerRequirements(
				tc.path, "dummy", "latest", logger)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestGolangRequirements(t *testing.T) {
	min_correct := []string{"go.mod", "go.sum"}
	testCases := map[string]struct {
		filenames      []string
		expectedResult bool
	}{
		"a": {
			min_correct,
			true,
		},
		"b": {
			min_correct[1:],
			false,
		},
		"c": {
			min_correct[:1],
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
			act, err := GolangRequirements(tc.filenames, logger)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

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

func TestPythonCodeStyle(t *testing.T) {

	testCases := map[string]struct {
		fn             func(string, *bytes.Buffer, arcalog.Logger) error
		expectedResult bool
	}{
		"a": {
			emptyPythonCodeStyle,
			true,
		},
		"b": {
			textPythonCodeStyle,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
			act, err := PythonCodeStyle(".", "dummy", "latest", tc.fn, logger)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}

func TestLanguageRequirements(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	act, err := LanguageRequirements(".", []string{"dummy_plugin.py"}, "dummy",
		"latest", logger, emptyPythonCodeStyle)
	if err != nil {
		log.Fatal(err)
	}
	assert.False(t, act)
}

func TestLookupEnvVar(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
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

func TestBuildImage(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineClient(ctrl)
	cec.EXPECT().
		Build("use", "the", []string{"forks"}).
		Return(nil).
		Times(1)
	BuildImage(true, true, cec, "use", "the", "forks", logger)
}

func TestBuildCmdMain(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineClient(ctrl)
	rg1 := Registry{
		Url:      "reg1.io",
		Username: "user1",
		Password: "secret1",
	}
	rg2 := Registry{
		Url:      "reg2.io",
		Username: "user2",
		Password: "secret2",
	}
	conf := config{
		Revision:         "20220928",
		Image_Name:       "dummy",
		Image_Tag:        "latest",
		Project_Filepath: ".",
		Registries:       []Registry{rg1, rg2},
	}
	python_filenames := []string{
		"plugin.py",
		"test_plugin.py",
		"Dockerfile",
		"requirements.txt",
		"pyproject.toml"}
	BuildCmdMain(
		true, true, cec, conf, ".",
		python_filenames, logger, emptyPythonCodeStyle)
}

func TestPushImage(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineClient(ctrl)
	rg1 := Registry{
		Url:       "reg1.io",
		Username:  "user1",
		Password:  "secret1",
		Namespace: "allyourbases",
	}
	image_name := "usethe"
	image_tag := "forks"

	destination := fmt.Sprintf("%s/%s/%s:%s", rg1.Url, rg1.Namespace, image_name, image_tag)
	cec.EXPECT().
		Tag(fmt.Sprintf("%s:%s", image_name, image_tag), destination).
		Return(nil).
		Times(1)
	cec.EXPECT().
		Push(destination, rg1.Username, rg1.Password, rg1.Url).
		Return(nil).
		Times(1)
	PushImage(true, true, true, cec, image_name, image_tag,
		rg1.Username, rg1.Password, rg1.Url, rg1.Namespace, logger)
}
