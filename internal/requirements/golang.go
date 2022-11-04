package requirements

import "go.arcalot.io/log"

func GolangRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true
	if has_, err := HasFilename(filenames, "go.mod"); err != nil {
		return false, err
	} else if !has_ {
		logger.Infof("Missing go.mod")
		meets_reqs = false
	}
	if has_, err := HasFilename(filenames, "go.sum"); err != nil {
		return false, err
	} else if !has_ {
		logger.Infof("Missing go.sum")
		meets_reqs = false
	}
	return meets_reqs, nil
}
