/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/dto"
	golog "log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/arcalot/arcaflow-plugin-image-builder/internal/ce_client"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/util"
	"github.com/spf13/cobra"
	"go.arcalot.io/log"
)

var Push bool
var Build bool

type ExternalProgramOnFile func(executable_filepath string, stdout *bytes.Buffer, logger log.Logger) error

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.PersistentFlags().BoolVarP(&Push, "push", "p", false, "push images to registry")
	buildCmd.PersistentFlags().BoolVarP(&Build, "build", "b", false, "validate requirements and build image")
}

var buildCmd = &cobra.Command{
	Use:   "build an image",
	Short: "build image",
	Run: func(cmd *cobra.Command, args []string) {

		cec, err := ce_client.NewCeClient("docker")

		if err != nil {
			rootLogger.Errorf("invalid container engine client (%w)", err)
			panic(err)
		}
		conf, err := dto.Unmarshal(rootLogger)
		if err != nil {
			rootLogger.Errorf("invalid carpenter config (%w)", err)
			panic(err)
		}
		abspath, err := filepath.Abs(conf.Project_Filepath)
		if err != nil {
			rootLogger.Errorf("invalid absolute path to project (%w)", err)
			panic(err)
		}
		files, err := os.Open(abspath)
		if err != nil {
			rootLogger.Errorf("error opening project directory (%w)", err)
			panic(err)
		}
		filenames, err := files.Readdirnames(0)
		if err != nil {
			rootLogger.Errorf("error reading project directory (%w)", err)
			panic(err)
		}
		err = files.Close()
		if err != nil {
			rootLogger.Errorf("error closing directory at %s (%w)", abspath, err)
			panic(err)
		}
		passed_reqs, err := BuildCmdMain(Build, Push, cec, conf, abspath, filenames,
			rootLogger,
			flake8PythonCodeStyle)
		if err != nil {
			panic(err)
		}
		if !passed_reqs {
			golog.Fatalf("failed requirements check, not building: %s %s", conf.Image_Name, conf.Image_Tag)
		}
	},
}

func BuildCmdMain(build_img bool, push_img bool, cec ce_client.ContainerEngineClient, conf dto.Carpenter, abspath string,
	filenames []string, logger log.Logger,
	pythonCodeStyleChecker func(abspath string, stdout *bytes.Buffer, logger log.Logger) error) (bool, error) {

	meets_reqs := make([]bool, 3)
	basic_reqs, err := BasicRequirements(filenames, logger)
	if err != nil {
		return false, err
	}
	meets_reqs[0] = basic_reqs
	container_reqs, err := ContainerRequirements(abspath, conf.Image_Name, conf.Image_Tag, logger)
	if err != nil {
		return false, err
	}
	meets_reqs[1] = container_reqs
	lang_req, err := LanguageRequirements(abspath, filenames, conf.Image_Name, conf.Image_Tag, logger,
		pythonCodeStyleChecker)
	if err != nil {
		return false, err
	}
	meets_reqs[2] = lang_req
	all_checks := AllTrue(meets_reqs)
	if !all_checks {
		return false, nil
	}
	if err := BuildImage(build_img, all_checks, cec, abspath, conf.Image_Name, conf.Image_Tag,
		logger); err != nil {
		return false, err
	}
	for _, registry := range conf.Registries {
		if err := PushImage(all_checks, build_img, push_img, cec, conf.Image_Name, conf.Image_Tag,
			registry.Username, registry.Password, registry.Url, registry.Namespace, logger); err != nil {
			logger.Errorf("(%w)", err)
		}
	}
	return true, nil
}

func BuildImage(build_img bool, all_checks bool, cec ce_client.ContainerEngineClient, abspath string, image_name string,
	image_tag string, logger log.Logger) error {

	if all_checks && build_img {
		logger.Infof("Passed all requirements: %s %s\n", image_name, image_tag)
		logger.Infof("Building %s %s from %v\n", image_name, image_tag, abspath)
		if err := cec.Build(abspath, image_name, []string{image_tag}); err != nil {
			return err
		}
	}
	return nil
}

func PushImage(all_checks, build_image, push_image bool, cec ce_client.ContainerEngineClient, name, version,
	username, password, registry_address, registry_namespace string, logger log.Logger) error {

	if all_checks && build_image && push_image {
		image_name_tag := name + ":" + version
		destination := strings.Join(
			[]string{registry_address, registry_namespace, image_name_tag},
			"/")
		logger.Infof("Pushing %s to %s", name, destination)
		err := cec.Tag(image_name_tag, destination)
		if err != nil {
			return err
		}
		err = cec.Push(destination, username, password, registry_address)
		if err != nil {
			return err
		}
	}
	return nil
}

func PythonRequirements(abspath string, filenames []string, name string, version string, logger log.Logger,
	pythonCodeStyleChecker func(abspath string, stdout *bytes.Buffer, logger log.Logger) error) (bool, error) {
	meets_reqs := true
	meets_reqs, err := PythonFileRequirements(filenames, logger)
	if err != nil {
		return false, err
	}
	good_style, err := PythonCodeStyle(abspath, name, version, pythonCodeStyleChecker, logger)
	if err != nil {
		return false, err
	} else if !good_style {
		meets_reqs = false
	}
	return meets_reqs, nil
}

func PythonCodeStyle(abspath string, name string, version string, checkPythonCodeStyle ExternalProgramOnFile, logger log.Logger) (bool, error) {
	stdout := &bytes.Buffer{}
	if err := checkPythonCodeStyle(abspath, stdout, logger); err != nil {
		if len(stdout.String()) > 0 {
			logger.Infof("python code style and quality check found these issues for %s version %s", name, version)
			logger.Infof("(%w)", stdout.String())
			return false, nil
		} else {
			return false, fmt.Errorf("error in executing python code style check for %s (%w)", name, err)
		}
	}
	return true, nil
}

func flake8PythonCodeStyle(abspath string, stdout *bytes.Buffer, logger log.Logger) error {
	err := os.Chdir(abspath)
	if err != nil {
		logger.Errorf("error changing current working directory to %s (%w)", abspath, err)
		panic(err)
	}
	return util.RunExternalProgram(
		"python3",
		[]string{
			"-m",
			"flake8",
			"--show-source",
			abspath,
		},
		nil,
		nil,
		stdout,
		stdout,
	)
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

func hasFilename(names []string, filename string) (bool, error) {
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

func BasicRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true
	output := ""

	if has_, err := hasFilename(filenames, "README.md"); err != nil {
		return false, err
	} else if !has_ {
		output = "Missing README.md\n"
		logger.Infof(output)
		meets_reqs = false
	}

	if has_, err := hasFilename(filenames, "Dockerfile"); err != nil {
		return false, err
	} else if !has_ {
		output = "Missing Dockerfile\n"
		logger.Infof(output)
		meets_reqs = false
	}

	if has_, err := hasFilename(filenames, "(?i).*test.*"); err != nil {
		return false, err
	} else if !has_ {
		// match case-insensitive 'test'?
		output = "Missing a test file\n"
		logger.Infof(output)
		meets_reqs = false
	}

	return meets_reqs, nil
}

func ContainerRequirements(abspath string, name string, version string, logger log.Logger) (bool, error) {
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
	has_, err := hasFilename(filenames, "Dockerfile")
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
			if has_, err := dockerfileHasLine(dockerfile, regexp_); err != nil {
				return false, err
			} else if !has_ {
				logger.Infof(loggerResp)
				meets_reqs = has_
			}
		}
	}
	return meets_reqs, nil
}

func PythonFileRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true
	has_reqs_txt, err := hasFilename(filenames, "requirements.txt")
	if err != nil {
		return false, err
	}
	has_pyproject, err := hasFilename(filenames, "pyproject.toml")
	if err != nil {
		return false, err
	}
	if !has_reqs_txt && !has_pyproject {
		logger.Infof("Missing a dependency manager: either add requirements.txt or pyproject.toml")
		meets_reqs = false
	}
	return meets_reqs, nil
}

func GolangRequirements(filenames []string, logger log.Logger) (bool, error) {
	meets_reqs := true
	if has_, err := hasFilename(filenames, "go.mod"); err != nil {
		return false, err
	} else if !has_ {
		logger.Infof("Missing go.mod")
		meets_reqs = false
	}
	if has_, err := hasFilename(filenames, "go.sum"); err != nil {
		return false, err
	} else if !has_ {
		logger.Infof("Missing go.sum")
		meets_reqs = false
	}
	return meets_reqs, nil
}

func AllTrue(checks []bool) bool {
	for _, v := range checks {
		if !v {
			return false
		}
	}
	return true
}

func ImageLanguage(filenames []string) (string, error) {
	ext2Lang := map[string]string{"go": "go", "py": "python"}
	for _, name := range filenames {
		matched, err := regexp.MatchString("(?i).*plugin.*", name)
		if err != nil {
			return "", err
		}
		if matched {
			ext := filepath.Ext(name)
			lang, ok := ext2Lang[ext[1:]]
			if ok {
				return lang, nil
			}
		}
	}
	// this seems like a bad way to finish this function
	return "", nil
}

func dockerfileHasLine(dockerfile string, line string) (bool, error) {
	matched, err := regexp.MatchString(line, dockerfile)
	if err != nil {
		return false, err
	}
	return matched, nil
}
