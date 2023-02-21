package main

import (
	"os/exec"
	"regexp"
	"time"

	win32 "github.com/rodrigocfd/windigo/win"
	win32Const "github.com/rodrigocfd/windigo/win/co"
)

const (
	retryInterval = 250 * time.Millisecond
	retryTimes    = 10
	waitingTime   = 10 * time.Second
)

func main() {
	var (
		gameTitleRegexp = regexp.MustCompile(`^アイドルマスター\s+シャイニーカラーズ$`)
		gameHWNDs       = make([]win32.HWND, 0)
		err             error
	)

	err = exec.Command(
		"C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
		"--profile-directory=Default",
		"--app=https://shinycolors.enza.fun/home",
	).Start()
	if err != nil {
		panic(err)
	}
	for i := 0; i < retryTimes; i++ {
		time.Sleep(retryInterval)
		win32.EnumWindows(func(hwnd win32.HWND) bool {
			if title := hwnd.GetWindowText(); gameTitleRegexp.MatchString(title) {
				gameHWNDs = append(gameHWNDs, hwnd)
			}
			return true
		})
		if len(gameHWNDs) > 0 {
			break
		}
	}
	if len(gameHWNDs) == 0 {
		panic("Cannot find any game windows")
	}
	for _, hwnd := range gameHWNDs {
		hwnd.ShowWindow(win32Const.SW_HIDE)
	}
	time.Sleep(waitingTime)
	for _, hwnd := range gameHWNDs {
		hwnd.SendMessage(win32Const.WM_CLOSE, 0, 0)
	}
}
