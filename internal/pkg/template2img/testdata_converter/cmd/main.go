package main

import (
	"context"
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img/testdata_converter/params"
)


func main() {
    pathToTestdata := flag.String("path", "testdata", "path to testdata")
    flag.Parse()
    
    createTmpIfNotExists()
    baseConverter := latex2img.NewLatexToImgConverter("tmp", true, params.ImageDPI())
    
    for _, testdata := range params.GetTesdataConvertationParams() {
        fmt.Printf("Processing testdata: %s", testdata.Name)

        path := filepath.Join(*pathToTestdata, testdata.Name)

        templateConverter := template2img.NewLatexTemplateToImgConverter(&baseConverter, testdata.Packages)

        img, err := templateConverter.Convert(context.Background(), testdata.Text, testdata.Params)
        if err != nil {
            fmt.Printf("Error converting file %s: %v\n", path, err)
            os.Exit(1)
        }

        outputFile := path + ".png"

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
    fmt.Printf("Create tmp directory\n")

    err := os.Mkdir("./tmp", 0644)
    if err == nil {
        return
    }
    
    if os.IsExist(err) {
        fmt.Printf("tmp directory already exists\n")
        return
    }

    fmt.Printf("Can't create tmp dir: %v\n", err)
    os.Exit(1)
}
