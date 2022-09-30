/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	username        string `default:""`
	password        string `default:""`
}

type Image struct {
	name    string
	context string
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.PersistentFlags().BoolVarP(&Push, "push", "p", false, "push images to registry")
	buildCmd.PersistentFlags().BoolVarP(&Build, "build", "b", false, "validate requirements and build image")
}

var buildCmd = &cobra.Command{
	Use:   "build an image",
	Short: "build image",
	Run: func(cmd *cobra.Command, args []string) {

		conf := getConfig()

		for _, registry := range conf.Registries {
			for _, img := range listImagesToBuild(&conf) {
				meets_reqs := make([]bool, 3)
				meets_reqs[0] = basicRequirements(img)
				meets_reqs[1] = containerRequirements(img)
				meets_reqs[2] = languageRequirements(img, conf.Image_Tag)
				all_checks := allTrue(meets_reqs)

				if !all_checks {
					fmt.Println("Please fix requirements.")
					os.Exit(1)
				}

				if Build {
					fmt.Printf("Building %s %s from %v\n", img.name, conf.Image_Tag, img.context)
					if err := buildVersion(img, conf.Image_Tag, conf.Revision); err != nil {
						log.Fatal(err)
					}
					if Push {
						fmt.Printf("Pushing %s version %s to registry %s\n", img.name, conf.Image_Tag, registry.Url)
						if err := pushImage(img, conf.Image_Tag, registry); err != nil {
							log.Fatal(err)
						}
					}
				} else if all_checks && !Build {
					fmt.Printf("Passed all requirements: %s %s\n", img.name, conf.Image_Tag)
				} else {
					fmt.Printf("Failed requirements check, not building: %s %s\n", img.name, conf.Image_Tag)
				}
			}
		}
	},
}

func pushImage(image Image, version string, registry Registry) error {
	image_tag := image.name + ":" + version
	env := []string{
		fmt.Sprintf("BLDIMG=%s/", image_tag),
	}
	stdout := &bytes.Buffer{}

	if err := runExternalProgram(
		"docker",
		[]string{
			"login",
			"--username",
			registry.username,
			"--password",
			registry.password,
			registry.Url,
		},
		env,
		nil,
		stdout,
		stdout,
	); err != nil {
		err := fmt.Errorf(
			"Error logging in for %s version %s (%w)",
			registry.username,
			version,
			err,
		)
		writeOutput(image.name, version, stdout, err)
		return err
	}

	destination := filepath.Join(registry.Url, registry.username, image.name)
	if userIsQuayRobot(registry.username) {
		robot_owner := strings.Split(registry.username, "+")
		destination = filepath.Join(registry.Url, robot_owner[0], image.name)
	}
	destination = destination + ":" + version

	if err := runExternalProgram(
		"docker",
		[]string{
			"tag",
			image_tag,
			destination,
		},
		env,
		nil,
		stdout,
		stdout,
	); err != nil {
		err := fmt.Errorf(
			"Error tagging for %s version %s (%w)",
			image.name,
			version,
			err,
		)
		writeOutput(image.name, version, stdout, err)
		return err
	}

	if err := runExternalProgram(
		"docker",
		[]string{
			"push",
			destination,
		},
		env,
		nil,
		stdout,
		stdout,
	); err != nil {
		err := fmt.Errorf(
			"Error pushing for %s version %s (%w)",
			image.name,
			version,
			err,
		)
		writeOutput(image.name, version, stdout, err)
		return err
	}
	return nil
}

func buildVersion(
	image Image,
	version string,
	date string,
) error {

	image_tag := image.name + ":" + version
	stdout := &bytes.Buffer{}
	env := []string{
		fmt.Sprintf("BLDIMG=%s/", image_tag),
	}
	os.Chdir(image.context)

	if err := runExternalProgram(
		"docker",
		[]string{
			"build",
			".",
			"--tag",
			image_tag,
		},
		env,
		nil,
		stdout,
		stdout,
	); err != nil {
		err := fmt.Errorf(
			"build failed for %s version %s (%w)",
			image.name,
			version,
			err,
		)
		writeOutput(image.name, version, stdout, err)
		return err
	}
	writeOutput(image.name, version, stdout, nil)

	return nil
}

func allTrue(checks []bool) bool {
	for _, v := range checks {
		if !v {
			return false
		}
	}
	return true
}

func lookupEnvVar(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("%s not set\n", key)
	} else if len(val) == 0 {
		fmt.Printf("%s is empty\n", key)
	}
	return val
}

func filterByIndex(list []Registry, remove map[string]Empty) []Registry {
	list2 := make([]Registry, 0, 5)
	for i := range list {
		_, ok := remove[strconv.FormatInt(int64(i), 10)]
		if !ok {
			list2 = append(list2, list[i])
		}
	}
	return list2
}

func getConfig() config {
	var Registries []Registry
	var PlaceHolder struct{}

	viper.UnmarshalKey("registries", &Registries)
	misconfigured_registries := make(map[string]Empty)
	for i := range Registries {
		username_envvar := Registries[i].Username_Envvar
		password_envvar := Registries[i].Password_Envvar
		username := lookupEnvVar(username_envvar)
		password := lookupEnvVar(password_envvar)
		if len(username) > 0 && len(password) > 0 {
			Registries[i].username = username
			Registries[i].password = password
		} else {
			fmt.Printf("Missing credentials for %s\n", Registries[i].Url)
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
		}
	}
	filteredRegistries := filterByIndex(Registries, misconfigured_registries)
	conf := config{
		Revision:         viper.GetString("revision"),
		Image_Name:       viper.GetString("image_name"),
		Project_Filepath: viper.GetString("project_filepath"),
		Image_Tag:        viper.GetString("image_tag"),
		Registries:       filteredRegistries}
	if err := defaults.Set(&conf); err != nil {
		log.Fatal(err)
	}
	return conf
}

func runExternalProgram(
	program string,
	args []string,
	env []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	programPath, err := exec.LookPath(program)
	if err != nil {
		return err
	}
	env = append(env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	env = append(env, fmt.Sprintf("TMP=%s", os.Getenv("TMP")))
	env = append(env, fmt.Sprintf("TEMP=%s", os.Getenv("TEMP")))
	cmd := &exec.Cmd{
		Path:   programPath,
		Args:   append([]string{programPath}, args...),
		Env:    env,
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func writeOutput(
	image string,
	version string,
	stdout *bytes.Buffer,
	err error,
) {
	output := ""
	output += fmt.Sprintf(
		"img=%s version=%s%\n",
		image,
		version,
	)
	output += stdout.String()
	if err != nil {
		output += fmt.Sprintf(err.Error())
	}
	if _, err := os.Stdout.Write([]byte(output)); err != nil {
		panic(err)
	}
}

func listPackagesFromFile(source_project string) []Image {
	var pwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var source_project_dir string = filepath.Join(pwd, source_project)
	list := make([]Image, 0, 10)
	source_project_file, err2 := os.Open(source_project_dir)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer source_project_file.Close()
	lst, _ := source_project_file.Readdir(-1)
	for _, file := range lst {
		if file.IsDir() {
			list = append(list, Image{file.Name(), filepath.Join(source_project_dir, file.Name())})
		}
	}
	return list
}

func filterContainerSelection(selection string, list []Image) []Image {
	if selection != "all" {
		list2 := make([]Image, 0, 10)
		for _, container := range list {
			if container.name == selection {
				list2 = append(list2, container)
			}
		}
		list = list2
	}
	return list
}

func listImagesToBuild(conf *config) []Image {
	files, err := os.Open(conf.Project_Filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer files.Close()
	filenames, err := files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}

	if !hasFilename(filenames, "Dockerfile") {
		list := listPackagesFromFile(conf.Project_Filepath)
		return filterContainerSelection(conf.Image_Name, list)
	}

	abspath, err := filepath.Abs(conf.Project_Filepath)
	if err != nil {
		log.Fatal(err)
	}
	return []Image{{
		name:    conf.Image_Name,
		context: abspath}}
}

func allDirectories(filepath string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lst, _ := file.Readdir(-1)
	for _, file := range lst {
		if !file.IsDir() {
			return false
		}
	}
	return true
}

func hasFilename(names []string, filename string) bool {
	for _, name := range names {
		matched, err4 := regexp.MatchString(filename, name)
		if err4 != nil {
			log.Fatal(err4)
		}
		if matched {
			return true
		}
	}
	return false
}

func hasMatchedFilename(names []string, match_name string) bool {
	for _, name := range names {
		matched, err4 := regexp.MatchString(match_name, name)
		if err4 != nil {
			log.Fatal(err4)
		}
		if matched {
			return true
		}
	}
	return false
}

func basicRequirements(img Image) bool {
	meets_reqs := true
	stdout := &bytes.Buffer{}
	files, err := os.Open(img.context)
	output := ""
	if err != nil {
		log.Fatal(err)
	}
	defer files.Close()
	filenames, err := files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}

	if !hasFilename(filenames, "README.md") {
		output = "Missing README.md\n"
		if _, err := os.Stdout.Write([]byte(output)); err != nil {
			panic(err)
		}
		writeOutput(img.name, "latest", stdout, nil)
		meets_reqs = false
	}
	if !hasFilename(filenames, "Dockerfile") {
		output = "Missing Dockerfile\n"
		if _, err := os.Stdout.Write([]byte(output)); err != nil {
			panic(err)
		}
		writeOutput(img.name, "latest", stdout, nil)
		meets_reqs = false
	}
	if !hasMatchedFilename(filenames, "(?i).*test.*") {
		// match case-insensitive 'test'?
		fmt.Print("Missing a test file\n")
		meets_reqs = false
	}
	return meets_reqs
}

func dockerfileHasLine(dockerfile string, line string) bool {
	matched, err := regexp.MatchString(line, dockerfile)
	if err != nil {
		log.Fatal(err)
	}
	return matched
}

func imageLanguage(filenames []string) string {
	ext2Lang := map[string]string{"go": "go", "py": "python"}
	for _, name := range filenames {
		matched, err := regexp.MatchString("(?i).*plugin.*", name)
		if err != nil {
			log.Fatal(err)
		}
		if matched {
			ext := filepath.Ext(name)
			lang, ok := ext2Lang[ext[1:]]
			if ok {
				return lang
			}
		}
	}
	// this seems like a bad way to finish this function
	return ""
}

func containerRequirements(img Image) bool {
	meets_reqs := true
	output := ""
	stdout := &bytes.Buffer{}
	project_files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer project_files.Close()
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	if !hasFilename(filenames, "Dockerfile") {
		fmt.Println("Missing Dockerfile")
		meets_reqs = false

	} else {
		file, err := os.ReadFile(filepath.Join(img.context, "Dockerfile"))
		if err != nil {
			log.Fatal(err)
		}
		dockerfile := string(file)

		if !dockerfileHasLine(dockerfile, "FROM quay\\.io/centos/centos:stream8") {
			output = "Dockerfile doesn't use 'FROM quay.io/centos/centos:stream8'\n"
			if _, err := os.Stdout.Write([]byte(output)); err != nil {
				panic(err)
			}
			writeOutput(img.name, "latest", stdout, nil)
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "(ADD|COPY) .*/LICENSE /.*") {
			// this regex could match on an invalid filepath
			output = "Dockerfile does not contain copy of arcaflow plugin license\n"
			if _, err := os.Stdout.Write([]byte(output)); err != nil {
				panic(err)
			}
			writeOutput(img.name, "latest", stdout, nil)
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "ENTRYPOINT \\[.*\".*plugin.*\".*\\]") {
			fmt.Println("Dockerfile enterypoint does not point to an executable that includes 'plugin' in its name")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "CMD \\[\\]") {
			fmt.Println("Dockerfile does not contain an empty command (i.e. CMD [])")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.source=\".*\"") {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.source")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.licenses=\"Apache-2\\.0.*\"") {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.licenses")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.vendor=\"Arcalot project\"") {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.vendor")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.authors=\"Arcalot contributors\"") {
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.authors")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL org.opencontainers.image.title=\".*\"") {
			// this title regular expression could match anything
			fmt.Println("Dockerfile is missing LABEL org.opencontainers.image.title")
			meets_reqs = false
		}
		if !dockerfileHasLine(dockerfile, "LABEL io.github.arcalot.arcaflow.plugin.version=\"(\\d*)(\\.?\\d*?)(\\.?\\d*?)\"") {
			fmt.Println("Dockerfile is missing LABEL io.github.arcalot.arcaflow.plugin.version that uses form X, X.Y, X.Y.Z")
			meets_reqs = false
		}
	}
	return meets_reqs
}

func golangRequirements(img Image) bool {
	meets_reqs := true
	project_files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer project_files.Close()
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}

	if !hasFilename(filenames, "go.mod") {
		fmt.Println("Missing go.md")
		meets_reqs = false
	}
	if !hasFilename(filenames, "go.sum") {
		fmt.Println("Missing go.sum")
		meets_reqs = false
	}
	// TODO: formatted to gofmt?
	return meets_reqs
}

func pythonRequirements(img Image, version string) bool {
	meets_reqs := true
	project_files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer project_files.Close()
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	has_reqs_txt := hasFilename(filenames, "requirements.txt")
	has_pyproject := hasFilename(filenames, "pyproject.toml")
	if !has_reqs_txt && !has_pyproject {
		if !has_reqs_txt {
			fmt.Println("Missing requirements.txt")
		}
		if !has_pyproject {
			fmt.Println("Missing pyproject.toml")
		}
		meets_reqs = false
	}
	// TODO: formatted to PEP 8?
	if !pythonCodeStyle(img, version) {
		meets_reqs = false
	}

	return meets_reqs
}

func pythonCodeStyle(image Image, version string) bool {
	meets_reqs := true

	image_tag := image.name + ":" + version
	stdout := &bytes.Buffer{}
	env := []string{
		fmt.Sprintf("BLDIMG=%s/", image_tag),
	}
	os.Chdir(image.context)

	if err := runExternalProgram(
		"python3",
		[]string{
			"-m",
			"flake8",
			"--show-source",
			image.context,
		},
		env,
		nil,
		stdout,
		stdout,
	); err != nil {
		err := fmt.Errorf(
			"Code style and quality check caused an error for %s version %s (%w)",
			image.name,
			version,
			err,
		)
		writeOutput(image.name, version, stdout, err)
	}
	// fail if code style checks returns anything besides whitespace to stdout
	if len(stdout.String()) > 0 {
		meets_reqs = false
	}
	return meets_reqs
}

func languageRequirements(img Image, version string) bool {
	meets_reqs := true
	project_files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer project_files.Close()
	filenames, err := project_files.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}

	switch lang := imageLanguage(filenames); lang {
	case "go":
		meets_reqs = golangRequirements(img)
	case "python":
		meets_reqs = pythonRequirements(img, "latest")
	default:
		fmt.Printf("Programming Language %s not supported\n", lang)
		meets_reqs = false
	}

	return meets_reqs
}

func userIsQuayRobot(username string) bool {
	matched, err := regexp.MatchString("^[a-z][a-z0-9_]{1,254}\\+[a-z][a-z0-9_]{1,254}$", username)
	if err != nil {
		log.Fatal(err)
	}
	if matched {
		return true
	}
	return false
}
