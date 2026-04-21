package requirements

import (
	"os"
	"path/filepath"
	"regexp"

	"go.arcalot.io/log/v2"
)

type labelCheck struct {
	strictPattern string
	existsPattern string
	labelName     string
}

var labelChecks = []labelCheck{
	{
		strictPattern: `LABEL org.opencontainers.image.source=".*"`,
		existsPattern: `LABEL org.opencontainers.image.source=`,
		labelName:     "org.opencontainers.image.source",
	},
	{
		strictPattern: `LABEL org.opencontainers.image.licenses="Apache-2\.0.*"`,
		existsPattern: `LABEL org.opencontainers.image.licenses=`,
		labelName:     "org.opencontainers.image.licenses",
	},
	{
		strictPattern: `LABEL org.opencontainers.image.vendor="Arcalot project"`,
		existsPattern: `LABEL org.opencontainers.image.vendor=`,
		labelName:     "org.opencontainers.image.vendor",
	},
	{
		strictPattern: `LABEL org.opencontainers.image.authors="Arcalot contributors"`,
		existsPattern: `LABEL org.opencontainers.image.authors=`,
		labelName:     "org.opencontainers.image.authors",
	},
	{
		strictPattern: `LABEL org.opencontainers.image.title=".*"`,
		existsPattern: `LABEL org.opencontainers.image.title=`,
		labelName:     "org.opencontainers.image.title",
	},
	{
		strictPattern: `LABEL io.github.arcalot.arcaflow.plugin.version="(\d*)(\.?\d*?)(\.?\d*?)"`,
		existsPattern: `LABEL io.github.arcalot.arcaflow.plugin.version=`,
		labelName:     "io.github.arcalot.arcaflow.plugin.version",
	},
}

func ContainerfileRequirements(abspath string, labelValidation string, logger log.Logger) (bool, error) {
	meets_reqs := true
	project_files, err := os.Open(filepath.Clean(abspath))
	if err != nil {
		return false, err
	}
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		return false, err
	}
	err = project_files.Close()
	if err != nil {
		return false, err
	}
	present, err := HasFilename(filenames, "Dockerfile")
	if err != nil {
		return false, err
	}
	if !present {
		logger.Infof("Missing Dockerfile")
		return false, nil
	}
	file, err := os.ReadFile(filepath.Join(filepath.Clean(abspath), "Dockerfile"))
	if err != nil {
		return false, err
	}
	dockerfile := string(file)

	nonLabelChecks := map[string]string{
		"FROM quay\\.io/(centos/centos:stream8|arcalot/arcaflow-plugin-baseimage-python-.*base)": "Dockerfile doesn't use a supported base image\n",
		"(ADD|COPY) .*/?LICENSE /.*": "Dockerfile does not contain copy of arcaflow plugin license\n",
		"CMD \\[\\]":                 "Dockerfile does not contain an empty command (i.e. CMD [])",
	}

	for regexp_, loggerResp := range nonLabelChecks {
		if has_, err := DockerfileHasLine(dockerfile, regexp_); err != nil {
			return false, err
		} else if !has_ {
			logger.Infof(loggerResp)
			meets_reqs = has_
		}
	}

	labelsOK, err := checkLabels(dockerfile, labelValidation, logger)
	if err != nil {
		return false, err
	}
	if !labelsOK {
		meets_reqs = false
	}

	return meets_reqs, nil
}

func checkLabels(dockerfile string, labelValidation string, logger log.Logger) (bool, error) {
	if labelValidation == "lenient" {
		return checkLabelsLenient(dockerfile, logger)
	}
	return checkLabelsStrict(dockerfile, logger)
}

func checkLabelsLenient(dockerfile string, logger log.Logger) (bool, error) {
	for _, lc := range labelChecks {
		has, err := DockerfileHasLine(dockerfile, lc.existsPattern)
		if err != nil {
			return false, err
		}
		if !has {
			logger.Warningf("Dockerfile is missing LABEL %s", lc.labelName)
		}
	}
	return true, nil
}

func checkLabelsStrict(dockerfile string, logger log.Logger) (bool, error) {
	passed := true
	for _, lc := range labelChecks {
		has, err := DockerfileHasLine(dockerfile, lc.strictPattern)
		if err != nil {
			return false, err
		}
		if !has {
			logger.Infof("Dockerfile is missing LABEL %s", lc.labelName)
			passed = false
		}
	}
	return passed, nil
}

func DockerfileHasLine(dockerfile string, line string) (bool, error) {
	matched, err := regexp.MatchString(line, dockerfile)
	if err != nil {
		return false, err
	}
	return matched, nil
}
