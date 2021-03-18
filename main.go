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
		log.Fatalf("failed to initiaplistze termui: %v", err)
	}
	defer ui.Close()

	plist := widgets.NewList()
	plist.Title = "PID"
	plist.Border = false
	plist.TextStyle = ui.NewStyle(ui.ColorYellow)
	plist.WrapText = true
	plist.SetRect(0, 0, 25, 40)

	pids, err := getPids()
	if err != nil {
		log.Fatalf("failed to get pid: %v", err)
	}
	for _, p := range pids {
		plist.Rows = append(plist.Rows, fmt.Sprintf("%d", p))
	}

	ui.Render(plist)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			plist.ScrollDown()
		case "k", "<Up>":
			plist.ScrollUp()
		case "<C-d>":
			plist.ScrollHalfPageDown()
		case "<C-u>":
			plist.ScrollHalfPageUp()
		case "<C-f>":
			plist.ScrollPageDown()
		case "<C-b>":
			plist.ScrollPageUp()
		}
		ui.Render(plist)
	}

}
