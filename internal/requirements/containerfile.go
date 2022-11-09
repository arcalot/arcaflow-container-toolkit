package requirements

import (
	"go.arcalot.io/log"
	"os"
	"path/filepath"
	"regexp"
)

func ContainerfileRequirements(abspath string, logger log.Logger) (bool, error) {
	meets_reqs := true
	project_files, err := os.Open(abspath)
	if err != nil {
		return false, err
	}
	defer project_files.Close()
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		return false, err
	}
	has_, err := HasFilename(filenames, "Dockerfile")
	if err != nil {
		return false, err
	}
	if !has_ {
		logger.Infof("Missing Dockerfile")
		meets_reqs = false

	} else {
		file, err := os.ReadFile(filepath.Join(abspath, "Dockerfile"))
		if err != nil {
			return false, err
		}
		dockerfile := string(file)

		// create map of regexp patterns to search for in Dockerfile as well as log information if not found
		m := map[string]string{
			"FROM quay\\.io/centos/centos:stream8":                                             "Dockerfile doesn't use 'FROM quay.io/centos/centos:stream8'\n",
			"(ADD|COPY) .*/LICENSE /.*":                                                        "Dockerfile does not contain copy of arcaflow plugin license\n",
			"ENTRYPOINT \\[.*\".*plugin.*\".*\\]":                                              "Dockerfile enterypoint does not point to an executable that includes 'plugin' in its name",
			"CMD \\[\\]":                                                                       "Dockerfile does not contain an empty command (i.e. CMD [])",
			"LABEL org.opencontainers.image.source=\".*\"":                                     "Dockerfile is missing LABEL org.opencontainers.image.source",
			"LABEL org.opencontainers.image.licenses=\"Apache-2\\.0.*\"":                       "Dockerfile is missing LABEL org.opencontainers.image.licenses",
			"LABEL org.opencontainers.image.vendor=\"Arcalot project\"":                        "Dockerfile is missing LABEL org.opencontainers.image.vendor",
			"LABEL org.opencontainers.image.authors=\"Arcalot contributors\"":                  "Dockerfile is missing LABEL org.opencontainers.image.authors",
			"LABEL org.opencontainers.image.title=\".*\"":                                      "Dockerfile is missing LABEL org.opencontainers.image.title",
			"LABEL io.github.arcalot.arcaflow.plugin.version=\"(\\d*)(\\.?\\d*?)(\\.?\\d*?)\"": "Dockerfile is missing LABEL io.github.arcalot.arcaflow.plugin.version that uses form X, X.Y, X.Y.Z",
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

func ImageLanguage(filenames []string) (string, error) {
	ext2Lang := map[string]string{"go": "go", "py": "python"}
	cr, err := regexp.Compile("(?i).*plugin.*")
	if err != nil {
		return "", err
	}
	for _, name := range filenames {
		if cr.MatchString(name) {
			ext := filepath.Ext(name)
			lang, ok := ext2Lang[ext[1:]]
			if ok {
				return lang, nil
			}
		}
	}
	return "", nil
}

func DockerfileHasLine(dockerfile string, line string) (bool, error) {
	matched, err := regexp.MatchString(line, dockerfile)
	if err != nil {
		return false, err
	}
	return matched, nil
}
