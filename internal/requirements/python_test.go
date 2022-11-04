package requirements

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	arcalog "go.arcalot.io/log"
	"log"
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

func TestPythonFileRequirements(t *testing.T) {
	min_correct := []string{"requirements.txt", "app.py", "main.py", "pyproject.toml"}
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
			true,
		},
		"c": {
			min_correct[2:],
			true,
		},
		"d": {
			min_correct[1:3],
			false,
		},
		"e": {
			[]string{},
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
