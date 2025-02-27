package latex2img

import (
    "context"
    "fmt"
    "image"
    "os"
    "os/exec"
    "path/filepath"
)

type PlainLatexToImgConverter struct {
    workDir string
    clearWorkDir bool

    imageDPI string
}

func NewPlainLatexToImgConverter(
    workDir string,
    clearWorkDir bool,
    imageDPI string,
) *PlainLatexToImgConverter {
    return &PlainLatexToImgConverter{
        workDir: workDir,
        clearWorkDir: clearWorkDir,
        imageDPI: imageDPI,
    }
}

func (c *PlainLatexToImgConverter) Convert(ctx context.Context, content []byte) (image.Image, error) {
    tempDir, err := os.MkdirTemp(c.workDir, "latex-")
    if err != nil {
        return nil, fmt.Errorf("cant create tempdir: %w", err)
    }
    if c.clearWorkDir {
        defer os.RemoveAll(tempDir)
    }

    const latexResFile = "result"
    err = compileLatex(ctx, tempDir, latexResFile, content)
    if err != nil {
        return nil, err
    }

    image, err := dvi2img(ctx, tempDir, latexResFile, c.imageDPI)
    if err != nil {
        return nil, err
    }

    return image, nil
}

// compileLatex compile latex file with context to `tempDir`/`resFile`.dvi file. 
// Some spin-off files can some side files may be created in tempDir.
func compileLatex(
    ctx context.Context,
    tempDir, resFile string,
    content []byte,
) error {
    texFile := filepath.Join(tempDir, "document.tex")
    if err := os.WriteFile(texFile, content, 0644); err != nil {
        return fmt.Errorf("cant write .tex file: %w", err)
    }

    latexCmd := exec.CommandContext(
        ctx,
        // latexmk instead latex because it automatically determines the number of compilations require 
        "latexmk", 
        "-dvi",
        "-interaction=nonstopmode",
        "-jobname=" + resFile,
        "document.tex",
        )
    latexCmd.Dir = tempDir
    output, err := latexCmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("latex compilation error: %w, output: %s", err, output)
    }

    return nil
}

// dvi2img get `tempDir`/`inFile`.dvi file and convert it to rasterized image.Image.
// Some spin-off files can some side files may be created in tempDir.
func dvi2img(
    ctx context.Context,
    tempDir, inFile string,
    dpi string,
) (image.Image, error) {
    pngFile := filepath.Join(tempDir, "output.png")
    dvipngCmd := exec.CommandContext(
        ctx,
        "dvipng",
        "-D", dpi,
        "-T", "Tight",
        "-bg", "Transparent",
        "-o", "output.png",
        inFile + ".dvi",
        )
    dvipngCmd.Dir = tempDir
    output, err := dvipngCmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error during dvipng running: %w, output: %s", err, output)
    }

    pngData, err := os.Open(pngFile)
    if err != nil {
        return nil, fmt.Errorf("can't read png file: %w", err)
    }

    image, _, err := image.Decode(pngData)
    if err != nil {
        return nil, fmt.Errorf("can't decode png file %w:", err)
    }

    return image, nil
}
