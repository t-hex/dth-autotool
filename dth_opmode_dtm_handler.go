package main

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DazToMayaExportDialogHandler struct {
	DthAutoToolOpModeBase
}

func NewDazToMayaExportDialogHandler(logger *logrus.Logger, cfg DthAutoToolConfig) DazToMayaExportDialogHandler {
	opMode := DazToMayaExportDialogHandler{}
	opMode.SetLogger(logger)
	opMode.SetConfiguration(cfg)
	return opMode
}

func (o DazToMayaExportDialogHandler) ValidateHandlingMethod(code HandlingMethodCode) (HandlingMethodCode, error) {
	switch code {
	case HandlingMethodKeySequencing, HandlingMethodVisualLookup:
		return code, nil
	default:
		return HandlingMethodUndefined, errors.New("unsupported handling method")
	}
}

func (o DazToMayaExportDialogHandler) StepSetAssetName(
	method HandlingMethodCode,
	dzToMayaExpWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Setting up 'Asset Name'"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 3 {
			err := robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}
		robotgo.TypeStr(o.cfg.DazStudio.DazToMayaExportHandler.AssetName)
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_MainExportOptionsGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzToMayaExpWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_AssetNameLbl, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Right,
			PostCaptureOffset:    Coordinates{10, 0},
		}}, dzToMayaExpWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("\"Asset Name\" input text box found at: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click("left", true)
		robotgo.TypeStr(o.cfg.DazStudio.DazToMayaExportHandler.AssetName)
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o DazToMayaExportDialogHandler) StepSetAssetType(
	method HandlingMethodCode,
	dzToMayaExpWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Setting up 'Asset Type'"))

	switch method {
	case HandlingMethodKeySequencing:
		err := robotgo.KeyTap(robotgo.Tab)
		if err != nil {
			return err
		}
		// while on asset type drop-down, we press 'A' to auto-select "animation" type
		err = robotgo.KeyTap(robotgo.KeyA)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_MainExportOptionsGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzToMayaExpWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_AssetTypeLbl, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Right,
			PostCaptureOffset:    Coordinates{10, 0},
		}}, dzToMayaExpWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
		// while on asset type drop-down, we press 'ESC' to close dropdown
		err = robotgo.KeyTap(robotgo.Esc)
		if err != nil {
			return err
		}

		// while on asset type drop-down, we press 'A' to auto-select "animation" type
		err = robotgo.KeyTap(robotgo.KeyA)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o DazToMayaExportDialogHandler) StepOpenSubdivisionsDialog(
	method HandlingMethodCode,
	dzToMayaExpWinInfo WinOnScreenInfo) error {

	o.logger.Info(fmt.Sprintf("Opening subdivisions setup dialog"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 4 {
			err := robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}

		err := robotgo.KeyTap(robotgo.Space)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_MainExportOptionsGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzToMayaExpWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_SubdLevelsBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
		}}, dzToMayaExpWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("\"Bake Subdivision Levels\" button found at: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		err := errors.New("unsupported handling method")
		return err
	}

	return nil
}

func (o DazToMayaExportDialogHandler) StepSetSubdivisionLevels(
	method HandlingMethodCode,
	dzSubdivisionWinInfo WinOnScreenInfo) error {
	o.logger.Info("Setting up subdivision levels")

	switch method {
	case HandlingMethodKeySequencing:
		// focus to the top drop-down
		for range 2 {
			err := robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}

		// iterate over all drop-downs and set subdivision levels to 1
		for range o.cfg.DazStudio.DazToMayaExportHandler.ExpectedSubdivisionShapes {
			err := robotgo.KeyTap(robotgo.Key1)
			if err != nil {
				return err
			}
			err = robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		// focus to the top drop-down
		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_SubdLevelsGrpTop, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: BottomRight,
			BoundingCornerOffset: dzSubdivisionWinInfo.Size,
			PostCaptureOffset:    Coordinates{-10, -10},
		}}, dzSubdivisionWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
		err = robotgo.KeyTap(robotgo.Esc)
		if err != nil {
			return err
		}

		// iterate over all drop-downs and set subdivision levels to 1
		witBtnSearchTraces, err := ScreenPositionSearchTrace([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_WhatIsThisBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzSubdivisionWinInfo.Size,
		}}, dzSubdivisionWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}
		// cache button right after the last subdivision drop-down
		witBtnFinalMatchTrace := witBtnSearchTraces[len(witBtnSearchTraces)-1]
		witBtnVisualTracker, err := NewScreenAreaVisualTracker(ScreenArea{
			Offset: witBtnFinalMatchTrace.OnScreenPosition,
			Size:   witBtnFinalMatchTrace.MatchedPatternSize,
		})
		if err != nil {
			return err
		}

		const maxIterations = 100 // to prevent infinite loop in case of error
		for i, c := 0, float32(1.0); i < maxIterations && c >= 0.95; i++ {
			o.logger.Debug(fmt.Sprintf("Changing subdivision level entry [%d]", i+1))
			err = robotgo.KeyTap(robotgo.Key1)
			if err != nil {
				return err
			}
			err = robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
			c, err = witBtnVisualTracker.Refresh()
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o DazToMayaExportDialogHandler) StepAcceptSubdivisionsDialog(
	method HandlingMethodCode,
	dzSubdivisionWinInfo WinOnScreenInfo) error {
	o.logger.Info("Confirm subdivision levels")

	switch method {
	case HandlingMethodKeySequencing:
		err := robotgo.KeyTap(robotgo.Tab)
		if err != nil {
			return err
		}
		err = robotgo.KeyTap(robotgo.Space)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_AcceptCancelBtnGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzSubdivisionWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_CommonHandler_AcceptBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
		}}, dzSubdivisionWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o DazToMayaExportDialogHandler) StepAcceptDazToMayaExportDialog(
	method HandlingMethodCode,
	dzToMayaExportWinInfo WinOnScreenInfo) error {
	o.logger.Info("Confirm asset export")

	switch method {
	case HandlingMethodKeySequencing:
		for range 8 {
			err := robotgo.KeyTap(robotgo.Tab, robotgo.Shift)
			if err != nil {
				return err
			}
		}
		err := robotgo.KeyTap(robotgo.Space)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_DtmHandler_AcceptCancelBtnGrp, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzToMayaExportWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				*NewSSIRP(RefPattern_CommonHandler_AcceptBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
		}}, dzToMayaExportWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("Accept button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}

	o.logger.Info("Exporting asset...")

	return nil
}

func (o DazToMayaExportDialogHandler) StepAcceptExportFinishedDialog(
	method HandlingMethodCode,
	dzToMayaExportFinishedWinInfo WinOnScreenInfo) error {

	o.logger.Info("Confirming export finished")

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
				*NewSSIRP(RefPattern_CommonHandler_OkBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
			BoundingCornerOffset: dzToMayaExportFinishedWinInfo.Size,
		}}, dzToMayaExportFinishedWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("OK button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}

	return nil
}

func (o DazToMayaExportDialogHandler) StepAcceptExportObjectBakingDialog(
	method HandlingMethodCode,
	dzToMayaExportObjectBakingWinInfo WinOnScreenInfo) error {

	o.logger.Info("Confirming export object baking")

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
				*NewSSIRP(RefPattern_CommonHandler_YesBtn, o.cfg.ImgPatternsPath),
			},
			PostCaptureAlignment: Center,
			BoundingCornerOffset: dzToMayaExportObjectBakingWinInfo.Size,
		}}, dzToMayaExportObjectBakingWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("Yes button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}

	return nil
}

func (o DazToMayaExportDialogHandler) Run() error {
	method, err := o.ValidateHandlingMethod(o.cfg.DazStudio.LoadAssetResizeTimelineHandler.HandlingMethod)
	o.logger.Info(fmt.Sprintf("Running DazToMayaExportDialogHandler [method:%d]", method))

	dzPid, err := FocusToProcess(o.cfg.DazStudio.ProcessName, o.logger)
	if err != nil {
		return err
	}
	o.logger.Info("Process found with PID: ", dzPid)

	winAwait := WinAwait{
		SleepDuration: 1 * time.Second,
		Logger:        o.logger,
	}

	dzToMayaExpWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.DazToMayaExportHandler.WindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.DazToMayaExportHandler.WindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepSetAssetName(method, dzToMayaExpWinInfo)
	if err != nil {
		return err
	}

	err = o.StepSetAssetType(method, dzToMayaExpWinInfo)
	if err != nil {
		return err
	}

	err = o.StepOpenSubdivisionsDialog(method, dzToMayaExpWinInfo)
	if err != nil {
		return err
	}

	dzSubdivisionWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.DazToMayaExportHandler.SubDivisionWindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.DazToMayaExportHandler.SubDivisionWindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepSetSubdivisionLevels(method, dzSubdivisionWinInfo)
	if err != nil {
		return err
	}

	err = o.StepAcceptSubdivisionsDialog(method, dzSubdivisionWinInfo)
	if err != nil {
		return err
	}

	dzToMayaExpWinInfo, err = winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.DazToMayaExportHandler.WindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.DazToMayaExportHandler.WindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepAcceptDazToMayaExportDialog(method, dzToMayaExpWinInfo)
	if err != nil {
		return err
	}

	dzToMayaExpWinInfoChan := make(chan func() (WinOnScreenInfo, error))
	dzToMayaObjectBakingWinInfoChan := make(chan func() (WinOnScreenInfo, error))

	// wait for either "baking object" prompt or straightly "export finished" window
	go func() {
		wa := WinAwait{
			SleepDuration: 1 * time.Second,
			Logger:        o.logger,
		}
		wi, err := wa.WithAwaitTimeout(time.Minute * time.Duration(o.cfg.DazStudio.DazToMayaExportHandler.MaxExportDurationMinutes)).
			WithTitle(o.cfg.DazStudio.DazToMayaExportHandler.ExpFinishedWindowTitle).
			AwaitOpen()
		dzToMayaExpWinInfoChan <- func() (WinOnScreenInfo, error) { return wi, err }
	}()
	go func() {
		wa := WinAwait{
			SleepDuration: 1 * time.Second,
			Logger:        o.logger,
		}
		wi, err := wa.
			WithAwaitTimeout(time.Minute * time.Duration(o.cfg.DazStudio.DazToMayaExportHandler.ExpObjectBakingWindowMaxWaitSeconds)).
			WithTitle(o.cfg.DazStudio.DazToMayaExportHandler.ExpObjectBakingWindowTitle).
			AwaitOpen()
		dzToMayaObjectBakingWinInfoChan <- func() (WinOnScreenInfo, error) { return wi, err }
	}()

AwaitExportFinishedOrObjectBakingPrompt:
	select {
	case dzToMayaObjectBakingWinInfoGetter := <-dzToMayaObjectBakingWinInfoChan:
		dzToMayaExpObjectBakingWinInfo, err := dzToMayaObjectBakingWinInfoGetter()
		if err != nil && !strings.Contains(err.Error(), "timeout") {
			return err
		}
		err = o.StepAcceptExportObjectBakingDialog(method, dzToMayaExpObjectBakingWinInfo)
		if err != nil {
			return err
		}
		goto AwaitExportFinishedOrObjectBakingPrompt
	case dzToMayaExpWinInfoGetter := <-dzToMayaExpWinInfoChan:
		dzToMayaExpWinInfo, err = dzToMayaExpWinInfoGetter()
		// continue with confirmation without awaiting "baking object" prompt
	}

	if err != nil {
		return err
	}

	err = o.StepAcceptExportFinishedDialog(method, dzToMayaExpWinInfo)
	if err != nil {
		return err
	}

	return nil
}
