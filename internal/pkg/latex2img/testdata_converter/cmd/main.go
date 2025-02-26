package main

import (
	"context"
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img/testdata_converter/params"
)

func main() {
    pathToTestdata := flag.String("path", "testdata", "path to testdata")
    flag.Parse()

    converter := latex2img.NewPlainLatexToImgConverter("tmp", params.ImageDPI())
    
    for _, file := range params.CorrectTestdataFiles() {
        path := filepath.Join(*pathToTestdata, file)
        fmt.Printf("Processing file: %s\n", path)

        content, err := os.ReadFile(path)
        if err != nil {
            fmt.Printf("Error reading file %s: %v\n", path, err)
            os.Exit(1)
        }

        img, err := converter.Convert(context.Background(), content)
        if err != nil {
            fmt.Printf("Error converting file %s: %v\n", path, err)
            os.Exit(1)
        }

        outputFile := strings.TrimSuffix(path, ".tex") + ".png"

        out, err := os.Create(outputFile)
        if err != nil {
            fmt.Printf("Error creating PNG file %s: %v\n", outputFile, err)
            os.Exit(1)
        }
        defer out.Close()

        err = png.Encode(out, img)
        if err != nil {
            fmt.Printf("Error encoding PNG file %s: %v\n", outputFile, err)
            os.Exit(1)
        }

        fmt.Printf("Successfully converted %s to %s\n", path, outputFile)
    }
}

func createTmpIfNotExists() {
    fmt.Printf("Create tmp directory")

    err := os.Mkdir("./tmp", 0644)
    if err == nil {
        return
    }
    
    if os.IsNotExist(err) {
        fmt.Printf("tmp directory already exists")
        return
    }

    fmt.Printf("Can't create tmp dir: %v", err)
    os.Exit(1)
}
