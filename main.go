package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const MAX_LENGTH = 30

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

func getCommand(pid int) (string, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func trim(str string) string {
	if len(str) > MAX_LENGTH {
		str = str[:MAX_LENGTH]
	}
	s := strings.TrimRightFunc(str, func(c rune) bool {
		//In windows newline is \r\n
		return c == '\r' || c == '\n'
	})
	if len(s) > MAX_LENGTH {
		s = s[:MAX_LENGTH] + "..."
	}
	return s
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

	clist := widgets.NewList()
	clist.Title = "Command"
	clist.Border = false
	clist.TextStyle = ui.NewStyle(ui.ColorYellow)
	clist.WrapText = true
	clist.SetRect(25, 0, 60, 40)

	pids, err := getPids()
	if err != nil {
		log.Fatalf("failed to get pid: %v", err)
	}
	for _, p := range pids {
		plist.Rows = append(plist.Rows, fmt.Sprintf("%d", p))
		c, err := getCommand(p)
		if err != nil {
			log.Fatalf("failed to get command: %v", err)
		}
		clist.Rows = append(clist.Rows, trim(c))
	}

	ui.Render(plist)
	ui.Render(clist)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			plist.ScrollDown()
			clist.ScrollDown()
		case "k", "<Up>":
			plist.ScrollUp()
			clist.ScrollUp()
		case "<C-d>":
			plist.ScrollHalfPageDown()
			clist.ScrollHalfPageDown()
		case "<C-u>":
			plist.ScrollHalfPageUp()
			clist.ScrollHalfPageUp()
		case "<C-f>":
			plist.ScrollPageDown()
			clist.ScrollPageDown()
		case "<C-b>":
			plist.ScrollPageUp()
			clist.ScrollPageUp()
		}
		ui.Render(plist)
		ui.Render(clist)
	}

}
