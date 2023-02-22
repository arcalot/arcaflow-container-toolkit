package dto

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"go.arcalot.io/log"
)

type Registry struct {
	Url              string
	Username_Envvar  string
	Password_Envvar  string
	Namespace_Envvar string
	Username         string `default:""`
	Password         string `default:""`
	Namespace        string `default:""`
	Quay_Custom_Repo string `default:""`
}

type Registries []Registry

type Empty struct{}

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

func UnmarshalRegistries(logger log.Logger) ([]Registry, error) {
	var registries Registries
	err := viper.UnmarshalKey("registries", &registries)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling registries from Carpenter file (%w)", err)
	}
	return registries.Parse(logger)
}

func (registries Registries) Parse(logger log.Logger) (Registries, error) {
	var PlaceHolder struct{}
	misconfigured_registries := make(map[string]Empty)
	for i := range registries {
		username_envvar := registries[i].Username_Envvar
		password_envvar := registries[i].Password_Envvar
		namespace_envvar := registries[i].Namespace_Envvar
		quay_custom_repo_envvar := registries[i].Quay_Custom_Repo
		username := LookupEnvVar(username_envvar, logger).Return_value
		password := LookupEnvVar(password_envvar, logger).Return_value
		namespace := LookupEnvVar(namespace_envvar, logger).Return_value
		quay_custom_repo := LookupEnvVar(quay_custom_repo_envvar, logger).Return_value
		if quay_custom_repo != "" && registries[i].Url == "quay.io" {
			registries[i].Namespace = quay_custom_repo
		}
		if registries[i].ValidCredentials(username) {
			registries[i].Username = username
			registries[i].Password = password
			inferred_namespace, err := InferNamespace(namespace, username)
			if err != nil {
				return nil, err
			}
			registries[i].Namespace = inferred_namespace
		} else {
			logger.Infof("Missing credentials for %s\n", registries[i].Url)
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
		}
	}
	filteredRegistries := FilterByIndex(registries, misconfigured_registries)
	return filteredRegistries, nil
}

func (registry Registry) ValidCredentials(username string) bool {
	return len(username) > 0
}

func InferNamespace(namespace string, username string) (string, error) {
	if len(namespace) == 0 {
		robot, err := UserIsQuayRobot(username)
		if err != nil {
			return "", err
		}
		if robot {
			robot_owner := strings.Split(username, "+")
			return robot_owner[0], nil
		} else {
			return username, nil
		}
	}
	return namespace, nil
}
