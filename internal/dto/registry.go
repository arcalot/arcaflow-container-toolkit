package dto

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"go.arcalot.io/log"
)

type ErrorPushArg struct {
	// location that will receive payload
	target string
	// argument required to upload to target location
	arg string
}

func (e ErrorPushArg) Error() string {
	return fmt.Sprintf("Push argument %q detected for %q, but error found. Not attempting to build or push.", e.arg, e.target)
}

type Registry struct {
	Url                          string
	Username_Envvar              string
	Password_Envvar              string
	Namespace_Envvar             string
	Quay_Custom_Namespace_Envvar string
	Username                     string `default:""`
	Password                     string `default:""`
	Namespace                    string `default:""`
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
		return nil, fmt.Errorf("error unmarshalling registries from act file ((%w))", err)
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
		quay_custom_namespace_envvar := registries[i].Quay_Custom_Namespace_Envvar
		quay_custom_namespace, _ := LookupEnvVar(registries[i].Url, quay_custom_namespace_envvar, logger)
		if quay_custom_namespace != "" && registries[i].Url == "quay.io" {
			logger.Infof("value is:%s", quay_custom_namespace_envvar)
			logger.Infof("QUAY_CUSTOM_NAMESPACE environment variable detected,"+
				"using value in place of QUAY_NAMESPACE for %s", registries[i].Url)
			namespace_envvar = quay_custom_namespace_envvar
		}
		username, err := LookupEnvVar(registries[i].Url, username_envvar, logger)
		if err != nil {
			logger.Errorf("%s", errors.Join(ErrorPushArg{registries[i].Url, "username"}, err))
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
			continue
		}
		password, err := LookupEnvVar(registries[i].Url, password_envvar, logger)
		if err != nil {
			logger.Errorf("%s", errors.Join(ErrorPushArg{registries[i].Url, "password"}, err))
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
			continue
		}
		namespace, err := LookupEnvVar(registries[i].Url, namespace_envvar, logger)
		if err != nil {
			logger.Errorf("%s", errors.Join(ErrorPushArg{registries[i].Url, "namespace"}, err))
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
			continue
		}
		if !registries[i].ValidCredentials(username) {
			logger.Errorf("Missing credentials for %s\n", registries[i].Url)
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
			continue
		}
		registries[i].Username = username
		registries[i].Password = password

		inferred_namespace, err := InferNamespace(namespace, username)
		if err != nil {
			return nil, err
		}
		registries[i].Namespace = inferred_namespace

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
