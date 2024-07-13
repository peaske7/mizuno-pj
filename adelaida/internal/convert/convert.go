package convert

import (
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

func Images(outputPath, format, path string) error {
	if filepath.Ext(path) != ".pdf" {
		fmt.Printf("Invalid file extension: %v. Please specify a PDF file.\n", filepath.Ext(path))
		return errors.New("invalid file extension")
	}

	doc, err := fitz.New(path)
	if err != nil {
		fmt.Printf("Failed to open PDF file at path: %v with error: %v\n", path, err)
		return err
	}

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(outputPath+"/", fmt.Sprintf("image-%05d.jpg", n)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		f.Close()

	}

	return nil
}
