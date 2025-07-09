package main

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"github.com/tailscale/win"
	"syscall"
	"time"
)

var (
	libUser32    = syscall.NewLazyDLL("user32.dll")
	procIsWindow = libUser32.NewProc("IsWindow")
)

type AsyncResult[T any] struct {
	Value T
	Error error
}

type WinOnScreenInfo struct {
	Position       Coordinates
	BoundingCorner Coordinates
	Size           Dimensions
	Handle         win.HWND
}

type WinAwait struct {
	Title         string
	AwaitTimeout  time.Duration
	SleepDuration time.Duration
	Logger        *logrus.Logger
}

func isWindow(hwnd win.HWND) bool {
	ret, _, _ := procIsWindow.Call(uintptr(hwnd))
	return ret != 0
}

func (c *WinAwait) WithTitle(title string) *WinAwait {
	c.Title = title
	return c
}

func (c *WinAwait) WithSleep(sleep time.Duration) *WinAwait {
	c.SleepDuration = sleep
	return c
}

func (c *WinAwait) WithAwaitTimeout(timeout time.Duration) *WinAwait {
	c.AwaitTimeout = timeout
	return c
}

func (c *WinAwait) AwaitOpen() (WinOnScreenInfo, error) {
	if c.Title == "" {
		return WinOnScreenInfo{}, errors.New("window title is empty")
	}
	if c.Logger == nil {
		return WinOnScreenInfo{}, errors.New("logger object not provided")
	}

	c.Logger.Info(fmt.Sprintf("Searching '%s' window", c.Title))

	winInfo, err := GetWindowInfo(c.Title, c.AwaitTimeout, c.Logger)
	if err != nil {
		return WinOnScreenInfo{}, err
	}

	maxFocusAttempts := 10
	win.SetForegroundWindow(winInfo.Handle)
	for notFocused := true; notFocused && isWindow(winInfo.Handle); notFocused = win.GetForegroundWindow() != winInfo.Handle {
		if maxFocusAttempts -= 1; maxFocusAttempts <= 0 {
			return WinOnScreenInfo{}, errors.New(fmt.Sprintf("Failed to focus '%s' window", c.Title))
		}
		time.Sleep(c.SleepDuration)
		win.SetForegroundWindow(winInfo.Handle)
	}

	return winInfo, nil
}

func (c *WinAwait) AwaitClose() error {
	if c.Title == "" {
		return errors.New("window title is empty")
	}
	if c.Logger == nil {
		return errors.New("logger object not provided")
	}
	c.Logger.Info(fmt.Sprintf("Searching '%s' window", c.Title))

	isClosed, err := AwaitWindowClosed(c.Title, c.AwaitTimeout, c.Logger)
	if err != nil {
		return err
	} else if isClosed {
		return nil
	}
	return errors.New("window not closed")
}

func FocusToProcess(processName string, logger *logrus.Logger) (int, error) {
	logger.Debugln("Searching for process name: ", processName)

	processIds, err := robotgo.FindIds(processName)
	if err != nil {
		return -1, err
	}
	pid := processIds[0]

	err = robotgo.ActivePid(pid)
	if err != nil {
		return -1, err
	}

	logger.Debugln("Process found with PID: ", pid)

	return pid, nil
}

func AwaitWindowClosed(title string, maxWaitDuration time.Duration, logger *logrus.Logger) (bool, error) {
	logger.Debugln(fmt.Sprintf("Awaiting for window: '%s' to close (max wait duration: %s)", title, maxWaitDuration))

	asyncResult := make(chan AsyncResult[win.HWND])
	maxWaitTimer := time.After(maxWaitDuration)

GetWindowInfo:
	go func(output chan<- AsyncResult[win.HWND]) {
		output <- AsyncResult[win.HWND]{
			Value: robotgo.FindWindow(title),
		}
	}(asyncResult)

	select {
	case r := <-asyncResult:
		if r.Value <= win.HWND(0) {
			return true, nil
		}
		time.Sleep(1 * time.Second)
		goto GetWindowInfo
	case <-maxWaitTimer:
		return false, fmt.Errorf("timeout exceeded when awaiting window: '%s' to close", title)
	}
}

func GetWindowInfo(title string, maxWaitDuration time.Duration, logger *logrus.Logger) (WinOnScreenInfo, error) {
	logger.Debugln(fmt.Sprintf("Awaiting for window: '%s' (max wait duration: %s)", title, maxWaitDuration))

	asyncResult := make(chan AsyncResult[WinOnScreenInfo])
	go AwaitWindow(title, asyncResult)
	maxWaitTimer := time.After(maxWaitDuration)

	select {
	case r := <-asyncResult:
		if r.Error != nil {
			return WinOnScreenInfo{}, r.Error
		}

		logger.Debugln(fmt.Sprintf("Window found with handle: %d, coords: %+v, size: %+v",
			r.Value.Handle, r.Value.Position, r.Value.Size))

		return r.Value, nil
	case <-maxWaitTimer:
		return WinOnScreenInfo{}, fmt.Errorf("timeout exceeded when awaiting window: '%s'", title)
	}
}

func AwaitWindow(title string, out chan<- AsyncResult[WinOnScreenInfo]) {
	var winHandle = win.HWND(0)
	for {
		winHandle = robotgo.FindWindow(title)
		if winHandle > win.HWND(0) {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	if !(winHandle > win.HWND(0)) {
		out <- AsyncResult[WinOnScreenInfo]{Error: errors.New("window handle not found")}
	}

	winRect := win.RECT{}
	win.GetWindowRect(winHandle, &winRect)

	winInfo := WinOnScreenInfo{}
	winInfo.Handle = winHandle
	winInfo.Size = Dimensions{
		W: int(winRect.Right - winRect.Left),
		H: int(winRect.Bottom - winRect.Top),
	}
	winInfo.Position = Coordinates{
		X: int(winRect.Left),
		Y: int(winRect.Top),
	}
	winInfo.BoundingCorner = Coordinates{
		X: int(winRect.Right),
		Y: int(winRect.Bottom),
	}

	out <- AsyncResult[WinOnScreenInfo]{
		Value: winInfo,
		Error: nil,
	}
}
