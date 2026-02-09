# Go Watermark（水印工具）

本项目提供 Go 语言的图片水印工具，功能包括：

- 重复平铺文字水印（可调整间距、角度、透明度、字号、颜色）
- 单点位置水印（根据背景亮度自动选用黑/白字，并绘制描边）
- 保存 JPEG 时自动进行背景合成，避免透明通道丢失

## 构建

```bash
go build ./cmd/watermark
```

## CLI 用法

重复平铺水印（需要指定字体路径）：

```bash
./watermark -mode repeat \
  -in input.jpg \
  -out out.jpg \
  -text "CONFIDENTIAL" \
  -font /path/to/font.ttf
```

单点位置水印：

```bash
./watermark -mode position \
  -in input.jpg \
  -out out.jpg \
  -text "CONFIDENTIAL"
```

## 作为库使用

```go
package main

import (
	"image/color"

	"watermark/pkg/watermark"
)

func main() {
	colorHex := "#4db6ac"
	space := 75
	angle := 30
	opacity := 0.5
	fontSize := 48
	fontHeightCrop := 1.0

	_, _ = watermark.AddRepeatWatermark(
		"input.jpg",
		"out.jpg",
		"CONFIDENTIAL",
		&watermark.RepeatOptions{
			Color:          &colorHex,
			Space:          &space,
			Angle:          &angle,
			Opacity:        &opacity,
			FontPath:       "/path/to/font.ttf",
			FontSize:       &fontSize,
			FontHeightCrop: &fontHeightCrop,
		},
	)

	margin := 0.04
	bg := colorNRGBA(255, 255, 255)
	_, _ = watermark.AddPositionWatermark(
		"input.jpg",
		"out.jpg",
		"CONFIDENTIAL",
		&watermark.PositionOptions{
			Opacity:       &opacity,
			Position:      watermark.BottomRight,
			MarginRatio:   &margin,
			JPGBackground: &bg,
		},
	)
}

func colorNRGBA(r, g, b uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}
```

## 说明

- `repeat` 模式要求提供字体路径（与 Python 版本一致）。
- `position` 模式优先使用传入字体，若为空或加载失败，会尝试常见的 Arial 路径，最后回退到 Go 内置字体。
