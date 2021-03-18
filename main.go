package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	li := widgets.NewList()
	li.Title = "PID"
	li.TextStyle = ui.NewStyle(ui.ColorYellow)
	li.WrapText = true
	li.SetRect(0, 0, 25, 40)

	pids, err := getPids()
	if err != nil {
		log.Fatalf("failed to get pid: %v", err)
	}
	for _, p := range pids {
		li.Rows = append(li.Rows, fmt.Sprintf("%d", p))
	}

	ui.Render(li)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			li.ScrollDown()
		case "k", "<Up>":
			li.ScrollUp()
		case "<C-d>":
			li.ScrollHalfPageDown()
		case "<C-u>":
			li.ScrollHalfPageUp()
		case "<C-f>":
			li.ScrollPageDown()
		case "<C-b>":
			li.ScrollPageUp()
		}
		ui.Render(li)
	}

}
