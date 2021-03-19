package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getPids(t *testing.T) {
	pids, err := getPids()
	assert.Nil(t, err)
	assert.Greater(t, len(pids), 0)
}

func Test_getUser(t *testing.T) {
	pids, err := getPids()
	assert.Nil(t, err)
	for _, p := range pids {
		_, err := getUser(p)
		assert.Nil(t, err)
	}
}

func Test_getProcStat(t *testing.T) {
	pre := map[int]pstat{}
	for k := 0; k < 10; k++ {
		pids, err := getPids()
		assert.Nil(t, err)
		for _, p := range pids {
			u, err := getProcStat(p)
			assert.Nil(t, err)
			assert.GreaterOrEqual(t, u.stime, pre[p].stime)
			assert.GreaterOrEqual(t, u.utime, pre[p].utime)
			assert.GreaterOrEqual(t, u.started, pre[p].started)
			pre[p] = *u
		}
		time.Sleep(time.Second)
	}

}

func Test_getCpuStat(t *testing.T) {
	pre := uint64(0)
	for k := 0; k < 10; k++ {
		var stat cpuStat
		err := getCpuStat(&stat)
		assert.Nil(t, err)
		ti := stat.getTotalTime()
		assert.GreaterOrEqual(t, ti, pre)
		pre = ti
		time.Sleep(time.Second)
	}

}
