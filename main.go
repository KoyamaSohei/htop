package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"

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

type cpuStat struct {
	usertime   uint64
	nicetime   uint64
	systemtime uint64
	idletime   uint64
	ioWait     uint64
	irq        uint64
	softIrq    uint64
	steal      uint64
	guest      uint64
	guestnice  uint64
}

func getCpuStat(stat *cpuStat) error {
	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("/proc/stat is empty")
	}
	_, err = fmt.Sscanf(lines[0],
		"cpu  %d %d %d %d %d %d %d %d %d %d",
		&stat.usertime,
		&stat.nicetime,
		&stat.systemtime,
		&stat.idletime,
		&stat.ioWait,
		&stat.irq,
		&stat.softIrq,
		&stat.steal,
		&stat.guest,
		&stat.guestnice)
	return err
}

func getCommand(pid int) (string, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getUser(pid int) (string, error) {
	info, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("failed to get uid or gid from pid")
	}
	uid := stat.Uid
	u, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		return "", err
	}
	return u.Name, nil
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
	if runtime.GOOS != "linux" {
		log.Fatalf("htop only support linux!\n")
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initiaplistze termui: %v", err)
	}
	defer ui.Close()

	plist := widgets.NewList()
	plist.Title = "PID"
	plist.Border = false
	plist.TextStyle = ui.NewStyle(ui.ColorYellow)
	plist.WrapText = true

	ulist := widgets.NewList()
	ulist.Title = "USER"
	ulist.Border = false
	ulist.TextStyle = ui.NewStyle(ui.ColorYellow)
	ulist.WrapText = true

	clist := widgets.NewList()
	clist.Title = "Command"
	clist.Border = false
	clist.TextStyle = ui.NewStyle(ui.ColorYellow)
	clist.WrapText = true

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewCol(1.0/3, plist),
		ui.NewCol(1.0/3, ulist),
		ui.NewCol(1.0/3, clist),
	)

	pids, err := getPids()
	if err != nil {
		log.Fatalf("failed to get pid: %v", err)
	}
	for _, p := range pids {
		plist.Rows = append(plist.Rows, fmt.Sprintf("%d", p))
		u, err := getUser(p)
		if err != nil {
			log.Fatalf("failed to get user: %v", err)
		}
		ulist.Rows = append(ulist.Rows, u)
		c, err := getCommand(p)
		if err != nil {
			log.Fatalf("failed to get command: %v", err)
		}
		clist.Rows = append(clist.Rows, trim(c))
	}

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			plist.ScrollDown()
			ulist.ScrollDown()
			clist.ScrollDown()
		case "k", "<Up>":
			plist.ScrollUp()
			ulist.ScrollUp()
			clist.ScrollUp()
		case "<C-d>":
			plist.ScrollHalfPageDown()
			ulist.ScrollHalfPageDown()
			clist.ScrollHalfPageDown()
		case "<C-u>":
			plist.ScrollHalfPageUp()
			ulist.ScrollHalfPageUp()
			clist.ScrollHalfPageUp()
		case "<C-f>":
			plist.ScrollPageDown()
			ulist.ScrollPageDown()
			clist.ScrollPageDown()
		case "<C-b>":
			plist.ScrollPageUp()
			ulist.ScrollPageUp()
			clist.ScrollPageUp()
		}
		ui.Render(grid)
	}

}
