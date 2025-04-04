package latex2img_test

import (
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img/testdata_converter/params"
	"github.com/stretchr/testify/suite"
)

type LatexToImgConverterTestSuite struct {
	suite.Suite

	pathToTestdata string
	converter      latex2img.LatexToImgConverter
}

func TestLatexToImgConverterSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(LatexToImgConverterTestSuite))
}

func (s *LatexToImgConverterTestSuite) SetupSuite() {
	tempDir := "tmp"

	err := os.RemoveAll(tempDir)
	s.Require().NoError(err)

	err = os.Mkdir("./tmp", 0644)
	s.Require().NoError(err)

	s.converter = latex2img.NewLatexToImgConverter(tempDir, false, params.ImageDPI())
	s.pathToTestdata = "testdata/"
}

func (s *LatexToImgConverterTestSuite) TestConvert_EqualToGeneratedPngs_WhenCorrectTexFiles() {
	for _, file := range params.CorrectTestdataFiles() {
		path := filepath.Join(s.pathToTestdata, file)

		s.Run(file, func() {
			img := s.correctlyConvertImgFromFile(path)

			s.imgEqualToImgFromFile(img, strings.TrimSuffix(path, ".tex")+".png")
		})
	}
}

func (s *LatexToImgConverterTestSuite) TestConvert_ReturnError_WhenLatexCompilationError() {
	tests := []struct {
		file          string
		expectedError *latex2img.LatexCompilationError
	}{
		{
			"not_closed_math_brace_error.tex",
			&latex2img.LatexCompilationError{
				Message: "Missing $ inserted.",
				Line:    4,
				Context: "<inserted text> \n                $\nl.4 \\end{document}",
			},
		},
		{
			"missing_package_error.tex",
			&latex2img.LatexCompilationError{
				Message: "LaTeX Error: File `nonexistpackage.sty' not found.",
				Line:    0,
				Context: "",
			},
		},
		{
			"not_ended_document_error.tex",
			&latex2img.LatexCompilationError{
				Message: "Emergency stop.",
				Line:    0,
				Context: "<*> document.tex",
			},
		},
		{
			"too_many_closed_brackets_error.tex",
			&latex2img.LatexCompilationError{
				Message: "Extra }, or forgotten $.",
				Line:    3,
				Context: "l.3 VkTeX \\( \\frac{1}{b}}\n                          \\)",
			},
		},
		{
			"undefined_control_sequence.tex",
			&latex2img.LatexCompilationError{
				Message: "Undefined control sequence.",
				Line:    3,
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

func (s *LatexToImgConverterTestSuite) imgEqualToImgFromFile(img image.Image, expectedImgPath string) {
	file, err := os.Open(expectedImgPath)
	s.Require().NoError(err)
	defer file.Close()

	expectedImg, err := png.Decode(file)
	s.Require().NoError(err)

	s.Require().Equal(expectedImg, img)
}

func (s *LatexToImgConverterTestSuite) correctlyConvertImgFromFile(path string) image.Image {
	content, err := os.ReadFile(path)
	s.Require().NoError(err)

	img, err := s.converter.Convert(context.Background(), content)
	s.Require().NoError(err)
	return img
}

func (s *LatexToImgConverterTestSuite) convertImgFromFileReturnLatexCompilationError(path string) *latex2img.LatexCompilationError {
	content, err := os.ReadFile(path)
	s.Require().NoError(err)

	var compilationError *latex2img.LatexCompilationError
	_, err = s.converter.Convert(context.Background(), content)

	s.Require().ErrorAs(err, &compilationError)
	return compilationError
}

type CompileLatexSuite struct {
	suite.Suite

	pathToTestdata string
}

func CompileLatexTestSuite(t *testing.T) {
	suite.Run(t, new(CompileLatexSuite))
}

type CompileLatex struct {
	suite.Suite

	tempDir string
}

func TestCompileLatex(t *testing.T) {
	suite.Run(t, new(CompileLatex))
}

func (s *CompileLatex) SetupTest() {
	tempDir, err := os.MkdirTemp("", "test-latex-permissions-")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *CompileLatex) TearDownTest() {
	err := os.RemoveAll(s.tempDir)
	s.Require().NoError(err)
}

func (s *CompileLatex) Test_ReturnOK_WhenDirectoryWithoutAccess() {
	restrictedFile := filepath.Join(s.tempDir, "..", "restricted_file.txt")
	err := os.WriteFile(restrictedFile, []byte("test\n"), 0644)
	s.Require().NoError(err)
	latexContent := []byte("\\documentclass{article}\n" +
		"\\usepackage{amsmath,amsthm,amssymb,amsfonts,mathtools,mathtext,physics,tikz,bigints}\n" +
		"\\usepackage[T1,T2A]{fontenc}\n" +
		"\\usepackage[utf8]{inputenc}\n" +
		"\\usepackage[english,russian]{babel}\n" +
		"\\usepackage{listings}\n" +
		"\\usepackage{xcolor}\n" +
		"\\begin{document}\n" +
		"lol\n" +
		"\\end{document}")

	err = latex2img.CompileLatex(context.Background(), s.tempDir, "test", latexContent)

	s.Require().NoError(err)
}

func (s *CompileLatex) Test_ReturnError_WhenDirectoryWithoutAccess() {
	restrictedFile := filepath.Join(s.tempDir, "..", "restricted_file.txt")
	err := os.WriteFile(restrictedFile, []byte("test\n"), 0644)
	s.Require().NoError(err)
	latexContent := []byte("\\documentclass{article}\n" +
		"\\begin{document}\n" +
		"\\input{" + restrictedFile + "}\n" +
		"\\end{document}")

	err = latex2img.CompileLatex(context.Background(), s.tempDir, "test", latexContent)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "File `/tmp/restricted_file.txt' not found.")
}
