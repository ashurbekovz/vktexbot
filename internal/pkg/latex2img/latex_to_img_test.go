package latex2img

// import (
// 	"os"
// 	"testing"
// 
// 	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img/testdata_converter/params"
// 	"github.com/stretchr/testify/suite"
// )
// 
// type PlainLatexToImgConverterTestSuite struct {
//     suite.Suite
// 
//     converter *PlainLatexToImgConverter
// }
// 
// func TestPlainLatexToImgConverterSuiteTestSuite(t *testing.T) {
//     suite.Run(t, new(PlainLatexToImgConverterTestSuite))
// }
// 
// func (s *PlainLatexToImgConverterTestSuite) SetupSuite() {
//     tempDir := "tmp"
//     os.RemoveAll(tempDir)
// 
//     s.converter = NewPlainLatexToImgConverter(tempDir, params.ImageDPI())
// }
// 
// func (s *PlainLatexToImgConverterTestSuite) TestCorrectTexFiles()
