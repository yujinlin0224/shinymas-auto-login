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
	relativePathOfMicrosoftEdge = "Microsoft/Edge/Application/msedge.exe"
	relativePathOfGoogleChrome  = "Google/Chrome/Application/chrome.exe"
	relativePathOfChromium      = "Chromium/Application/chrome.exe"
	relativePathOfBrave         = "BraveSoftware/Brave-Browser/Application/brave.exe"

	gameEntryURL              = "https://shinycolors.enza.fun/home"
	profileDirectoryOfBrowser = "Default"

	retryInterval = 100 * time.Millisecond
	retryTimes    = 50

	waitingTimeForSigningGame = 10 * time.Second
)

var (
	pathOfProgramDirectories = []string{
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("LocalAppData"),
	}
	relativePathOfBrowsers = []string{
		relativePathOfMicrosoftEdge,
		relativePathOfGoogleChrome,
		relativePathOfChromium,
		relativePathOfBrave,
	}

	gameTitleRegexp = regexp.MustCompile(`^\s*アイドルマスター\s+シャイニーカラーズ\s*$`)

	browserPath         string
	commandToLaunchGame *exec.Cmd
)

func setBrowserPath() error {
	var err error

	for _, pathOfProgramDirectory := range pathOfProgramDirectories {
		if pathOfProgramDirectory == "" {
			continue
		}
		for _, relativePathOfBrowser := range relativePathOfBrowsers {
			path := filepath.Join(pathOfProgramDirectory, relativePathOfBrowser)
			if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
				continue
			} else if err != nil {
				return err
			} else {
				browserPath = path
				return nil
			}
		}
	}
	if browserPath == "" {
		return errors.New("cannot find browser")
	}
	return err
}

func getGameHWNDs() ([]win32.HWND, error) {
	var (
		gameHWNDs = make([]win32.HWND, 0)
		err       error
	)

	win32.EnumWindows(func(hwnd win32.HWND) bool {
		if title := hwnd.GetWindowText(); gameTitleRegexp.MatchString(title) {
			_, processID := hwnd.GetWindowThreadProcessId()
			process, err := win32.OpenProcess(
				win32Const.PROCESS_QUERY_INFORMATION|win32Const.PROCESS_VM_READ, false, processID,
			)
			if err != nil {
				return true
			}
			processBaseName, err := process.GetModuleBaseName(0)
			if err != nil {
				return true
			}
			if processBaseName != filepath.Base(browserPath) {
				return true
			}
			gameHWNDs = append(gameHWNDs, hwnd)
		}
		return true
	})
	return gameHWNDs, err
}

func checkWindowVisible(hwnd win32.HWND) bool {
	windowStyle := win32Const.WS(hwnd.GetWindowLongPtr(win32Const.GWLP_STYLE))
	isWindowVisible := (windowStyle & win32Const.WS_VISIBLE) == win32Const.WS_VISIBLE
	return isWindowVisible
}

func init() {
	if err := setBrowserPath(); err != nil {
		panic(err)
	}
	commandToLaunchGame = exec.Command(
		browserPath,
		"--profile-directory="+profileDirectoryOfBrowser,
		"--app="+gameEntryURL,
	)
}

func main() {
	var (
		gameHWNDs []win32.HWND
		err       error
	)

	// Make sure there is no game window opened
	for i := 0; i < retryTimes; i++ {
		gameHWNDs, err = getGameHWNDs()
		if err != nil {
			panic(err)
		}
		if len(gameHWNDs) == 0 {
			break
		} else {
			countOfOpenedWindows := 0
			for _, hwnd := range gameHWNDs {
				if !hwnd.IsWindow() {
					continue
				}
				isWindowVisible := checkWindowVisible(hwnd)
				if isWindowVisible {
					countOfOpenedWindows += 1
				} else {
					hwnd.SendMessage(win32Const.WM_CLOSE, 0, 0)
				}
			}
			if countOfOpenedWindows > 0 {
				panic("game is already running")
			}
		}
		if i < retryTimes-1 {
			time.Sleep(retryInterval)
		}
	}

	// Launch game and find its window
	err = commandToLaunchGame.Start()
	if err != nil {
		panic(err)
	}
	for i := 0; i < retryTimes; i++ {
		gameHWNDs, err = getGameHWNDs()
		if err != nil {
			panic(err)
		}
		if len(gameHWNDs) > 0 {
			break
		}
		if i < retryTimes-1 {
			time.Sleep(retryInterval)
		}
	}
	if len(gameHWNDs) == 0 {
		panic("cannot find any game windows")
	}

	// Hide game window and wait for signing
	for _, hwnd := range gameHWNDs {
		if !hwnd.IsWindow() {
			continue
		}
		hwnd.ShowWindow(win32Const.SW_HIDE)
	}
	time.Sleep(waitingTimeForSigningGame)

	// Close game window if it is still hidden
	for _, hwnd := range gameHWNDs {
		if !hwnd.IsWindow() {
			continue
		}
		isWindowVisible := checkWindowVisible(hwnd)
		if !isWindowVisible {
			hwnd.SendMessage(win32Const.WM_CLOSE, 0, 0)
		}
	}
}
