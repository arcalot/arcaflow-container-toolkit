package requirements

import "go.arcalot.io/log"

func GolangRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true
	if present, err := HasFilename(filenames, "go.mod"); err != nil {
		return false, err
	} else if !present {
		logger.Infof("Missing go.mod")
		meets_reqs = false
	}
	if present, err := HasFilename(filenames, "go.sum"); err != nil {
		return false, err
	} else if !present {
		logger.Infof("Missing go.sum")
		meets_reqs = false
	}
	return meets_reqs, nil
}
