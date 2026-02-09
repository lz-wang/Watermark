package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"watermark/pkg/watermark"
)

func main() {
	mode := flag.String("mode", "repeat", "watermark mode: repeat or position")
	input := flag.String("in", "", "input image path (required)")
	output := flag.String("out", "", "output image path (required)")
	text := flag.String("text", "", "watermark text (required)")

	colorHex := flag.String("color", "#4db6ac", "repeat: watermark color hex")
	space := flag.Int("space", 75, "repeat: spacing between tiles")
	angle := flag.Int("angle", 30, "repeat: rotation angle")
	opacity := flag.Float64("opacity", 0.5, "opacity 0..1")
	fontPath := flag.String("font", "", "font path (.ttf/.otf)")
	fontSize := flag.Int("font-size", 48, "repeat: font size")
	fontHeightCrop := flag.Float64("font-height-crop", 1.0, "repeat: font height crop factor")

	position := flag.String("position", "bottom-right", "position: bottom-right|bottom-left|top-right|top-left|center")
	marginRatio := flag.Float64("margin-ratio", 0.04, "position: margin ratio relative to width")
	jpgBG := flag.String("jpg-bg", "255,255,255", "jpeg background RGB, e.g. 255,255,255")

	flag.Parse()

	if err := validateRequired(*input, *output, *text); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	bg, err := parseRGB(*jpgBG)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid -jpg-bg:", err)
		os.Exit(2)
	}

	switch strings.ToLower(*mode) {
	case "repeat":
		if strings.TrimSpace(*fontPath) == "" {
			fmt.Fprintln(os.Stderr, "repeat mode requires -font to be set")
			os.Exit(2)
		}
		opts := &watermark.RepeatOptions{
			Color:          colorHex,
			Space:          space,
			Angle:          angle,
			Opacity:        opacity,
			FontPath:       *fontPath,
			FontSize:       fontSize,
			FontHeightCrop: fontHeightCrop,
		}
		_, err := watermark.AddRepeatWatermark(*input, *output, *text, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "position":
		opts := &watermark.PositionOptions{
			Opacity:       opacity,
			Position:      watermark.Position(strings.ToLower(*position)),
			FontPath:      *fontPath,
			MarginRatio:   marginRatio,
			JPGBackground: &bg,
		}
		_, err := watermark.AddPositionWatermark(*input, *output, *text, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "unsupported mode:", *mode)
		os.Exit(2)
	}
}

func validateRequired(input, output, text string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("missing -in")
	}
	if strings.TrimSpace(output) == "" {
		return errors.New("missing -out")
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("missing -text")
	}
	return nil
}

func parseRGB(raw string) (color.NRGBA, error) {
	parts := strings.Split(raw, ",")
	if len(parts) != 3 {
		return color.NRGBA{}, errors.New("expected format r,g,b")
	}
	vals := [3]uint8{}
	for i := 0; i < 3; i++ {
		p := strings.TrimSpace(parts[i])
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 || v > 255 {
			return color.NRGBA{}, fmt.Errorf("invalid channel: %q", p)
		}
		vals[i] = uint8(v)
	}
	return color.NRGBA{R: vals[0], G: vals[1], B: vals[2], A: 255}, nil
}
