package main

import (
	"context"
	"flag"
	"fmt"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
)

func main() {
    pathToTestdata := flag.String("path", "testdata", "path to testdata")
    flag.Parse()

    converter := latex2img.NewPlainLatexToImgConverter("", "400")

	err := filepath.Walk(*pathToTestdata, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tex") {
			fmt.Printf("Processing file: %s\n", path)

			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}

			img, err := converter.Convert(context.Background(), content)
			if err != nil {
				fmt.Printf("Error converting file %s: %v\n", path, err)
				return nil
			}

			outputFile := strings.TrimSuffix(path, ".tex") + ".png"

			out, err := os.Create(outputFile)
			if err != nil {
				fmt.Printf("Error creating PNG file %s: %v\n", outputFile, err)
				return nil
			}
			defer out.Close()

			err = png.Encode(out, img)
			if err != nil {
				fmt.Printf("Error encoding PNG file %s: %v\n", outputFile, err)
				return nil
			}

			fmt.Printf("Successfully converted %s to %s\n", path, outputFile)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error during processing: %v\n", err)
	}
}
