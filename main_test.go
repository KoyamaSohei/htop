package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPids(t *testing.T) {
	pids, err := getPids()
	assert.Nil(t, err)
	assert.Greater(t, len(pids), 0)
}
