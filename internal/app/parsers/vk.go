package parsers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/shopspring/decimal"
)

type Vk struct {
}

func NewVk() *Vk {
	return &Vk{}
}

type ParseRes struct {
	Message     string
	ImageParams template2img.ImageParams
	Mention     bool
}

var looksLikeMention *regexp.Regexp = regexp.MustCompile(`^\[.*\|@.*\]$`)

func (e *Vk) Parse(rawMessage string) (ParseRes, error) {
	// TODO(ashurbekovz): all constants should be taken from receiver
	// TODO(ashurbekovz): separate to parser + imageParams converter
	parseRes := ParseRes{}
	chunks := strings.Split(rawMessage, " ")

	if len(chunks) > 0 && looksLikeMention.MatchString(chunks[0]) {
		parseRes.Mention = true
		chunks = chunks[1:]
	}

	// TODO(ashurbekovz): make flags processing more strict (e.g. disallow
	// multiple flags with string size)
	fontSizePt := decimal.NewFromInt(10)
	opts := []template2img.Opt{template2img.FontSize(fontSizePt)}
	i := 0
	hasCrop := false
	hasHa, hasWa := false, false
loop:
	for ; i < len(chunks); i++ {
		c := chunks[i]
		switch c {
		case "-c":
			hasCrop = true
			opts = append(opts, template2img.Crop())
		case "-l":
			opts = append(opts, template2img.TextWidth(decimal.NewFromInt(300)))
		case "-l2":
			opts = append(opts, template2img.TextWidth(decimal.NewFromInt(400)))
		case "-ha":
			hasHa = true
		case "-wa":
			hasWa = true
		default:
			break loop
		}
	}
	if !hasCrop {
		opts = append(opts, template2img.MinImageSize(fontSizePt, fontSizePt))
		opts = append(opts, template2img.MinImageWidthToFontSizeRatio(decimal.NewFromInt(22)))
		opts = append(opts, template2img.MinImageHeightToFontSizeRatio(decimal.NewFromInt(7)))

		wa := decimal.NewFromInt(5)
		if hasWa {
			wa = decimal.NewFromInt(15)
		}

		ha := decimal.NewFromInt(5)
		if hasHa {
			ha = decimal.NewFromInt(15)
		}

		opts = append(opts, template2img.AdditionalBorders(ha, wa))
	}

	parseRes.Message = strings.Join(chunks[i:], " ")

	var err error
	parseRes.ImageParams, err = template2img.NewImageParams(opts...)
	if err != nil {
		return ParseRes{}, fmt.Errorf("failed to create image params: %w", err)
	}

	return parseRes, nil
}
