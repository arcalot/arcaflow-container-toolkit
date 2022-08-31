package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type config struct {
	Revision         string `yaml:"revision"`
	Target           string `default:"all"`
	Project_Filepath string
}

type Image struct {
	name    string
	context string
}

func runExternalProgram(
	program string,
	args []string,
	env []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	_, _ = stdout.Write([]byte(fmt.Sprintf("\033[0;32m⚙ Running %s...\u001B[0m\n", program)))
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
	prefix := "\033[0;32m✅ "
	if err != nil {
		prefix = "\033[0;31m❌ "
	}
	output += fmt.Sprintf(
		"::group::%s img=%s version=%s\n",
		prefix,
		image,
		version,
	)
	output += stdout.String()
	if err != nil {
		output += fmt.Sprintf("\033[0;31m%s\033[0m\n", err.Error())
	}
	output += "::endgroup::\n"
	if _, err := os.Stdout.Write([]byte(output)); err != nil {
		panic(err)
	}
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

func getConfig(configYamlPath string) *config {
	fh, err := os.Open(configYamlPath)
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(fh)
	if err != nil {
		log.Fatal(err)
	}
	conf := &config{}
	if err := yaml.Unmarshal(data, conf); err != nil {
		log.Fatal(err)
	}
	if err := defaults.Set(conf); err != nil {
		log.Fatal(err)
	}
	return conf
}

func listImagesToBuild(conf *config) []Image {
	if allDirectories(conf.Project_Filepath) {
		list := listPackagesFromFile(conf.Project_Filepath)
		return filterContainerSelection(conf.Target, list)
	}

	abspath, err := filepath.Abs(conf.Project_Filepath)
	if err != nil {
		log.Fatal(err)
	}
	return []Image{{
		name:    filepath.Base(conf.Project_Filepath),
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
	files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer files.Close()
	filenames, _ := files.Readdirnames(0)

	if !hasFilename(filenames, "README.md") {
		fmt.Println("Missing README.md")
		meets_reqs = false
	}
	if !hasFilename(filenames, "Dockerfile") {
		fmt.Println("Missing Dockerfile")
		meets_reqs = false
	}
	if !hasMatchedFilename(filenames, "(?i).*test.*") {
		// match case-insensitive 'test'?
		fmt.Println("Missing a test file")
		meets_reqs = false
	}
	return meets_reqs
}

func licenseRequirements(img Image) bool {
	meets_reqs := true
	files, err := os.Open(img.context)
	if err != nil {
		log.Fatal(err)
	}
	defer files.Close()
	filenames, _ := files.Readdirnames(0)
	if !hasFilename(filenames, "LICENSE") {
		fmt.Println("Missing LICENSE")
		meets_reqs = false
	} else {
		license, err2 := os.ReadFile(filepath.Join(img.context, "LICENSE"))
		if err2 != nil {
			log.Fatal(err2)
		}
		matched, err3 := regexp.MatchString("Apache License\\s*Version 2.0, January 2004", string(license))
		if err3 != nil {
			log.Fatal(err3)
		}
		if matched {
			return true
		}
	}
	return meets_reqs
}

func containerRequirements(img Image) bool {
	return true
}

func allTrue(checks []bool) bool {
	for _, v := range checks {
		if !v {
			return false
		}
	}
	return true
}

func main() {
	conf := getConfig("build.yaml")
	for _, img := range listImagesToBuild(conf) {
		fmt.Printf("Building %s from %v\n", img.name, img.context)
		meets_reqs := make([]bool, 4)
		meets_reqs[0] = basicRequirements(img)
		meets_reqs[1] = containerRequirements(img)
		if allTrue(meets_reqs) {
			// if err := buildVersion(img, "latest", conf.Revision); err != nil {
			//     log.Fatal(err)
			// }
		}
	}
}
