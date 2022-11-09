package dto

import (
	"github.com/stretchr/testify/assert"
	arcalog "go.arcalot.io/log"
	"log"
	"os"
	"testing"
)

func TestFilterByIndex(t *testing.T) {
	a := Registry{Url: "a"}
	b := Registry{Url: "b"}
	c := Registry{Url: "c"}
	d := Registry{Url: "d"}
	e := Registry{Url: "e"}
	var PlaceHolder struct{}
	list := []Registry{a, b, c, d, e}
	remove := map[string]Empty{
		"1": PlaceHolder,
		"3": PlaceHolder,
	}
	actualList := FilterByIndex(list, remove)
	assert.Equal(t, actualList[0], a)
	assert.Equal(t, actualList[1], c)
	assert.Equal(t, actualList[2], e)
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
				log.Fatal(err)
			}
			assert.Equal(t, tc.expectedResult, act)
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

	a := Registry{
		Url:              "a",
		Username_Envvar:  "reg1_username",
		Password_Envvar:  "reg1_password",
		Namespace_Envvar: "reg1_namespace",
	}
	b := Registry{
		Url:              "b",
		Username_Envvar:  "reg2_username",
		Password_Envvar:  "reg2_password",
		Namespace_Envvar: reg2_namespace,
	}

	for envvar_key, envvar_val := range envvars {
		os.Setenv(envvar_key, envvar_val)
	}

	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	rs := Registries{a, b}
	rs2, err := rs.Parse(logger)
	if err != nil {
		panic(err)
	}

	v := envvars["reg1_username"]
	assert.Equal(t, v, rs2[0].Username)
	v = envvars["reg1_password"]
	assert.Equal(t, v, rs2[0].Password)
	assert.Equal(t, "reg1_username", rs2[0].Namespace)

	v = envvars["reg2_username"]
	assert.Equal(t, v, rs2[1].Username)
	v = envvars["reg2_password"]
	assert.Equal(t, v, rs2[1].Password)
	assert.Equal(t, "reg2_username", rs2[1].Namespace)

	for envvar_key := range envvars {
		os.Unsetenv(envvar_key)
	}
}
