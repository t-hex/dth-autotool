package main

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type SaganAlembicExportDialogHandler struct {
	DthAutoToolOpModeBase
}

func NewSaganAlembicExportDialogHandler(logger *logrus.Logger, cfg DthAutoToolConfig) SaganAlembicExportDialogHandler {
	opMode := SaganAlembicExportDialogHandler{}
	opMode.SetLogger(logger)
	opMode.SetConfiguration(cfg)
	return opMode
}

func (o SaganAlembicExportDialogHandler) ValidateHandlingMethod(code HandlingMethodCode) (HandlingMethodCode, error) {
	switch code {
	case HandlingMethodKeySequencing, HandlingMethodVisualLookup:
		return code, nil
	default:
		return HandlingMethodUndefined, errors.New("unsupported handling method")
	}
}

func (o SaganAlembicExportDialogHandler) StepLoadSaganAlembicExportConfig(
	method HandlingMethodCode,
	dzSaganAlembicExportWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Loading 'Sagan Alembic Export' configuration"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 7 {
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
		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp.png",
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp_load_highlighted.png",
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp_save_highlighted.png",
				}},
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzSaganAlembicExportWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/load_cfg_btn.png",
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/load_cfg_btn_highlighted.png"},
					ValidationText: "Load Config"},
			},
			PostCaptureAlignment: Center,
		}}, dzSaganAlembicExportWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("Accept button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o SaganAlembicExportDialogHandler) StepSelectSaganAlembicExportConfig(
	method HandlingMethodCode,
	dzSaganAlembicLoadCfgWinInfo WinOnScreenInfo,
	configFileAbsPath string) error {
	o.logger.Info(fmt.Sprintf("Selecting 'Sagan Alembic Export' configuration file'"))

	switch method {
	case HandlingMethodKeySequencing:
		robotgo.TypeStr(configFileAbsPath)
		err := robotgo.KeyTap(robotgo.Enter)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{o.cfg.ImgPatternsPath + "/common/win10_dlg_filename_lbl.png"}},
			},
			PostCaptureAlignment: Right,
			PostCaptureOffset:    Coordinates{20, 0},
			BoundingCornerOffset: dzSaganAlembicLoadCfgWinInfo.Size,
		}}, dzSaganAlembicLoadCfgWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("\"File Name\" input text box found at: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
		robotgo.TypeStr(configFileAbsPath)

		position, err = ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/common/win10_dlg_open_cancel_btn_grp.png",
					o.cfg.ImgPatternsPath + "/common/win11_dlg_open_cancel_btn_grp.png",
					o.cfg.ImgPatternsPath + "/common/win10_dlg_open_cancel_btn_grp_open_highlighted.png",
					o.cfg.ImgPatternsPath + "/common/win11_dlg_open_cancel_btn_grp_open_highlighted.png",
					o.cfg.ImgPatternsPath + "/common/win10_dlg_open_cancel_btn_grp_cancel_highlighted.png",
					o.cfg.ImgPatternsPath + "/common/win11_dlg_open_cancel_btn_grp_cancel_highlighted.png",
				}},
			},
			PostCaptureAlignment: TopLeft,
			BoundingCornerOffset: dzSaganAlembicLoadCfgWinInfo.Size,
		}, {
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/common/win10_dlg_open_btn.png",
					o.cfg.ImgPatternsPath + "/common/win11_dlg_open_btn.png",
					o.cfg.ImgPatternsPath + "/common/win10_dlg_open_btn_highlighted.png",
					o.cfg.ImgPatternsPath + "/common/win11_dlg_open_btn_highlighted.png",
				}, ValidationText: "Open"},
			},
			PostCaptureAlignment: Center,
		}}, dzSaganAlembicLoadCfgWinInfo.Position, o.cfg.ScreenSearchValidation)
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

func (o SaganAlembicExportDialogHandler) StepGenerateTemporaryExportConfig() (string, error) {
	o.logger.Info("Setting up timeline range ")

	cfg, err := ini.Load(o.cfg.DazStudio.SaganAlembicExportHandler.SaganConfigFilePath)
	if err != nil {
		return "", err
	}

	tempCfgAbsFilePath, err := filepath.Abs("./.temp.export.sagan")
	if err != nil {
		return "", err
	}

	cfg.Section("General").
		Key("end_frame").
		SetValue(strconv.Itoa(int(o.cfg.DazStudio.SaganAlembicExportHandler.TimelineEndFrame)))
	cfg.Section("General").
		Key("base_directory").
		SetValue(o.cfg.DazStudio.SaganAlembicExportHandler.OutputPath)

	err = cfg.SaveTo(tempCfgAbsFilePath)
	if err != nil {
		return "", err
	}

	return tempCfgAbsFilePath, nil
}

func (o SaganAlembicExportDialogHandler) StepSetupTimelineRange(
	method HandlingMethodCode,
	dzSaganAlembicExportWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Setting up timeline range 0-%d", o.cfg.DazStudio.SaganAlembicExportHandler.TimelineEndFrame))

	switch method {
	case HandlingMethodKeySequencing:
		for range 7 {
			err := robotgo.KeyTap(robotgo.Tab)
			if err != nil {
				return err
			}
		}

		robotgo.TypeStr(strconv.Itoa(int(o.cfg.DazStudio.SaganAlembicExportHandler.TimelineEndFrame)))
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/sagan_alembic_exporter_dialog/end_frame_lbl.png",
				}, ValidationText: "End Frame"},
			},
			PostCaptureAlignment: Right,
			PostCaptureOffset:    Coordinates{20, 0},
			BoundingCornerOffset: dzSaganAlembicExportWinInfo.Size,
		}}, dzSaganAlembicExportWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("\"End Frame\" input range box found at: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click(robotgo.Mleft, true)
		robotgo.TypeStr(strconv.Itoa(int(o.cfg.DazStudio.SaganAlembicExportHandler.TimelineEndFrame)))
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o SaganAlembicExportDialogHandler) StepAcceptSaganAlembicExport(
	method HandlingMethodCode,
	dzSaganAlembicExportWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Starting 'Sagan Alembic Export' process"))

	switch method {
	case HandlingMethodKeySequencing:
		for range 2 {
			err := robotgo.KeyTap(robotgo.Tab)
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
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/common/export_btn.png",
					o.cfg.ImgPatternsPath + "/common/export_btn_highlighted.png"},
					ValidationText: "Export"},
			},
			PostCaptureAlignment: Center,
			BoundingCornerOffset: dzSaganAlembicExportWinInfo.Size,
		}}, dzSaganAlembicExportWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("Export button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o SaganAlembicExportDialogHandler) StepCloseSaganAlembicExportDialog(
	method HandlingMethodCode,
	dzSaganAlembicExportWinInfo WinOnScreenInfo) error {
	o.logger.Info(fmt.Sprintf("Closing 'Sagan Alembic Export' process dialog"))

	switch method {
	case HandlingMethodKeySequencing:
		err := robotgo.KeyTap(robotgo.Tab)
		if err != nil {
			return err
		}
		err = robotgo.KeyTap(robotgo.Enter)
		if err != nil {
			return err
		}
	case HandlingMethodVisualLookup:
		robotgo.Move(robotgo.GetScreenSize()) // move cursor out of the way

		position, err := ScreenPositionSearch([]ScreenPositionSearchDefinition{{
			Patterns: []ScreenSearchImgReferencePattern{
				{ImageFilePaths: []string{
					o.cfg.ImgPatternsPath + "/common/done_btn.png",
					o.cfg.ImgPatternsPath + "/common/done_btn_highlighted.png"},
					ValidationText: "Done"},
			},
			PostCaptureAlignment: Center,
			BoundingCornerOffset: dzSaganAlembicExportWinInfo.Size,
		}}, dzSaganAlembicExportWinInfo.Position, o.cfg.ScreenSearchValidation)
		if err != nil {
			return err
		}

		o.logger.Debug("Done button found: ", position)

		robotgo.Move(position.X, position.Y)
		robotgo.Click()
	default:
		return errors.New("unsupported handling method")
	}
	return nil
}

func (o SaganAlembicExportDialogHandler) Run() error {
	method, err := o.ValidateHandlingMethod(o.cfg.DazStudio.SaganAlembicExportHandler.HandlingMethod)
	o.logger.Info(fmt.Sprintf("Running SaganAlembicExportDialogHandler [method:%d]", method))

	dzPid, err := FocusToProcess(o.cfg.DazStudio.ProcessName, o.logger)
	if err != nil {
		return err
	}
	o.logger.Info("Process found with PID: ", dzPid)

	winAwait := WinAwait{
		SleepDuration: 1 * time.Second,
		Logger:        o.logger,
	}

	dzSaganAlembicExportWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.SaganAlembicExportHandler.WindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.SaganAlembicExportHandler.WindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	err = o.StepLoadSaganAlembicExportConfig(method, dzSaganAlembicExportWinInfo)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second) // file dialog may take longer to open and adjust dimensions

	dzSaganAlembicLoadCfgWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.SaganAlembicExportHandler.LoadSaganConfigWindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.SaganAlembicExportHandler.LoadSaganConfigWindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}

	cfgAbsPath, err := o.StepGenerateTemporaryExportConfig()
	if err != nil {
		return err
	}
	defer os.Remove(cfgAbsPath)

	err = o.StepSelectSaganAlembicExportConfig(method, dzSaganAlembicLoadCfgWinInfo, cfgAbsPath)
	if err != nil {
		return err
	}

	//deprecated: kept just for a reference
	//err = o.StepSetupTimelineRange(method, dzSaganAlembicLoadCfgWinInfo)
	//if err != nil {
	//	return err
	//}

	err = o.StepAcceptSaganAlembicExport(method, dzSaganAlembicExportWinInfo)
	if err != nil {
		return err
	}

	dzExpProgressWinInfo, err := winAwait.
		WithAwaitTimeout(time.Second * time.Duration(o.cfg.DazStudio.SaganAlembicExportHandler.ExpProgressWindowMaxWaitSeconds)).
		WithTitle(o.cfg.DazStudio.SaganAlembicExportHandler.ExpProgressWindowTitle).
		AwaitOpen()
	if err != nil {
		return err
	}
	o.logger.Info(fmt.Sprintf("Waiting for export window to finish [%d]", dzExpProgressWinInfo.Handle))
	err = winAwait.
		WithAwaitTimeout(time.Minute * time.Duration(o.cfg.DazStudio.SaganAlembicExportHandler.MaxExportDurationMinutes)).
		AwaitClose()
	if err != nil {
		return err
	}

	err = o.StepCloseSaganAlembicExportDialog(method, dzSaganAlembicExportWinInfo)
	if err != nil {
		return err
	}

	return nil
}
