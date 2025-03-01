package template2img_test

import (
	"context"
	"image"
	"image/png"
	"os"
	"path"
	"testing"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img/testdata_converter/params"
	"github.com/stretchr/testify/suite"
)

type TemplateToImgConverterTestSuite struct {
	suite.Suite

	pathToTestdata string
	baseConverter latex2img.LatexToImgConverter
}

func TestTemplateToImgConverterSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateToImgConverterTestSuite))
}

func (s *TemplateToImgConverterTestSuite) SetupSuite() {
	tempDir := "tmp"

	err := os.RemoveAll(tempDir)
	s.Require().NoError(err)

	err = os.Mkdir("./tmp", 0644)
	s.Require().NoError(err)

    s.baseConverter = latex2img.NewLatexToImgConverter(tempDir, false, params.ImageDPI())

	s.pathToTestdata = "testdata/"
}

func (s *TemplateToImgConverterTestSuite) TestConvert_EqualToGeneratedPngs_WhenOk() {
	for _, testdata := range params.GetTesdataConvertationParams() {

        s.Run(testdata.Name, func() {
            converter := template2img.NewLatexTemplateToImgConverter(&s.baseConverter, testdata.Packages)

            img, err := converter.Convert(context.Background(), testdata.Text, testdata.Params)
            s.Require().NoError(err)

            s.imgEqualToImgFromFile(img, path.Join(s.pathToTestdata, testdata.Name + ".png")) 
        })
	}
}

// Helpers

func (s *TemplateToImgConverterTestSuite) imgEqualToImgFromFile(img image.Image, expectedImgPath string) {
    file, err := os.Open(expectedImgPath)
    s.Require().NoError(err)
    defer file.Close()

    expectedImg, err := png.Decode(file)
    s.Require().NoError(err)

    s.Require().Equal(expectedImg, img)
}
