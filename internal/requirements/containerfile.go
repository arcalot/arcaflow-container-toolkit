package requirements

import (
	"go.arcalot.io/log"
	"os"
	"path/filepath"
	"regexp"
)

func ContainerfileRequirements(abspath string, logger log.Logger) (bool, error) {
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

		// create map of regexp patterns to search for in Dockerfile as well as log information if not found
		m := map[string]string{
			"FROM quay\\.io/(centos/centos:stream8|arcalot/arcaflow-plugin-baseimage-python-.*base$)": "Dockerfile doesn't use a supported base image\n",
			"(ADD|COPY) .*/?LICENSE /.*":                                                                "Dockerfile does not contain copy of arcaflow plugin license\n",
			"CMD \\[\\]":                                                                                "Dockerfile does not contain an empty command (i.e. CMD [])",
			"LABEL org.opencontainers.image.source=\".*\"":                                              "Dockerfile is missing LABEL org.opencontainers.image.source",
			"LABEL org.opencontainers.image.licenses=\"Apache-2\\.0.*\"":                                "Dockerfile is missing LABEL org.opencontainers.image.licenses",
			"LABEL org.opencontainers.image.vendor=\"Arcalot project\"":                                 "Dockerfile is missing LABEL org.opencontainers.image.vendor",
			"LABEL org.opencontainers.image.authors=\"Arcalot contributors\"":                           "Dockerfile is missing LABEL org.opencontainers.image.authors",
			"LABEL org.opencontainers.image.title=\".*\"":                                               "Dockerfile is missing LABEL org.opencontainers.image.title",
			"LABEL io.github.arcalot.arcaflow.plugin.version=\"(\\d*)(\\.?\\d*?)(\\.?\\d*?)\"":          "Dockerfile is missing LABEL io.github.arcalot.arcaflow.plugin.version that uses form X, X.Y, X.Y.Z",
		}

		for regexp_, loggerResp := range m {
			if has_, err := DockerfileHasLine(dockerfile, regexp_); err != nil {
				return false, err
			} else if !has_ {
				logger.Infof(loggerResp)
				meets_reqs = has_
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
