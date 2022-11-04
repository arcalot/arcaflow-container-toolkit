package requirements

import (
	"bytes"
	"fmt"
	"go.arcalot.io/log"
	"regexp"
)

func HasFilename(names []string, filename string) (bool, error) {
	for _, name := range names {
		matched, err := regexp.MatchString(filename, name)
		if err != nil {
			return false, fmt.Errorf("error when looking for %s and found %s (%w)", filename, name, err)
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func LanguageRequirements(abspath string, filenames []string, name string, version string, logger log.Logger,
	pythonCodeStyleChecker func(abspath string, stdout *bytes.Buffer, logger log.Logger) error) (bool, error) {
	meets_reqs := true
	lang, err := ImageLanguage(filenames)
	if err != nil {
		return false, err
	}
	switch lang {
	case "go":
		meets_reqs, err = GolangRequirements(filenames, logger)
		if err != nil {
			return false, err
		}
	case "python":
		meets_reqs, err = PythonRequirements(abspath, filenames, name, version, logger, pythonCodeStyleChecker)
		if err != nil {
			return false, err
		}
	default:
		logger.Infof("Programming Language %s not supported\n", lang)
		meets_reqs = false
	}

	return meets_reqs, nil
}
