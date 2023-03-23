package dto_test

import (
	"log"
	"os"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/dto"
	"go.arcalot.io/assert"
	arcalog "go.arcalot.io/log"
)

func TestFilterByIndex(t *testing.T) {
	a := dto.Registry{Url: "a"}
	b := dto.Registry{Url: "b"}
	c := dto.Registry{Url: "c"}
	d := dto.Registry{Url: "d"}
	e := dto.Registry{Url: "e"}
	var PlaceHolder struct{}
	list := dto.Registries{a, b, c, d, e}
	remove := map[string]dto.Empty{
		"1": PlaceHolder,
		"3": PlaceHolder,
	}
	actualList := dto.FilterByIndex(list, remove)
	assert.Equals(t, actualList[0], a)
	assert.Equals(t, actualList[1], c)
	assert.Equals(t, actualList[2], e)
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
			act, err := dto.UserIsQuayRobot(tc.username)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equals(t, tc.expectedResult, act)
		})
	}
}

func TestRegistries_Parse(t *testing.T) {

	envvars := map[string]string{
		"reg1_username":  "reg1_username",
		"reg1_password":  "reg1_password",
		"reg1_namespace": "",
		"reg2_username":  "reg2_username+robot",
		"reg2_password":  "reg2_password",
		"reg2_namespace": "",
	}

	reg2_namespace := envvars["reg2_namespace"]

	a := dto.Registry{
		Url:              "a",
		Username_Envvar:  "reg1_username",
		Password_Envvar:  "reg1_password",
		Namespace_Envvar: "reg1_namespace",
	}
	b := dto.Registry{
		Url:              "b",
		Username_Envvar:  "reg2_username",
		Password_Envvar:  "reg2_password",
		Namespace_Envvar: reg2_namespace,
	}

	for envvar_key, envvar_val := range envvars {
		err := os.Setenv(envvar_key, envvar_val)
		if err != nil {
			log.Fatal(err)
		}
	}

	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	rs := dto.Registries{a, b}
	rs2, err := rs.Parse(logger)
	if err != nil {
		panic(err)
	}

	v := envvars["reg1_username"]
	assert.Equals(t, v, rs2[0].Username)
	v = envvars["reg1_password"]
	assert.Equals(t, v, rs2[0].Password)
	assert.Equals(t, "reg1_username", rs2[0].Namespace)

	v = envvars["reg2_username"]
	assert.Equals(t, v, rs2[1].Username)
	v = envvars["reg2_password"]
	assert.Equals(t, v, rs2[1].Password)
	assert.Equals(t, "reg2_username", rs2[1].Namespace)

	for envvar_key := range envvars {
		err := os.Unsetenv(envvar_key)
		if err != nil {
			log.Fatal(err)
		}
	}
}
