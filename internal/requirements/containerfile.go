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
		meets_reqs = false
	} else {
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

		labelChecks := []labelCheck{
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

		for _, lc := range labelChecks {
			if labelValidation == "lenient" {
				if has_, err := DockerfileHasLine(dockerfile, lc.existsPattern); err != nil {
					return false, err
				} else if !has_ {
					logger.Warningf("Dockerfile is missing LABEL %s", lc.labelName)
				}
			} else {
				if has_, err := DockerfileHasLine(dockerfile, lc.strictPattern); err != nil {
					return false, err
				} else if !has_ {
					logger.Infof("Dockerfile is missing LABEL %s", lc.labelName)
					meets_reqs = has_
				}
			}
		}
	}
	return meets_reqs, nil
}

func DockerfileHasLine(dockerfile string, line string) (bool, error) {
	matched, err := regexp.MatchString(line, dockerfile)
	if err != nil {
		return false, err
	}
	return matched, nil
}
