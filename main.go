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
		log.Fatalf("failed to initiapidze termui: %v", err)
	}
	defer ui.Close()

	pid := widgets.NewList()
	pid.Title = "PID"
	pid.Border = false
	pid.TextStyle = ui.NewStyle(ui.ColorYellow)
	pid.WrapText = true

	user := widgets.NewList()
	user.Title = "USER"
	user.Border = false
	user.TextStyle = ui.NewStyle(ui.ColorYellow)
	user.WrapText = true

	command := widgets.NewList()
	command.Title = "Command"
	command.Border = false
	command.TextStyle = ui.NewStyle(ui.ColorYellow)
	command.WrapText = true

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewCol(1.0/3, pid),
		ui.NewCol(1.0/3, user),
		ui.NewCol(1.0/3, command),
	)

	update := func() {
		pids, err := getPids()
		if err != nil {
			log.Fatalf("failed to get pid: %v", err)
		}
		pid.Rows = make([]string, 0)
		user.Rows = make([]string, 0)
		command.Rows = make([]string, 0)
		for _, p := range pids {
			pid.Rows = append(pid.Rows, fmt.Sprintf("%d", p))
			u, err := getUser(p)
			if err != nil {
				log.Fatalf("failed to get user: %v", err)
			}
			user.Rows = append(user.Rows, u)
			c, err := getCommand(p)
			if err != nil {
				log.Fatalf("failed to get command: %v", err)
			}
			command.Rows = append(command.Rows, trim(c))
		}
		ui.Render(grid)
	}

	update()

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			pid.ScrollDown()
			user.ScrollDown()
			command.ScrollDown()
		case "k", "<Up>":
			pid.ScrollUp()
			user.ScrollUp()
			command.ScrollUp()
		case "<C-d>":
			pid.ScrollHalfPageDown()
			user.ScrollHalfPageDown()
			command.ScrollHalfPageDown()
		case "<C-u>":
			pid.ScrollHalfPageUp()
			user.ScrollHalfPageUp()
			command.ScrollHalfPageUp()
		case "<C-f>":
			pid.ScrollPageDown()
			user.ScrollPageDown()
			command.ScrollPageDown()
		case "<C-b>":
			pid.ScrollPageUp()
			user.ScrollPageUp()
			command.ScrollPageUp()
		}
		ui.Render(grid)
	}

}
