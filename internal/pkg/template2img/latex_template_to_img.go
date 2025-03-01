package template2img

import (
	"context"
	"fmt"
	"image"

	"github.com/ashurbekovz/vktexbot/internal/utils/resize"
	"github.com/shopspring/decimal"
)

type ImageParams struct {
    crop bool
    fontSizePt decimal.Decimal
    minImageHeightPt decimal.Decimal
    minImageWidthPt decimal.Decimal
    minImageWidthToFontSizeRatio decimal.Decimal
    textWidthPt decimal.Decimal
    additionalImageHeightPt decimal.Decimal
    additionalImageWidthPt decimal.Decimal
}

type opt func(*ImageParams)

func Crop() opt {
    return func(imageParams *ImageParams) {
        imageParams.crop = true
    }
}

func FontSize(fontSizePt decimal.Decimal) opt {
    return func(imageParams *ImageParams) {
        imageParams.fontSizePt = fontSizePt
    }
}

func MinImageSize(heightPt, widthPt decimal.Decimal) opt {
    return func(imageParams *ImageParams) {
        imageParams.minImageHeightPt = heightPt
        imageParams.minImageWidthPt = widthPt
    }
}

func MinImageWidthToFontSizeRatio(ratio decimal.Decimal) opt {
    return func(imageParams *ImageParams) {
        imageParams.minImageWidthToFontSizeRatio = ratio
    }
}

func AdditionalBorders(heightPt, widthPt decimal.Decimal) opt {
    return func(imageParams *ImageParams) {
        imageParams.additionalImageHeightPt = heightPt
        imageParams.additionalImageWidthPt = widthPt
    }
}

func NewImageParams(opts ...opt) (ImageParams, error) {
    imageParams := ImageParams{
        crop: false,
        fontSizePt: decimal.RequireFromString("10"),
        minImageHeightPt: decimal.Zero,
        minImageWidthPt: decimal.Zero,
        minImageWidthToFontSizeRatio: decimal.Zero,
        textWidthPt: decimal.RequireFromString("236"),
        additionalImageHeightPt: decimal.Zero,
        additionalImageWidthPt: decimal.Zero,
    }

    for _, opt := range opts {
        opt(&imageParams)  
    }

    if imageParams.crop {
        if !imageParams.minImageWidthToFontSizeRatio.IsZero() {
            return ImageParams{}, fmt.Errorf("lineHeightToImageWidthRatio can be specified only if image not cropped")
        }

        if !imageParams.minImageHeightPt.IsZero() {
            return ImageParams{}, fmt.Errorf("minHeightPt can be specified only if image not cropped")
        }

        if !imageParams.minImageWidthPt.IsZero() {
            return ImageParams{}, fmt.Errorf("minWidthPt can be specified only if image not cropped")
        }
    }

    return imageParams, nil
}

type Converter interface {
    Convert(ctx context.Context, content []byte) (image.Image, error)
    GetDPI() decimal.Decimal
}

type LatexTemplateToImgConverter struct {
    baseConverter Converter
    
    packages string
}

func NewLatexTemplateToImgConverter(
    baseConverter Converter,
    packages string,
) LatexTemplateToImgConverter {
    return LatexTemplateToImgConverter{
        baseConverter: baseConverter,
        packages: packages,
    }
}

func (c *LatexTemplateToImgConverter) Convert(
    ctx context.Context,
    bodyText string,
    params ImageParams,
) (image.Image, error) {
    content := substituteToTemplate(bodyText, c.packages, params.textWidthPt, params.fontSizePt)

    img, err := c.baseConverter.Convert(ctx, []byte(content))
    if err != nil {
        return nil, fmt.Errorf("can't convert text to img: %w", err)
    }

    processedImg, err := processImage(img, c.baseConverter.GetDPI(), params)
    if err != nil {
        return nil, fmt.Errorf("can't process image: %w", err)
    }

    return processedImg, nil
}

func processImage(
	img image.Image,
	dpi decimal.Decimal,
	params ImageParams,
) (image.Image, error) {
    img, err := resize.CropToBoundingBox(img)
    if err != nil {
        return nil, fmt.Errorf("error while ctop img: %w", err)
    }
    
    finalHeightPt := resize.PixelToPt(int64(img.Bounds().Dy()), dpi)
    finalWidthPt := resize.PixelToPt(int64(img.Bounds().Dx()), dpi)
    if !params.crop {
        if !params.minImageWidthToFontSizeRatio.IsZero() {
            finalWidthPt = decimal.Max(finalWidthPt, params.fontSizePt.Mul(params.minImageWidthToFontSizeRatio))
        }

        if !params.minImageHeightPt.IsZero() {
            finalHeightPt = decimal.Max(finalHeightPt, params.minImageHeightPt)
        }

        if !params.minImageWidthPt.IsZero() {
            finalWidthPt = decimal.Max(finalWidthPt, params.minImageWidthPt)
        }
    }

    finalWidthPt = finalWidthPt.Add(params.additionalImageWidthPt)
    finalHeightPt = finalHeightPt.Add(params.additionalImageHeightPt)

    img = resize.EnlargeAndCenterImage(
    	img,
    	int(resize.PtToPixel(finalHeightPt, dpi)),
    	int(resize.PtToPixel(finalWidthPt, dpi)),
    )

    return img, nil
}

func substituteToTemplate(
	bodyText, packages string,
	textWidthPt, fontSizePt decimal.Decimal,
) string {

    lineSpacing := fontSizePt.Mul(decimal.NewFromFloat(1.2))

    content := `
\documentclass[preview]{standalone}
\setlength{\textwidth}{` + textWidthPt.String() + `pt}
\fontsize{` + fontSizePt.String() + `pt}{` + lineSpacing.String() + `pt}\selectfont
` + packages + `
\begin{document}
` + bodyText + `
\end{document}`

    return content
}

