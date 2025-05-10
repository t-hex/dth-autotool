package main

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"time"
)

type LoadCharacterDialogHandler struct {
	DthAutoToolOpModeBase
}

func NewLoadCharacterDialogHandler(logger *logrus.Logger, cfg DthAutoToolConfig) LoadCharacterDialogHandler {
	opMode := LoadCharacterDialogHandler{}
	opMode.SetLogger(logger)
	opMode.SetConfiguration(cfg)
	return opMode
}

func (o LoadCharacterDialogHandler) ValidateHandlingMethod(code HandlingMethodCode) (HandlingMethodCode, error) {
	switch code {
	case HandlingMethodKeySequencing, HandlingMethodVisualLookup:
		return code, nil
	default:
		return HandlingMethodUndefined, errors.New("unsupported handling method")
	}
}

func (o LoadCharacterDialogHandler) StepApplyCharacterToCurrentSelection(
	method HandlingMethodCode,
	dzLoadCharacterOptionsWinInfo WinOnScreenInfo) error {
	o.logger.Infoln(fmt.Sprintf("Choosing to apply character to the current selection"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 2 {
			err := robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}

		// reset radio button to the first item by pressing down & up while focused
		err := robotgo.KeyTap(robotgo.Down)
		if err != nil {
			return err
		}
		err = robotgo.KeyTap(robotgo.Up)
		if err != nil {
			return err
		}

		// select option by its order from the top
		const rBtnGrpOptOrder = 2
		for range rBtnGrpOptOrder - 1 {
			err := robotgo.KeyTap(robotgo.Down)
			if err != nil {
				return err
			}
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_LchHandler_ApplyToSelectedRadioBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzLoadCharacterOptionsWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_CommonHandler_RadioBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
		}}, dzLoadCharacterOptionsWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debugln("Radio button 'Apply this Character to the currently selected Figure(s)' found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o LoadCharacterDialogHandler) StepAcceptCharacterLoadingOptions(
	method HandlingMethodCode,
	dzLoadCharacterOptionsWinInfo WinOnScreenInfo) error {
	o.logger.Infoln(fmt.Sprintf("Accepting character loading options"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 2 {
			err := robotgo.KeyTap(robotgo.Tab, robotgo.Lshift)
			if err != nil {
				return err
			}
		}
		err := robotgo.KeyTap(robotgo.Enter)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_CommonHandler_AcceptBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
			BoundingCornerOffset: dzLoadCharacterOptionsWinInfo.Size,
		}}, dzLoadCharacterOptionsWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debugln("Accept button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o LoadCharacterDialogHandler) Run() error {
	method, err := o.ValidateHandlingMethod(o.cfg.DazStudio.LoadCharacterHandler.HandlingMethod)
	o.logger.Infoln(fmt.Sprintf("Running LoadCharacterDialogHandler [method:%d]", method))

	dzPid, err := FocusToProcess(o.cfg.DazStudio.ProcessName, o.logger)
	if err != nil {
		return err
	}
	o.logger.Infoln("Process found with PID: ", dzPid)

	winAwait := WinAwait{
		SleepDuration: 1 * time.Second,
		Logger:        o.logger,
	}

	dzLoadCharacterOptionsWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.LoadCharacterHandler.WindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.LoadCharacterHandler.WindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepApplyCharacterToCurrentSelection(method, dzLoadCharacterOptionsWinInfo)
	if err != nil {
		return err
	}

	err = o.StepAcceptCharacterLoadingOptions(method, dzLoadCharacterOptionsWinInfo)
	if err != nil {
		return err
	}

	return nil
}
