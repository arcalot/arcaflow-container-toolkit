/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/arcalot/arcaflow-plugin-image-builder/internal/ce_client"
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/util"
	"github.com/creasty/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Push bool
var Build bool

type Empty struct{}

type config struct {
	Revision         string `yaml:"revision"`
	Image_Name       string `default:"all"`
	Project_Filepath string
	Image_Tag        string `default:"latest"`
	Registries       []Registry
}

type Registry struct {
	Url             string
	Username_Envvar string
	Password_Envvar string
	Username        string `default:""`
	Password        string `default:""`
}

type verbose struct {
	msg          string
	return_value string
}

type ExternalProgramOnFile func(executable_filepath string, stdout *bytes.Buffer) error

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
			log.Fatal(err)
		}
		conf, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}
		abspath, err := filepath.Abs(conf.Project_Filepath)
		if err != nil {
			log.Fatal(err)
		}
		files, err := os.Open(abspath)
		if err != nil {
			log.Fatal(err)
		}
		defer files.Close()
		filenames, err := files.Readdirnames(0)
		if err != nil {
			log.Fatal(err)
		}

		if err := BuildCmdMain(Build, Push, cec, conf, abspath, filenames); err != nil {
			log.Fatal(err)
		}
	},
}

func BuildCmdMain(build_img bool, push_img bool, cec ce_client.ContainerEngineClient, conf config, abspath string, filenames []string) error {
	for _, registry := range conf.Registries {
		meets_reqs := make([]bool, 3)
		basic_reqs, err := BasicRequirements(filenames)
		if err != nil {
			return err
		}
		meets_reqs[0] = basic_reqs
		container_reqs, err := ContainerRequirements(abspath, conf.Image_Name, conf.Image_Tag)
		if err != nil {
			return err
		}
		meets_reqs[1] = container_reqs
		lang_req, err := LanguageRequirements(abspath, filenames, conf.Image_Name, conf.Image_Tag)
		if err != nil {
			return err
		}
		meets_reqs[2] = lang_req
		all_checks := AllTrue(meets_reqs)
		if err := BuildImage(build_img, all_checks, cec, abspath, conf.Image_Name, conf.Image_Tag); err != nil {
			return err
		}
		if err := PushImage(all_checks, build_img, push_img, cec, conf.Image_Name, conf.Image_Tag, registry.Username, registry.Password, registry.Url); err != nil {
			return err
		}
		if all_checks && !build_img {
			fmt.Printf("Passed all requirements: %s %s\n", conf.Image_Name, conf.Image_Tag)
		} else {
			fmt.Printf("Failed requirements check, not building: %s %s\n", conf.Image_Name, conf.Image_Tag)
		}
	}
	return nil
}

func BuildImage(build_img bool, all_checks bool, cec ce_client.ContainerEngineClient, abspath string, image_name string, image_tag string) error {
	if all_checks && build_img {
		fmt.Printf("Building %s %s from %v\n", image_name, image_tag, abspath)
		if err := cec.Build(abspath, image_name, []string{image_tag}); err != nil {
			return err
		}
	}
	return nil
}

func PushImage(all_checks, build_image, push_image bool, cec ce_client.ContainerEngineClient, name, version, username, password, registry_address string) error {
	if all_checks && build_image && push_image {
		fmt.Printf("Pushing %s version %s to registry %s\n", name, version, registry_address)
		image_name_tag := name + ":" + version

		destination := filepath.Join(registry_address, username, name)
		if robot, err := UserIsQuayRobot(username); err != nil {
			return err
		} else if robot {
			robot_owner := strings.Split(username, "+")
			destination = filepath.Join(registry_address, robot_owner[0], name)
		}
		destination = destination + ":" + version

		err2 := cec.Tag(image_name_tag, destination)
		if err2 != nil {
			return err2
		}

		err3 := cec.Push(destination, username, password, registry_address)
		if err3 != nil {
			return err3
		}
	}
	return nil
}

func getConfig() (config, error) {
	var Registries []Registry
	var PlaceHolder struct{}

	viper.UnmarshalKey("registries", &Registries)
	misconfigured_registries := make(map[string]Empty)
	for i := range Registries {
		username_envvar := Registries[i].Username_Envvar
		password_envvar := Registries[i].Password_Envvar
		username := LookupEnvVar(username_envvar).return_value
		password := LookupEnvVar(password_envvar).return_value
		if len(username) > 0 && len(password) > 0 {
			Registries[i].Username = username
			Registries[i].Password = password
		} else {
			fmt.Printf("Missing credentials for %s\n", Registries[i].Url)
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
		}
	}
	filteredRegistries := FilterByIndex(Registries, misconfigured_registries)
	conf := config{
		Revision:         viper.GetString("revision"),
		Image_Name:       viper.GetString("image_name"),
		Project_Filepath: viper.GetString("project_filepath"),
		Image_Tag:        viper.GetString("image_tag"),
		Registries:       filteredRegistries}
	if err := defaults.Set(&conf); err != nil {
		log.Fatal(err)
	}
	return conf, nil
}

func PythonRequirements(abspath string, filenames []string, name string, version string) (bool, error) {
	meets_reqs := true
	meets_reqs, err := PythonFileRequirements(filenames)
	if err != nil {
		return false, err
	}

	// TODO: formatted to PEP 8?
	good_style, err := PythonCodeStyle(abspath, name, version, flake8PythonCodeStyle)
	if err != nil {
		return false, err
	} else if !good_style {
		meets_reqs = false
	}

	return meets_reqs, nil
}

func PythonCodeStyle(abspath string, name string, version string, checkPythonCodeStyle ExternalProgramOnFile) (bool, error) {
	meets_reqs := true
	stdout := &bytes.Buffer{}
	os.Chdir(abspath)

	if err := checkPythonCodeStyle(abspath, stdout); err != nil {
		err := fmt.Errorf(
			"Code style and quality check caused an error for %s version %s (%w)",
			name,
			version,
			err,
		)
		util.WriteOutput(name, version, stdout, err)
		return false, err
	}
	// fail if code style checks returns anything besides whitespace to stdout
	if len(stdout.String()) > 0 {
		meets_reqs = false
	}
	return meets_reqs, nil
}

func flake8PythonCodeStyle(abspath string, stdout *bytes.Buffer) error {
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

func LanguageRequirements(abspath string, filenames []string, name string, version string) (bool, error) {
	meets_reqs := true
	lang, err := ImageLanguage(filenames)
	if err != nil {
		return false, err
	}
	switch lang {
	case "go":
		meets_reqs, err = GolangRequirements(filenames)
		if err != nil {
			return false, err
		}
	case "python":
		meets_reqs, err = PythonRequirements(abspath, filenames, name, version)
		if err != nil {
			return false, err
		}
	default:
		fmt.Printf("Programming Language %s not supported\n", lang)
		meets_reqs = false
	}

	return meets_reqs, nil
}

func hasFilename(names []string, filename string) (bool, error) {
	for _, name := range names {
		matched, err := regexp.MatchString(filename, name)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func hasMatchedFilename(names []string, match_name string) (bool, error) {
	for _, name := range names {
		matched, err := regexp.MatchString(match_name, name)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func BasicRequirements(filenames []string) (bool, error) {
	meets_reqs := true
	output := ""

	if has_, err := hasFilename(filenames, "README.md"); err != nil {
		return false, err
	} else if !has_ {
		output = "Missing README.md\n"
		fmt.Println(output)
		meets_reqs = false
	}

	if has_, err := hasFilename(filenames, "Dockerfile"); err != nil {
		return false, err
	} else if !has_ {
		output = "Missing Dockerfile\n"
		fmt.Println(output)
		meets_reqs = false
	}

	if has_, err := hasMatchedFilename(filenames, "(?i).*test.*"); err != nil {
		return false, err
	} else if !has_ {
		// match case-insensitive 'test'?
		fmt.Print("Missing a test file\n")
		meets_reqs = false
	}

	return meets_reqs, nil
}

func ContainerRequirements(abspath string, name string, version string) (bool, error) {
	meets_reqs := true
	output := ""
	stdout := &bytes.Buffer{}
	project_files, err := os.Open(abspath)
	if err != nil {
		log.Fatal(err)
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
		fmt.Println("Missing Dockerfile")
		meets_reqs = false

	} else {
		file, err := os.ReadFile(filepath.Join(abspath, "Dockerfile"))
		if err != nil {
			return false, err
		}
		dockerfile := string(file)

		if has_, err := dockerfileHasLine(dockerfile, "FROM quay\\.io/centos/centos:stream8"); err != nil {
			return false, err
		} else if !has_ {
			output = "Dockerfile doesn't use 'FROM quay.io/centos/centos:stream8'\n"
			if _, err := os.Stdout.Write([]byte(output)); err != nil {
				panic(err)
			}
			util.WriteOutput(name, "latest", stdout, nil)
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "(ADD|COPY) .*/LICENSE /.*"); err != nil {
			return false, err
		} else if !has_ {
			// this regex could match on an invalid filepath
			output = "Dockerfile does not contain copy of arcaflow plugin license\n"
			if _, err := os.Stdout.Write([]byte(output)); err != nil {
				panic(err)
			}
			util.WriteOutput(name, "latest", stdout, nil)
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "ENTRYPOINT \\[.*\".*plugin.*\".*\\]"); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile enterypoint does not point to an executable that includes 'plugin' in its name")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "CMD \\[\\]"); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile does not contain an empty command (i.e. CMD [])")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.source=\".*\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.source")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.licenses=\"Apache-2\\.0.*\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.licenses")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.vendor=\"Arcalot project\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.vendor")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.authors=\"Arcalot contributors\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.authors")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.title=\".*\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.title")
			meets_reqs = false
		}

		if has_, err := dockerfileHasLine(dockerfile, "LABEL io.github.arcalot.arcaflow.plugin.version=\"(\\d*)(\\.?\\d*?)(\\.?\\d*?)\""); err != nil {
			return false, err
		} else if !has_ {
			fmt.Println("Dockerfile is missing LABEL io.github.arcalot.arcaflow.plugin.version that uses form X, X.Y, X.Y.Z")
			meets_reqs = false
		}
	}
	return meets_reqs, nil
}

func PythonFileRequirements(filenames []string) (bool, error) {
	meets_reqs := true
	has_reqs_txt, err := hasFilename(filenames, "requirements.txt")
	if err != nil {
		return false, err
	}
	has_pyproject, err := hasFilename(filenames, "pyproject.toml")
	if err != nil {
		return false, err
	}
	if !has_reqs_txt || !has_pyproject {
		if !has_reqs_txt {
			fmt.Println("Missing requirements.txt")
		}
		if !has_pyproject {
			fmt.Println("Missing pyproject.toml")
		}
		meets_reqs = false
	}
	return meets_reqs, nil
}

func GolangRequirements(filenames []string) (bool, error) {
	meets_reqs := true
	if has_, err := hasFilename(filenames, "go.mod"); err != nil {
		return false, err
	} else if !has_ {
		fmt.Println("Missing go.mod")
		meets_reqs = false
	}
	if has_, err := hasFilename(filenames, "go.sum"); err != nil {
		return false, err
	} else if !has_ {
		fmt.Println("Missing go.sum")
		meets_reqs = false
	}
	// TODO: formatted to gofmt?
	return meets_reqs, nil
}

func UserIsQuayRobot(username string) (bool, error) {
	matched, err := regexp.MatchString("^[a-z][a-z0-9_]{1,254}\\+[a-z][a-z0-9_]{1,254}$", username)
	if err != nil {
		return false, err
	}
	if matched {
		return true, nil
	}
	return false, nil
}

func LookupEnvVar(key string) verbose {
	val, ok := os.LookupEnv(key)
	verbose := verbose{return_value: val}
	if !ok {
		verbose.msg = fmt.Sprintf("%s not set", key)
	} else if len(val) == 0 {
		verbose.msg = fmt.Sprintf("%s is empty", key)
	}
	fmt.Println(verbose.msg)
	return verbose
}

func FilterByIndex(list []Registry, remove map[string]Empty) []Registry {
	list2 := make([]Registry, 0, 5)
	for i := range list {
		_, ok := remove[strconv.FormatInt(int64(i), 10)]
		if !ok {
			list2 = append(list2, list[i])
		}
	}
	return list2
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
