package latex2img

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shopspring/decimal"
)

type LatexToImgConverter struct {
	workDir      string
	clearWorkDir bool

	imageDPI decimal.Decimal
}

func NewLatexToImgConverter(
	workDir string,
	clearWorkDir bool,
	imageDPI decimal.Decimal,
) LatexToImgConverter {
	return LatexToImgConverter{
		workDir:      workDir,
		clearWorkDir: clearWorkDir,
		imageDPI:     imageDPI,
	}
}

func (c *LatexToImgConverter) GetDPI() decimal.Decimal {
	return c.imageDPI
}

func (c *LatexToImgConverter) Convert(ctx context.Context, content []byte) (image.Image, error) {
	tempDir, err := os.MkdirTemp(c.workDir, "latex-")
	if err != nil {
		return nil, fmt.Errorf("cant create tempdir: %w", err)
	}
	if c.clearWorkDir {
		defer os.RemoveAll(tempDir)
	}

	const latexResFile = "result"
	err = CompileLatex(ctx, tempDir, latexResFile, content)
	if err != nil {
		return nil, err
	}

	image, err := dvi2img(ctx, tempDir, latexResFile, c.imageDPI)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// CompileLatex compile latex file with context to `tempDir`/`resFile`.dvi file.
// Some spin-off files can some side files may be created in tempDir.
func CompileLatex(
	ctx context.Context,
	tempDir, resFile string,
	content []byte,
) error {
	texFile := filepath.Join(tempDir, "document.tex")
	if err := os.WriteFile(texFile, content, 0644); err != nil {
		return fmt.Errorf("cant write .tex file: %w", err)
	}

	tempDirFull, err := filepath.Abs(tempDir)
	if err != nil {
		return fmt.Errorf("can't get absolute path for tempDir: %w", err)
	}

	latexCmd := exec.CommandContext(
		ctx,
		"bwrap",
		"--proc", "/proc",
		"--dev", "/dev",
		"--tmpfs", "/tmp",
		"--ro-bind", "/usr", "/usr",
		"--ro-bind", "/lib", "/lib",
		"--ro-bind", "/bin", "/bin",
		"--ro-bind", "/etc", "/etc",
		"--ro-bind", "/var/lib/texmf", "/var/lib/texmf",
		"--bind", tempDirFull, tempDirFull,
		"/usr/bin/latexmk",
		"-dvi",
		"-interaction=nonstopmode",
		"-jobname="+resFile,
		"document.tex",
	)
	latexCmd.Dir = tempDir
	output, err := latexCmd.CombinedOutput()
	if err != nil {
		log.Println(latexCmd.Args)
		log.Println(string(output))
		return fmt.Errorf("latex compilation error: %w", parseLatexError(output))
	}

	return nil
}

// dvi2img get `tempDir`/`inFile`.dvi file and convert it to rasterized image.Image.
// Some spin-off files can some side files may be created in tempDir.
func dvi2img(
	ctx context.Context,
	tempDir, inFile string,
	dpi decimal.Decimal,
) (image.Image, error) {
	dvipngCmd := exec.CommandContext(
		ctx,
		"dvipng",
		"-D", dpi.String(),
		"-T", "Tight",
		"-bg", "Transparent",
		"-o", "output.png",
		inFile+".dvi",
	)
	dvipngCmd.Dir = tempDir
	output, err := dvipngCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error during dvipng running:\n %w\n, output:\n %s\n", err, output)
	}

	pngFilename := filepath.Join(tempDir, "output.png")
	pngFile, err := os.Open(pngFilename)
	if err != nil {
		return nil, fmt.Errorf("can't read png file: %w", err)
	}
	defer pngFile.Close()

	image, err := png.Decode(pngFile)
	if err != nil {
		return nil, fmt.Errorf("can't decode png file '%s': %w", pngFilename, err)
	}

	return image, nil
}
