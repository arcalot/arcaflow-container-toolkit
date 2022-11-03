package dto

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.arcalot.io/log"
	log2 "log"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type Registry struct {
	Url              string
	Username_Envvar  string
	Password_Envvar  string
	Namespace_Envvar string
	Username         string `default:""`
	Password         string `default:""`
	Namespace        string `default:""`
}

type Registries []Registry

type Empty struct{}

func (s *Registry) SetDefaults() {
	if len(s.Namespace) == 0 {
		s.Namespace = s.Username
	}
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
		username := LookupEnvVar(username_envvar, logger).return_value
		password := LookupEnvVar(password_envvar, logger).return_value
		namespace := LookupEnvVar(namespace_envvar, logger).return_value
		if len(username) > 0 && len(password) > 0 {
			registries[i].Username = username
			registries[i].Password = password
			if len(namespace) == 0 {
				if robot, err := UserIsQuayRobot(username); err != nil {
					return nil, err
				} else if robot {
					robot_owner := strings.Split(username, "+")
					registries[i].Namespace = robot_owner[0]
				} else {
					registries[i].Namespace = registries[i].Username
				}
			} else {
				registries[i].Namespace = namespace
			}
		} else {
			logger.Infof("Missing credentials for %s\n", registries[i].Url)
			misconfigured_registries[strconv.FormatInt(int64(i), 10)] = PlaceHolder
		}
	}
	filteredRegistries := FilterByIndex(registries, misconfigured_registries)
	return filteredRegistries, nil
}

func TestUserIsQuayRobot(t *testing.T) {
	testCases := map[string]struct {
		username       string
		expectedResult bool
	}{
		"a": {
			"river+robot",
			true,
		},
		"b": {
			"river+",
			false,
		},
		"c": {
			"river",
			false,
		},
		"d": {
			"+robot",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			act, err := UserIsQuayRobot(tc.username)
			if err != nil {
				log2.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
		})
	}
}
