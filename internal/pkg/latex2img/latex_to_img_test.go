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

	s.converter = NewPlainLatexToImgConverter(tempDir, false, params.ImageDPI())
	s.pathToTestdata = "testdata/"
}

func (s *PlainLatexToImgConverterTestSuite) TestConvert_EqualToGeneratedPngs_WhenCorrectTexFiles() {
	for _, file := range params.CorrectTestdataFiles() {
        path := filepath.Join(s.pathToTestdata, file)

        s.Run(file, func() {
            img := s.correctlyConvertImgFromFile(path)

            s.imgEqualToImgFromFile(img, strings.TrimSuffix(path, ".tex") + ".png") 
        })
	}
}

func (s *PlainLatexToImgConverterTestSuite) TestConvert_ReturnError_WhenLatexCompilationError() {
    tests := []struct {
        file string
        expectedError *LatexCompilationError
    } {
        {
            "not_closed_math_brace_error.tex",
            &LatexCompilationError{
                Message: "Missing $ inserted.",
                Line: 4,
                Context: "<inserted text> \n                $\nl.4 \\end{document}",
            },
        },
        {
            "missing_package_error.tex",
            &LatexCompilationError{
                Message: "LaTeX Error: File `nonexistpackage.sty' not found.",
                Line: 0,
                Context: "",
            },
        },
        {
            "not_ended_document_error.tex",
            &LatexCompilationError{
                Message: "Emergency stop.",
                Line: 0,
                Context: "<*> document.tex",
            },
        },
        {
            "too_many_closed_brackets_error.tex",
            &LatexCompilationError{
                Message: "Extra }, or forgotten $.",
                Line: 3,
                Context: "l.3 VkTeX \\( \\frac{1}{b}}\n                          \\)",
            },
        },
        {
            "undefined_control_sequence.tex",
            &LatexCompilationError{
                Message: "Undefined control sequence.",
                Line: 3,
                Context: "l.3 \\dtae",
            },
        },
    }

    for _, test := range tests {
        path := filepath.Join(s.pathToTestdata, test.file)

        s.Run(test.file, func() {
            err := s.convertImgFromFileReturnLatexCompilationError(path)

            s.Require().Equal(err, test.expectedError)
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

func (s *PlainLatexToImgConverterTestSuite) correctlyConvertImgFromFile(path string) image.Image {
    content, err := os.ReadFile(path)
    s.Require().NoError(err)

    img, err := s.converter.Convert(context.Background(), content)
    s.Require().NoError(err)
    return img
}


func (s *PlainLatexToImgConverterTestSuite) convertImgFromFileReturnLatexCompilationError(path string) *LatexCompilationError {
    content, err := os.ReadFile(path)
    s.Require().NoError(err)
    
    var compilationError *LatexCompilationError
    _, err = s.converter.Convert(context.Background(), content)
    
    s.Require().ErrorAs(err, &compilationError)
    return compilationError
}
