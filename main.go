package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	win32 "github.com/rodrigocfd/windigo/win"
	win32Const "github.com/rodrigocfd/windigo/win/co"
)

const (
	gameEntryURL = "https://shinycolors.enza.fun/home"

	retryInterval = 250 * time.Millisecond
	retryTimes    = 10
	waitingTime   = 10 * time.Second
)

var (
	gameTitleRegexp = regexp.MustCompile(`^アイドルマスター\s+シャイニーカラーズ$`)

	commandToLaunchGame *exec.Cmd
)

func getPathOfMicrosoftEdge() (string, error) {
	programFilesPath := os.Getenv("ProgramFiles(x86)")
	if programFilesPath == "" {
		programFilesPath = os.Getenv("ProgramFiles")
	}
	if programFilesPath == "" {
		return "", errors.New("cannot find path of \"Program Files\" directory")
	}
	return filepath.Join(programFilesPath, "Microsoft/Edge/Application/msedge.exe"), nil
}

func init() {
	var err error

	pathOfMicrosoftEdge, err := getPathOfMicrosoftEdge()
	if err != nil {
		panic(err)
	}
	commandToLaunchGame = exec.Command(
		pathOfMicrosoftEdge,
		"--profile-directory=Default",
		"--app="+gameEntryURL,
	)
}

func main() {
	var (
		gameHWNDs = make([]win32.HWND, 0)
		err       error
	)

	err = commandToLaunchGame.Start()
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
		panic("cannot find any game windows")
	}
	for _, hwnd := range gameHWNDs {
		hwnd.ShowWindow(win32Const.SW_HIDE)
	}
	time.Sleep(waitingTime)
	for _, hwnd := range gameHWNDs {
		hwnd.SendMessage(win32Const.WM_CLOSE, 0, 0)
	}
}
