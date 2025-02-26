package latex2img

import (
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img/testdata_converter/params"
	"github.com/stretchr/testify/suite"
)

type PlainLatexToImgConverterTestSuite struct {
	suite.Suite

	pathToTestdata string
	converter      *PlainLatexToImgConverter
}

func TestPlainLatexToImgConverterSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(PlainLatexToImgConverterTestSuite))
}

func (s *PlainLatexToImgConverterTestSuite) SetupSuite() {
	tempDir := "tmp"

	err := os.RemoveAll(tempDir)
	s.Require().NoError(err)

	err = os.Mkdir("./tmp", 0644)
	s.Require().NoError(err)

	s.converter = NewPlainLatexToImgConverter(tempDir, params.ImageDPI())
	s.pathToTestdata = "testdata/"
}

func (s *PlainLatexToImgConverterTestSuite) TestConvert_EqualToGeneratedPngs_WhenCorrectTexFiles() {
	for _, file := range params.CorrectTestdataFiles() {
        s.Run(s.pathToTestdata, func() {
            path := filepath.Join(s.pathToTestdata, file)

            img := s.convertImgFromFile(path)

            s.imgEqualToImgFromFile(img, strings.TrimSuffix(path, ".tex") + ".png") 
        })
	}
}

// Helpers

func (s *PlainLatexToImgConverterTestSuite) imgEqualToImgFromFile(img image.Image, expectedImgPath string) {
    file, err := os.Open(expectedImgPath)
    s.Require().NoError(err)
    defer file.Close()

    expectedImg, err := png.Decode(file)
    s.Require().NoError(err)

    s.Require().Equal(expectedImg, img)
}

func (s *PlainLatexToImgConverterTestSuite) convertImgFromFile(path string) image.Image {
    content, err := os.ReadFile(path)
    s.Require().NoError(err)

    img, err := s.converter.Convert(context.Background(), content)
    s.Require().NoError(err)
    return img
}
