package resize

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/shopspring/decimal"
)

func createNewImage(src image.Image, rect image.Rectangle) image.Image {
	newImg := image.NewNRGBA(rect)

	draw.Draw(newImg, newImg.Bounds(), src, rect.Min, draw.Src)

	return newImg
}

func CropToBoundingBox(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if (minX > maxX) || (minY > maxY) {
		return image.Rectangle{}, fmt.Errorf("image fully transparent")
	}

	return createNewImage(img, image.Rect(minX, minY, maxX+1, maxY+1)), nil
}

func PixelToPt(pixels int64, dpi decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(int64(pixels)).Mul(decimal.NewFromInt(72)).Div(dpi)
}

func PtToPixel(pt decimal.Decimal, dpi decimal.Decimal) int64 {
	return pt.Mul(dpi).Div(decimal.NewFromInt(72)).Round(0).IntPart()
}

func EnlargeAndCenterImage(src image.Image, height, width int) image.Image {
	width = max(width, src.Bounds().Dx())
	height = max(height, src.Bounds().Dy())

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	transparent := color.NRGBA{0, 0, 0, 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{transparent}, image.Point{}, draw.Src)

	srcBounds := src.Bounds()
	dx := (width - srcBounds.Dx()) / 2
	dy := (height - srcBounds.Dy()) / 2

	draw.Draw(dst, image.Rect(dx, dy, dx+srcBounds.Dx(), dy+srcBounds.Dy()), src, srcBounds.Min, draw.Over)

	return dst
}
