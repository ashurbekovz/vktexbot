package latex2img

import (
	"context"
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

func (s *PlainLatexToImgConverterTestSuite) TestCorrectTexFiles() {
	for _, file := range params.CorrectTestdataFiles() {
		path := filepath.Join(s.pathToTestdata, file)

		content, err := os.ReadFile(path)
		s.Require().NoError(err)

		img, err := s.converter.Convert(context.Background(), content)
		s.Require().NoError(err)

		file, err := os.Open(strings.TrimSuffix(path, ".tex") + ".png")
		s.Require().NoError(err)

		expectedImg, err := png.Decode(file)
		file.Close()
		s.Require().NoError(err)

		s.Require().Equal(expectedImg, img)
	}
}
