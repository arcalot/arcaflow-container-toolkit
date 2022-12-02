package requirements

import (
	"bytes"
	"fmt"
	"go.arcalot.io/log"
	"path/filepath"
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
	lang, err := PluginLanguage(filenames)
	if err != nil {
		return false, err
	}
	switch lang {
	case "go":
		return GolangRequirements(filenames, logger)
	case "python":
		return PythonRequirements(abspath, filenames, name, version, logger, pythonCodeStyleChecker)
	default:
		return false, fmt.Errorf("Programming Language %s not supported\n", lang)
	}
}

func PluginLanguage(filenames []string) (string, error) {
	ext2Lang := map[string]string{"go": "go", "py": "python"}
	cr := regexp.MustCompile(`(?i).*plugin.*\.`)
	for _, name := range filenames {
		if cr.MatchString(name) {
			ext := filepath.Ext(name)
			fmt.Println(name)
			lang, ok := ext2Lang[ext[1:]]
			if ok {
				return lang, nil
			}
		}
		if name == "pyproject.toml" {
			return "python", nil
		}
	}
	return "", nil
}
