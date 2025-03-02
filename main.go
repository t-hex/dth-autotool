package main

import (
	"encoding/json"
	"github.com/go-vgo/robotgo"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	_, err := flags.ParseArgs(&DthAutoToolCmdOptions, os.Args)
	if err != nil {
		panic(err)
	}

	var logger = logrus.New()
	if DthAutoToolCmdOptions.LogLevel != "" {
		logLevel, err := logrus.ParseLevel(DthAutoToolCmdOptions.LogLevel)
		if err != nil {
			panic(err)
		}
		logger.SetLevel(logLevel)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupt
		os.Exit(0)
	}()

	appConfigAbsPath, err := filepath.Abs(DthAutoToolCmdOptions.AppConfigPath)
	if err != nil {
		logger.Fatal(err)
	}
	cfgFile, err := os.Open(appConfigAbsPath)
	if err != nil {
		logger.Fatal(err)
	}

	// configuration defaults
	var cfg = DthAutoToolConfig{
		MouseSleepMs: 10,
		KeySleepMs:   10,
		DazStudio: DthAutoToolDazStudioConfig{
			ProcessName: "dazstudio",
			SaganAlembicExportHandler: DthAutoToolSaganAlembicExportHandlerConfig{
				HandlingMethod:                      HandlingMethodKeySequencing,
				SaganConfigFilePath:                 DthAutoToolCmdOptions.SaganConfigPath,
				WindowTitle:                         "Sagan Alembic Exporter v3",
				WindowMaxWaitSeconds:                10,
				LoadSaganConfigWindowTitle:          "Select Config File",
				LoadSaganConfigWindowMaxWaitSeconds: 10,
				ExpProgressWindowTitle:              "Export Progress",
				ExpProgressWindowMaxWaitSeconds:     10,
				MaxExportDurationMinutes:            5,
				TimelineEndFrame:                    DthAutoToolCmdOptions.EndFrame,
				OutputPath:                          DthAutoToolCmdOptions.SaganOutputPath,
			},
			LoadAssetResizeTimelineHandler: DthAutoToolLoadAssetResizeTimelineHandlerConfig{
				HandlingMethod:       HandlingMethodKeySequencing,
				WindowTitle:          "Animation Range : Total Frames",
				WindowMaxWaitSeconds: 10,
			},
			LoadCharacterHandler: DthAutoToolLoadCharacterHandlerConfig{
				HandlingMethod:       HandlingMethodKeySequencing,
				WindowTitle:          "Character Loading Options",
				WindowMaxWaitSeconds: 10,
			},
			DazToMayaExportHandler: DthAutoToolDazToMayaExportHandlerConfig{
				HandlingMethod:                      HandlingMethodKeySequencing,
				WindowTitle:                         "Maya Export Options",
				WindowMaxWaitSeconds:                10,
				SubDivisionWindowTitle:              "Bake Subdivision Levels",
				SubDivisionWindowMaxWaitSeconds:     10,
				ExpFinishedWindowTitle:              "Daz To Maya Bridge",
				MaxExportDurationMinutes:            5,
				ExpObjectBakingWindowTitle:          "Object Baking Recommended",
				ExpObjectBakingWindowMaxWaitSeconds: 10,
				ExpectedSubdivisionShapes:           DthAutoToolCmdOptions.ShapesCount,
				AssetName:                           DthAutoToolCmdOptions.AssetName,
			},
		},
		ScreenSearchValidation: ScreenPositionSearchValidationConfig{
			PatternMaxRetryAttempts:          3,
			PatternRetryAttemptDelayMs:       500,
			GrayScalePatternsEnabled:         true,
			GrayScalePatternsLuminanceLevels: 5,
			TesseractValidation: TesseractValidationConfig{
				Enabled:       IsTesseractAvailable(),
				TesseractPSM:  TesseractPsm_SingleWord,
				TesseractLang: TesseractLang_Eng,
				JaroWrinklerValidation: JaroWrinklerValidationConfig{
					Enabled:           true,
					BoostThreshold:    0.7,
					DistanceThreshold: 0.75,
					PrefixSize:        4,
				},
			},
		},
	}
	func(cfg *DthAutoToolConfig) {
		jsonParser := json.NewDecoder(cfgFile)
		defer cfgFile.Close()
		err = jsonParser.Decode(cfg)
		if err != nil {
			logger.Fatal(err)
		}
	}(&cfg)

	robotgo.MouseSleep = cfg.MouseSleepMs
	robotgo.KeySleep = cfg.KeySleepMs

	logger.Out = os.Stdout
	var mode DthAutoToolOpMode

	switch DthAutoToolCmdOptions.OpMode {
	case OpModeHandleDazToMayaExportDialog:
		mode = NewDazToMayaExportDialogHandler(logger, cfg)
		break
	case OpModeHandleAlembicExportDialog:
		mode = NewSaganAlembicExportDialogHandler(logger, cfg)
		break
	case OpModeHandleResizeTimelineDialog:
		mode = NewTimelineResizeDialogHandler(logger, cfg)
		break
	case OpModeHandleLoadCharacterDialog:
		mode = NewLoadCharacterDialogHandler(logger, cfg)
		break
	default:
		logger.Fatal("Unrecognized mode")
	}

	err = mode.Run()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
