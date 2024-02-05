package requirements

import "go.arcalot.io/log/v2"

func BasicRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true

	if present, err := HasFilename(filenames, "README.md"); err != nil {
		return false, err
	} else if !present {
		logger.Errorf("Missing required file README.md")
		meets_reqs = false
	}

	if present, err := HasFilename(filenames, "Dockerfile"); err != nil {
		return false, err
	} else if !present {
		logger.Errorf("Missing required file Dockerfile")
		meets_reqs = false
	}

	if present, err := HasFilename(filenames, "(?i).*test.*"); err != nil {
		return false, err
	} else if !present {
		// match case-insensitive 'test'?
		logger.Errorf("Missing required test file")
		meets_reqs = false
	}
	return meets_reqs, nil
}
