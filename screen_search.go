package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"image"
	"image/color"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"time"

	//"github.com/agnivade/levenshtein"
	"github.com/go-vgo/robotgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/vcaesar/gcv"
	"github.com/xrash/smetrics"
)

type TesseractLang string

const (
	TesseractLang_Eng = "eng"
)

type TesseractPSM int // Tesseract's page segmentation model (PSM)
const (
	TesseractPsm_SingleLine TesseractPSM = 7  // Treat the image as a single text line.
	TesseractPsm_SingleWord              = 8  // Treat the image as a single word.
	TesseractPsm_SingleChar              = 10 // Treat the image as a single character.
	TesseractPsm_RawLine                 = 13 // Raw line. Treat the image as a single text line, bypassing hacks that are Tesseract-specific.
)

type Coordinates struct {
	X, Y int
}

type Dimensions struct {
	W, H int
}

type ScreenArea struct {
	Offset Coordinates
	Size   Dimensions
}

type Alignment uint64

const (
	Center Alignment = iota
	TopLeft
	Top
	TopRight
	BottomLeft
	Bottom
	BottomRight
	Left
	Right
)

type ScreenSearchImgReferencePattern struct {
	ImageFilePaths []string
	ValidationText string
	Psm            TesseractPSM
	Language       TesseractLang
}

type ScreenPositionSearchDefinition struct {
	Patterns             []ScreenSearchImgReferencePattern // Reference patterns to search for in CaptureArea
	BoundingCornerOffset Dimensions                        // Offset applied to the starting coordinates to define rectangular search area
	PostCaptureAlignment Alignment                         // Auto-adjustment applied to the result position
	PostCaptureOffset    Coordinates                       // Custom offset to be applied the result position after PostSearchAlignment
}

type ScreenPositionSearchResult struct {
	MatchedPattern     ScreenSearchImgReferencePattern // First successful matched pattern
	MatchedPatternSize Dimensions                      // First successful matched pattern rectangle size
	OnScreenPosition   Coordinates                     // Screen coordinates of MatchedPattern
	SearchAreaPosition Coordinates                     // Search area relative coordinates of MatchedPattern
	Position           Coordinates                     // Resulted coordinates after alignment and offset applied
}

type ScreenAreaVisualTracker struct {
	CachedScreenAreaCapture image.Image
	TrackedAreaInfo         ScreenArea
}

func NewScreenAreaVisualTracker(area ScreenArea) (ScreenAreaVisualTracker, error) {
	img, err := robotgo.CaptureImg(area.Offset.X, area.Offset.Y, area.Size.W, area.Size.H)
	if err != nil {
		return ScreenAreaVisualTracker{}, err
	}
	return ScreenAreaVisualTracker{
		CachedScreenAreaCapture: img,
		TrackedAreaInfo:         area,
	}, nil
}

func (o *ScreenAreaVisualTracker) Refresh() (float32, error) {
	img, err := robotgo.CaptureImg(o.TrackedAreaInfo.Offset.X, o.TrackedAreaInfo.Offset.Y, o.TrackedAreaInfo.Size.W, o.TrackedAreaInfo.Size.H)
	if err != nil {
		return 0.0, err
	}
	_, correlation, _, _ := gcv.FindImg(img, o.CachedScreenAreaCapture)
	return correlation, nil
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func EvalAlignment(position Coordinates, rect Dimensions, align Alignment) Coordinates {
	result := position

	switch align {
	case TopLeft:
		break // special case - do nothing
	case Top:
		result.X = result.X + rect.W/2
		break
	case TopRight:
		result.X = result.X + rect.W
		break
	case BottomLeft:
		result.Y = result.Y + rect.H
		break
	case Bottom:
		result.X = result.X + rect.W/2
		result.Y = result.Y + rect.H
		break
	case BottomRight:
		result.X = result.X + rect.W
		result.Y = result.Y + rect.H
		break
	case Left:
		result.Y = result.Y + rect.H/2
		break
	case Right:
		result.X = result.X + rect.W
		result.Y = result.Y + rect.H/2
		break
	case Center:
		result.X = result.X + rect.W/2
		result.Y = result.Y + rect.H/2
		break
	default:
		panic(errors.New("invalid alignment"))
	}

	return result
}

func ScreenPositionSearch(
	searchDefinitionStack []ScreenPositionSearchDefinition,
	startFrom Coordinates,
	config ScreenPositionSearchValidationConfig,
) (Coordinates, error) {
	searchResults, err := ScreenPositionSearchTrace(searchDefinitionStack, startFrom, config)
	if err != nil {
		return Coordinates{}, err
	}
	return searchResults[len(searchResults)-1].Position, nil
}

func ScreenPositionSearchTrace(
	searchDefinitionStack []ScreenPositionSearchDefinition,
	startFrom Coordinates,
	config ScreenPositionSearchValidationConfig,
) ([]ScreenPositionSearchResult, error) {
	retryAttempts := config.PatternMaxRetryAttempts

RetrySearch:
	searchResults, err := ScreenPositionSearchImpl(searchDefinitionStack, startFrom, ScreenPositionSearchResult{}, config)
	if err != nil {
		if retryAttempts--; retryAttempts >= 0 { // retry pattern after timeout (e.g. window is still loading, not fully rendered)
			time.Sleep(time.Duration(config.PatternRetryAttemptDelayMs) * time.Millisecond)
			goto RetrySearch
		}
		return []ScreenPositionSearchResult{}, err
	}
	return searchResults, nil
}

func IsTesseractAvailable() bool {
	cmd := exec.Command("tesseract", "--version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return false
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if errStr != "" {
		return false
	}

	verLineFormat := regexp.MustCompile(`tesseract\sv(\d+\.\d+\.\d+\.\d+)`)

	version := verLineFormat.FindStringSubmatch(outStr)
	if version == nil || len(version) != 2 {
		return false
	}

	return true
}

func TesseractValidation(imgFilePath, validationText string, config TesseractValidationConfig) error {
	if !config.Enabled {
		return nil
	}

	tCmd := exec.Command("tesseract", imgFilePath, "stdout", "-l", string(config.TesseractLang), "--psm", strconv.Itoa(int(config.TesseractPSM)))

	var stdout, stderr bytes.Buffer
	tCmd.Stdout = &stdout
	tCmd.Stderr = &stderr

	err := tCmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if err != nil {
		return errors.New(fmt.Sprintf("tesseract validation: %s -- [STDERR]: %s -- [STDOUT]: %s", err.Error(), errStr, outStr))
	}

	jaroWinklerDistance := 0.0
	if config.JaroWrinklerValidation.Enabled {
		jaroWinklerDistance = smetrics.JaroWinkler(outStr, validationText, config.JaroWrinklerValidation.BoostThreshold, config.JaroWrinklerValidation.PrefixSize)
	}
	fuzzyMatch := fuzzy.Match(validationText, outStr)

	if !fuzzyMatch || (config.JaroWrinklerValidation.Enabled && jaroWinklerDistance < config.JaroWrinklerValidation.DistanceThreshold) {
		return errors.New(fmt.Sprintf("tesseract validation: insufficient matching (fuzzy match: %t, jaro/winkler distance: %f)", fuzzyMatch, jaroWinklerDistance))
	}

	return nil
}

func ScreenPositionSearchImpl(
	searchDefinitionStack []ScreenPositionSearchDefinition,
	startFrom Coordinates,
	prevSearchResult ScreenPositionSearchResult,
	config ScreenPositionSearchValidationConfig,
) ([]ScreenPositionSearchResult, error) {
	defer func() {
		_ = os.Remove(".t.png")
		_ = os.Remove(".debug.png")
	}()

	if len(searchDefinitionStack) == 0 {
		return nil, errors.New("at least one search definition is required")
	}

	searchDefinition := searchDefinitionStack[0]

	boundingCornerOffset := searchDefinition.BoundingCornerOffset
	if boundingCornerOffset.W < 1 || boundingCornerOffset.H < 1 {
		boundingCornerOffset.W = prevSearchResult.MatchedPatternSize.W
		boundingCornerOffset.H = prevSearchResult.MatchedPatternSize.H
	}

	var searchResult ScreenPositionSearchResult

SearchPatternLoop:
	for index, searchPattern := range searchDefinition.Patterns {
		patternImages := make([]image.Image, len(searchPattern.ImageFilePaths))
		for idx, imageFilePath := range searchPattern.ImageFilePaths {
			loadedPatternImg, _, err := robotgo.DecodeImg(imageFilePath)
			if err != nil {
				return nil, err
			}
			patternImages[idx] = loadedPatternImg
		}

		screenCaptureArea := ScreenArea{
			Offset: Coordinates{
				X: Min(boundingCornerOffset.W+startFrom.X, startFrom.X),
				Y: Min(boundingCornerOffset.H+startFrom.Y, startFrom.Y),
			},
			Size: Dimensions{
				W: Abs(boundingCornerOffset.W), // must use abs(), if user specifies custom (negative) bounds
				H: Abs(boundingCornerOffset.H),
			},
		}

		safeguardBorder := CalculateSafeguardBorder(patternImages, screenCaptureArea)

		screenCaptureImg, err := robotgo.CaptureImg(
			screenCaptureArea.Offset.X-safeguardBorder.W,
			screenCaptureArea.Offset.Y-safeguardBorder.H,
			screenCaptureArea.Size.W+safeguardBorder.W,
			screenCaptureArea.Size.H+safeguardBorder.H)
		//_ = robotgo.Save(screenCaptureImg, ".debug.png") // debug only
		if err != nil {
			return nil, err
		}

		var unorderedPotentialMatches [][]gcv.Result
		const threshold = 0.75

		//for index, patternImage := range patternImages {
		//	_ = robotgo.Save(patternImage, fmt.Sprintf(".debug.png-%d.png", index)) // debug only
		//}
		if config.GrayScalePatternsEnabled {
			unorderedPotentialMatches = gcv.FindMultiAllImg(
				ConvertToGrayscaleAll(patternImages, config.GrayScalePatternsLuminanceLevels),
				ConvertToGrayscale(screenCaptureImg, config.GrayScalePatternsLuminanceLevels),
				threshold)
		} else {
			unorderedPotentialMatches = gcv.FindMultiAllImg(
				patternImages,
				screenCaptureImg,
				threshold)
		}

		sortedPotentialMatches, bestMatch, err := SortMatchesByRelevance(unorderedPotentialMatches)
		if err != nil {
			if index == len(searchDefinition.Patterns)-1 { // last pattern exhausted
				return nil, err
			}
			continue SearchPatternLoop // try next pattern
		}

		if searchPattern.ValidationText != "" { // if validation enabled, try to find the best matchToValidate if any
			validatedMatchFound := false
			for _, matchToValidate := range sortedPotentialMatches {

				patterImgScreenPosition := Coordinates{
					X: screenCaptureArea.Offset.X + matchToValidate.TopLeft.X,
					Y: screenCaptureArea.Offset.Y + matchToValidate.TopLeft.Y,
				}

				err = robotgo.SaveCapture(".t.png", patterImgScreenPosition.X, patterImgScreenPosition.Y, matchToValidate.ImgSize.W, matchToValidate.ImgSize.H)
				if err != nil {
					return nil, err
				}

				err = TesseractValidation(".t.png", searchPattern.ValidationText, config.TesseractValidation)
				if err != nil {
					continue
				}
				validatedMatchFound = true
				bestMatch = matchToValidate
				break
			}

			if !validatedMatchFound { // any of potential matches did not pass
				if index == len(searchDefinition.Patterns)-1 { // last pattern exhausted
					return nil, err
				}
				continue SearchPatternLoop // try next pattern
			}
		}

		alignedPosition := EvalAlignment(Coordinates{
			X: bestMatch.TopLeft.X,
			Y: bestMatch.TopLeft.Y,
		}, Dimensions{
			W: bestMatch.ImgSize.W,
			H: bestMatch.ImgSize.H,
		}, searchDefinition.PostCaptureAlignment)

		alignedPosition.X = alignedPosition.X + screenCaptureArea.Offset.X + searchDefinition.PostCaptureOffset.X
		alignedPosition.Y = alignedPosition.Y + screenCaptureArea.Offset.Y + searchDefinition.PostCaptureOffset.Y

		searchResult = ScreenPositionSearchResult{
			MatchedPattern: searchPattern,
			MatchedPatternSize: Dimensions{
				W: bestMatch.ImgSize.W,
				H: bestMatch.ImgSize.H,
			},
			OnScreenPosition: Coordinates{
				X: bestMatch.TopLeft.X + screenCaptureArea.Offset.X,
				Y: bestMatch.TopLeft.Y + screenCaptureArea.Offset.Y,
			},
			SearchAreaPosition: Coordinates{
				X: bestMatch.TopLeft.X,
				Y: bestMatch.TopLeft.Y,
			},
			Position: alignedPosition,
		}

		break SearchPatternLoop
	}

	// make sure recursion ends
	if len(searchDefinitionStack) > 1 {
		searchResults, err := ScreenPositionSearchImpl(searchDefinitionStack[1:], searchResult.Position, searchResult, config)
		if err != nil {
			if searchResults == nil {
				searchResults = []ScreenPositionSearchResult{}
			}
			return append(searchResults, searchResult), err
		}
		return searchResults, err
	}
	return []ScreenPositionSearchResult{searchResult}, nil // final successful search
}

func CalculateSafeguardBorder(patternImages []image.Image, screenCaptureArea ScreenArea) Dimensions {
	var patternImageMaxDimension Dimensions

	for _, patternImage := range patternImages {
		if patternImage.Bounds().Dx() > patternImageMaxDimension.W {
			patternImageMaxDimension.W = patternImage.Bounds().Dx()
		}
		if patternImage.Bounds().Dy() > patternImageMaxDimension.H {
			patternImageMaxDimension.H = patternImage.Bounds().Dy()
		}
	}

	var safeguardBorder Dimensions

	if screenCaptureArea.Size.W <= patternImageMaxDimension.W {
		safeguardBorder.W = (patternImageMaxDimension.W-screenCaptureArea.Size.W)/2 + 1
	}
	if screenCaptureArea.Size.H <= patternImageMaxDimension.H {
		safeguardBorder.H = (patternImageMaxDimension.H-screenCaptureArea.Size.H)/2 + 1
	}

	return safeguardBorder
}

func ConvertToGrayscale(img image.Image, luminanceLevels uint8) image.Image {
	imgBounds := img.Bounds()
	grayImg := image.NewGray(imgBounds)

	lumLevelDivider := luminanceLevels - 1
	if lumLevelDivider == 0 {
		// cannot divide by 0, luminanceLevels = 1 means one-color in image
		for x := 0; x < imgBounds.Max.X; x++ {
			for y := 0; y < imgBounds.Max.Y; y++ {
				grayImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
		return grayImg
	}

	lumLevelStep := int16(^uint8(0)) / int16(lumLevelDivider)
	lumLevelHalfStep := lumLevelStep / 2
	clampToUint8 := func(v int16) uint8 {
		if v > int16(^uint8(0)) {
			return ^uint8(0)
		} else if v < 0 {
			return 0
		}
		return uint8(v)
	}

	for x := 0; x < imgBounds.Max.X; x++ {
		for y := 0; y < imgBounds.Max.Y; y++ {
			rgba := img.At(x, y)
			gray := color.GrayModel.Convert(rgba).(color.Gray)

			lum16 := int16(gray.Y)
			lumLevelOffset := lum16 % lumLevelStep // 103
			if lumLevelOffset <= lumLevelHalfStep {
				grayImg.SetGray(x, y, color.Gray{Y: clampToUint8(lum16 - lumLevelOffset)})
			} else {
				grayImg.SetGray(x, y, color.Gray{Y: clampToUint8(lum16 + (lumLevelStep - lumLevelOffset))})
			}
		}
	}

	return grayImg
}

func ConvertToGrayscaleAll(images []image.Image, luminanceLevels uint8) []image.Image {
	var grayScaleImages []image.Image
	for _, img := range images {
		grayScaleImages = append(grayScaleImages, ConvertToGrayscale(img, luminanceLevels))
	}
	return grayScaleImages
}

func SortMatchesByRelevance(matches [][]gcv.Result) ([]gcv.Result, gcv.Result, error) {
	var flattenedMatches []gcv.Result
	for _, innerMatches := range matches {
		if innerMatches == nil || len(innerMatches) == 0 {
			continue
		}
		flattenedMatches = append(flattenedMatches, innerMatches...)
	}

	if len(flattenedMatches) == 0 {
		return []gcv.Result{}, gcv.Result{}, errors.New("no matches found")
	}

	sort.Slice(flattenedMatches, func(i, j int) bool {
		return slices.Max(flattenedMatches[i].MaxVal) > slices.Max(flattenedMatches[j].MaxVal)
	})

	return flattenedMatches, flattenedMatches[0], nil
}
