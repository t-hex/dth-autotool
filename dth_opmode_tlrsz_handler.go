package main

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"time"
)

type TimelineResizeDialogHandler struct {
	DthAutoToolOpModeBase
}

func NewTimelineResizeDialogHandler(logger *logrus.Logger, cfg DthAutoToolConfig) TimelineResizeDialogHandler {
	opMode := TimelineResizeDialogHandler{}
	opMode.SetLogger(logger)
	opMode.SetConfiguration(cfg)
	return opMode
}

func (o TimelineResizeDialogHandler) ValidateHandlingMethod(code HandlingMethodCode) (HandlingMethodCode, error) {
	switch code {
	case HandlingMethodKeySequencing, HandlingMethodVisualLookup:
		return code, nil
	default:
		return HandlingMethodUndefined, errors.New("unsupported handling method")
	}
}

func (o TimelineResizeDialogHandler) StepConfirmResizeTimeline(
	method HandlingMethodCode,
	dzTimelineResizeWinInfo WinOnScreenInfo) error {
	o.logger.Infoln(fmt.Sprintf("Confirming timeline resize"))

	switch method {
	case HandlingMethodKeySequencing:
		err := robotgo.KeyTap(robotgo.Enter)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_TlrsHandler_YesNoBtnGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzTimelineResizeWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_CommonHandler_YesBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
		}}, dzTimelineResizeWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debugln("Yes button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o TimelineResizeDialogHandler) Run() error {
	method, err := o.ValidateHandlingMethod(o.cfg.DazStudio.LoadAssetResizeTimelineHandler.HandlingMethod)
	o.logger.Infoln(fmt.Sprintf("Running TimelineResizeDialogHandler [method:%d]", method))

	dzPid, err := FocusToProcess(o.cfg.DazStudio.ProcessName, o.logger)
	if err != nil {
		return err
	}
	o.logger.Infoln("Process found with PID: ", dzPid)

	winAwait := WinAwait{
		SleepDuration: 1 * time.Second,
		Logger:        o.logger,
	}

	dzTimelineResizeWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.LoadAssetResizeTimelineHandler.WindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.LoadAssetResizeTimelineHandler.WindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepConfirmResizeTimeline(method, dzTimelineResizeWinInfo)
	if err != nil {
		return err
	}

	return nil
}
