package test

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertAppErr(t *testing.T, err error, code string) {
	assert.Error(t, err)
	appEr, ok := er.Is(err)
	assert.True(t, ok)
	assert.Equal(t, code, appEr.Code())
}
