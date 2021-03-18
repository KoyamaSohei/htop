package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func getPids() ([]int, error) {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0)
	for _, f := range files {
		if k, err := strconv.Atoi(f.Name()); err != nil {
			continue
		} else {
			pids = append(pids, k)
		}
	}
	return pids, nil
}

func main() {
	pids, err := getPids()
	if err != nil {
		panic(err)
	}
	for _, p := range pids {
		fmt.Println(p)
	}
}
