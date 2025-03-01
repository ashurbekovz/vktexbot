package params

import (
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/ashurbekovz/vktexbot/internal/utils/must"
	"github.com/shopspring/decimal"
)

func ImageDPI() decimal.Decimal {
    return decimal.NewFromInt(400)
}

type Testdata struct {
    Name string
    Text string
    Packages string
    Params template2img.ImageParams
}

var defaultFontSize = decimal.NewFromInt(10)

func GetTesdataConvertationParams() []Testdata {
    return []Testdata{
        {
            Name: "simple",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.MinImageSize(defaultFontSize, defaultFontSize),
                )),
        },
        {
            Name: "crop",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.Crop(),
                )),
        },
        {
            Name: "add_height",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.Crop(),
                template2img.AdditionalBorders(decimal.NewFromInt(30), decimal.Zero),
                )),
        },
        {
            Name: "add_width",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.Crop(),
                template2img.AdditionalBorders(decimal.Zero, decimal.NewFromInt(30)),
                )),
        },
        {
            Name: "min_image_height",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.MinImageSize(decimal.NewFromInt(100), decimal.Zero),
                )),
        },
        {
            Name: "min_image_width",
            Text: "some text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.MinImageSize(decimal.NewFromInt(100), decimal.Zero),
                )),
        },
        {
            Name: "min_image_width_to_font_size_ratio",
            Text: "text text text text text text text text text text text text text text text text text text text",
            Packages: "",
            Params: must.Get(template2img.NewImageParams(
                template2img.MinImageWidthToFontSizeRatio(defaultFontSize.Mul(decimal.NewFromInt(4))),
                )),
        },
    }
}
