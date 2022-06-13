package pdf

import (
	"github.com/signintech/gopdf"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func propBoxWidth(t *Text, pageWidth float64) {
	var (
		boxWidth float64
		props    = t.Props
	)

	if box, ok := props["box-width"]; ok {
		boxWidth, ok = box.(float64)
		if !ok {
			return
		}
	} else {
		boxWidth = pageWidth
	}

	props["box-width"] = boxWidth
}

func propAlign(t *Text) {
	var (
		boxWidth    float64
		stringWidth = float64(utf8.RuneCountInString(t.Text))
		props       = t.Props
	)

	// 检查盒子是否有宽度
	if box, ok := props["box-width"]; ok {
		boxWidth, ok = box.(float64)
		if !ok {
			return
		}
	} else {
		return
	}

	if a, ok := props["align"]; ok {
		align := a.(string)

		switch align {
		case "left":
			// 无需修改
			return
		case "center":
			// 居中
			t.SetX(t.OffsetX + (boxWidth-stringWidth)/2)
		case "right":
			// 右对齐
			t.SetX(t.OffsetX + (boxWidth - stringWidth))
			return
		}
	}
}

func propHeight(i *Image) {
	var (
		rectHeight float64
		props      = i.Props
	)

	if height, ok := props["height"]; ok {
		rectHeight, ok = height.(float64)
		if !ok {
			i.Rect = nil
		}

		if i.Rect != nil {
			i.Rect.H = rectHeight
		} else {
			i.Rect = &gopdf.Rect{H: rectHeight}
		}
	} else {
		i.Rect = nil
	}
}

func propWidth(i *Image) {
	var (
		rectWidth float64
		props     = i.Props
	)

	if width, ok := props["width"]; ok {
		rectWidth, ok = width.(float64)
		if !ok {
			i.Rect = nil
		}

		if i.Rect != nil {
			i.Rect.W = rectWidth
		} else {
			i.Rect = &gopdf.Rect{W: rectWidth}
		}
	} else {
		i.Rect = nil
	}
}

func propFontFamily(t *Text) string {
	var (
		props = t.Props
	)

	if family, ok := props["family"]; ok {
		return family.(string)
	} else {
		return _config.Default.Font
	}
}

func propFontStyle(t *Text) string {
	return ""
}

func propFontSize(t *Text) uint {
	var (
		props = t.Props
		size  uint
	)

	if s, ok := props["size"]; ok {
		size = uint(s.(float64))

	} else {
		size = _config.Default.FontSize
	}

	return size
}

func propFontColor(t *Text) (uint8, uint8, uint8) {
	var (
		props = t.Props
		color string
	)

	if c, ok := props["color"]; ok {
		color, ok = c.(string)
		if !ok {
			color = "000000"
		}
	} else {
		color = "000000"
	}

	reg := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	if reg == nil {
		return 0, 0, 0
	}

	if !reg.MatchString(color) {
		return 0, 0, 0
	}

	c := strings.Trim(color, "#")
	c = strings.ToLower(c)
	r, err := strconv.ParseInt(c[0:2], 16, 0)
	g, err := strconv.ParseInt(c[2:4], 16, 0)
	b, err := strconv.ParseInt(c[4:6], 16, 0)

	if err != nil {
		return 0, 0, 0
	}

	return uint8(r), uint8(g), uint8(b)
}
