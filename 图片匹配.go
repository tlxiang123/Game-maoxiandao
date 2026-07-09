package main

import (
	"bytes"
	stdimage "image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"

	agimages "github.com/Dasongzi1366/AutoGo/images"
)

func (z *Zg) findPicPure(feature *Pic, template *[]byte) (int, int) {
	if feature == nil || template == nil || len(*template) == 0 {
		return -1, -1
	}

	x1, y1, x2, y2 := z.offsetRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	screen := agimages.CaptureScreen(x1, y1, x2, y2, z.displayID())
	if screen == nil {
		return -1, -1
	}

	tplImg, _, err := stdimage.Decode(bytes.NewReader(*template))
	if err != nil {
		return -1, -1
	}
	tpl := toNRGBA(tplImg)
	if tpl == nil {
		return -1, -1
	}

	return findTemplate(screen, tpl, x1, y1, defaultSim(feature.Sim), feature.Gray, feature.Transparent)
}

func findTemplate(screen, tpl *stdimage.NRGBA, baseX, baseY int, sim float32, gray, transparent bool) (int, int) {
	sw := screen.Bounds().Dx()
	sh := screen.Bounds().Dy()
	tw := tpl.Bounds().Dx()
	th := tpl.Bounds().Dy()
	if sw <= 0 || sh <= 0 || tw <= 0 || th <= 0 || tw > sw || th > sh {
		return -1, -1
	}

	transparentColor := tpl.NRGBAAt(tpl.Bounds().Min.X, tpl.Bounds().Min.Y)
	for y := 0; y <= sh-th; y++ {
		for x := 0; x <= sw-tw; x++ {
			if templateMatchesAt(screen, tpl, x, y, sim, gray, transparent, transparentColor) {
				return baseX + x, baseY + y
			}
		}
	}
	return -1, -1
}

func templateMatchesAt(screen, tpl *stdimage.NRGBA, startX, startY int, sim float32, gray, transparent bool, transparentColor color.NRGBA) bool {
	tb := tpl.Bounds()
	for ty := 0; ty < tb.Dy(); ty++ {
		for tx := 0; tx < tb.Dx(); tx++ {
			tp := tpl.NRGBAAt(tb.Min.X+tx, tb.Min.Y+ty)
			if transparent && isTransparentTemplatePixel(tp, transparentColor) {
				continue
			}
			sp := screen.NRGBAAt(screen.Bounds().Min.X+startX+tx, screen.Bounds().Min.Y+startY+ty)
			if !pixelMatch(sp, tp, sim, gray) {
				return false
			}
		}
	}
	return true
}

func isTransparentTemplatePixel(pixel, transparentColor color.NRGBA) bool {
	if pixel.A < 16 {
		return true
	}
	return pixel.R == transparentColor.R && pixel.G == transparentColor.G && pixel.B == transparentColor.B
}

func pixelMatch(a, b color.NRGBA, sim float32, gray bool) bool {
	if gray {
		ga := luminance(a)
		gb := luminance(b)
		return 1-float32(absInt(ga-gb))/255 >= sim
	}
	diff := absInt(int(a.R)-int(b.R)) + absInt(int(a.G)-int(b.G)) + absInt(int(a.B)-int(b.B))
	return 1-float32(diff)/765 >= sim
}

func luminance(c color.NRGBA) int {
	return int(math.Round(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)))
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func toNRGBA(src stdimage.Image) *stdimage.NRGBA {
	if src == nil {
		return nil
	}
	if img, ok := src.(*stdimage.NRGBA); ok {
		return img
	}
	bounds := src.Bounds()
	dst := stdimage.NewNRGBA(stdimage.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			dst.Set(x, y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}
