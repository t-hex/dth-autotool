package main

import (
	"errors"
	"regexp"
	"strings"
)

type DthAutoToolOpModeName string

const (
	OpModeHandleResizeTimelineDialog  DthAutoToolOpModeName = "AutoConfirm/ResizeTimeline"
	OpModeHandleLoadCharacterDialog   DthAutoToolOpModeName = "AutoConfirm/LoadCharacter"
	OpModeHandleAlembicExportDialog   DthAutoToolOpModeName = "AlembicExport/HandleDialog"
	OpModeHandleDazToMayaExportDialog DthAutoToolOpModeName = "DazToMaya/HandleDialog"
)

type HandlingMethodCode uint32

const (
	HandlingMethodUndefined HandlingMethodCode = iota
	HandlingMethodKeySequencing
	HandlingMethodVisualLookup
)

func ParseHandlingMethod(method string) (HandlingMethodCode, error) {
	switch regexp.MustCompile(`[\s_-]`).ReplaceAllString(strings.ToLower(method), "") {
	case "keysequencing":
		return HandlingMethodKeySequencing, nil
	case "visuallookup":
		return HandlingMethodVisualLookup, nil
	default:
		return HandlingMethodUndefined, errors.New("invalid handling method")
	}
}

func (d *HandlingMethodCode) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	value := strings.Trim(string(data), "\"")
	method, err := ParseHandlingMethod(value)
	if err == nil {
		*d = method
	}
	return nil // keep default method (skip this property parsing)
}

var DthAutoToolCmdOptions struct {
	OpMode          DthAutoToolOpModeName `short:"m" long:"operation-mode"`
	AppConfigPath   string                `short:"c" long:"configuration"`
	LogLevel        string                `long:"loglevel"`
	AssetName       string                `short:"n" long:"asset-name"`
	ShapesCount     uint32                `short:"s" long:"shapes-count"`
	EndFrame        uint32                `short:"e" long:"end-frame"`
	SaganConfigPath string                `short:"a" long:"sagan-alembic-export-config" default:"template.config.sagan"`
	SaganOutputPath string                `short:"o" long:"sagan-alembic-export-output"`
}

type DthAutoToolDazToMayaExportHandlerConfig struct {
	HandlingMethod                      HandlingMethodCode `json:"handling_method"`
	WindowTitle                         string             `json:"window_title"`
	WindowMaxWaitSeconds                int                `json:"window_max_wait_duration_s"`
	SubDivisionWindowTitle              string             `json:"subdivision_window_title"`
	SubDivisionWindowMaxWaitSeconds     int                `json:"subdivision_window_max_wait_duration_s"`
	ExpFinishedWindowTitle              string             `json:"export_finished_window_title"`
	MaxExportDurationMinutes            int                `json:"max_export_duration_m"`
	ExpObjectBakingWindowTitle          string             `json:"object_baking_window_title"`
	ExpObjectBakingWindowMaxWaitSeconds int                `json:"object_baking_window_max_wait_duration_s"`
	ExpectedSubdivisionShapes           uint32
	AssetName                           string
}

type DthAutoToolSaganAlembicExportHandlerConfig struct {
	HandlingMethod                      HandlingMethodCode `json:"handling_method"`
	WindowTitle                         string             `json:"window_title"`
	WindowMaxWaitSeconds                int                `json:"window_max_wait_duration_s"`
	LoadSaganConfigWindowTitle          string             `json:"load_sagan_config_window_title"`
	LoadSaganConfigWindowMaxWaitSeconds int                `json:"load_sagan_config_window_max_wait_duration_s"`
	ExpProgressWindowTitle              string             `json:"export_progress_window_title"`
	ExpProgressWindowMaxWaitSeconds     int                `json:"export_progress_window_max_wait_duration_s"`
	MaxExportDurationMinutes            int                `json:"max_export_duration_m"`
	SaganConfigFilePath                 string
	OutputPath                          string
	TimelineEndFrame                    uint32
}

type DthAutoToolLoadAssetResizeTimelineHandlerConfig struct {
	HandlingMethod       HandlingMethodCode `json:"handling_method"`
	WindowTitle          string             `json:"window_title"`
	WindowMaxWaitSeconds int                `json:"window_max_wait_duration_s"`
}

type DthAutoToolLoadCharacterHandlerConfig struct {
	HandlingMethod       HandlingMethodCode `json:"handling_method"`
	WindowTitle          string             `json:"window_title"`
	WindowMaxWaitSeconds int                `json:"window_max_wait_duration_s"`
}

type DthAutoToolDazStudioConfig struct {
	ProcessName                    string                                          `json:"process_name"`
	DazToMayaExportHandler         DthAutoToolDazToMayaExportHandlerConfig         `json:"daz_to_maya_export_handler"`
	SaganAlembicExportHandler      DthAutoToolSaganAlembicExportHandlerConfig      `json:"sagan_alembic_export_handler"`
	LoadAssetResizeTimelineHandler DthAutoToolLoadAssetResizeTimelineHandlerConfig `json:"load_asset_resize_timeline_handler"`
	LoadCharacterHandler           DthAutoToolLoadCharacterHandlerConfig           `json:"load_character_asset_handler"`
}

type JaroWrinklerValidationConfig struct {
	Enabled           bool    `json:"enabled"`
	BoostThreshold    float64 `json:"boost_threshold"`
	DistanceThreshold float64 `json:"distance_threshold"`
	PrefixSize        int     `json:"prefix_size"`
}

type TesseractValidationConfig struct {
	Enabled                bool                         `json:"enabled"`
	TesseractPSM           TesseractPSM                 `json:"psm"`
	TesseractLang          TesseractLang                `json:"lang"`
	JaroWrinklerValidation JaroWrinklerValidationConfig `json:"jaro_wrinkler_validation"`
}

type ScreenPositionSearchValidationConfig struct {
	PatternMaxRetryAttempts          int                       `json:"pattern_max_retry_attempts"`
	PatternRetryAttemptDelayMs       int                       `json:"pattern_retry_attempts_delay_ms"`
	GrayScalePatternsEnabled         bool                      `json:"gray_scale_patterns_enabled"`
	GrayScalePatternsLuminanceLevels uint8                     `json:"gray_scale_patterns_luminance_levels"`
	TesseractValidation              TesseractValidationConfig `json:"tesseract_validation"`
}

type DthAutoToolConfig struct {
	ImgPatternsPath        string                               `json:"img_patterns_path"`
	MouseSleepMs           int                                  `json:"mouse_sleep_ms"`
	KeySleepMs             int                                  `json:"key_sleep_ms"`
	DazStudio              DthAutoToolDazStudioConfig           `json:"daz_studio"`
	ScreenSearchValidation ScreenPositionSearchValidationConfig `json:"screen_search_validation"`
}
