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

func Test_getCpuStat(t *testing.T) {
	var stat cpuStat
	err := getCpuStat(&stat)
	assert.Nil(t, err)
}
