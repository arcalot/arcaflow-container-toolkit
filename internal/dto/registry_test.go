package dto

import (
	"github.com/stretchr/testify/assert"
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
