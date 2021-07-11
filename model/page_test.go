package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPage_String(t *testing.T) {

	scenario := []struct {
		desc   string
		input  Page
		output string
	}{
		{"Just key", NewPage("Go"), "<h1>Go</h1>"},
	}

	for _, v := range scenario {
		assert.Equal(t, v.output, v.input.String(), v.desc)
	}

}
