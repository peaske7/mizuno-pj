package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	internalConvert "adelaida/internal/convert"
)

var convert = &cli.Command{
	Name:    "convert",
	Aliases: []string{"c"},
	Usage:   "Convert PDF files into images",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "Path where the PDF file is located",
			Required: true,
		}, &cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "The output directory where the images will be saved",
			Required: true,
		}, &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Image format",
			Value:   "jpg",
		},
	},
	Action: func(c *cli.Context) error {
		// determine whether inputPath is file or dir
		inputPath := c.String("input")
		outputPath := c.String("output")
		format := c.String("format")

		switch format {
		case "jpg", "jpeg":
			format = "jpg"
		case "png":
			format = "png"
		default:
			fmt.Printf("Invalid format: %v. Please specify 'jpg' or 'png'.\n", format)
			return errors.New("invalid format")
		}

		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Failed to open file at path: %v with error: %v\n", inputPath, err)
			return err
		}
		defer inputFile.Close()

		fileInfo, err := inputFile.Stat()
		if err != nil {
			fmt.Printf("Failed to get file info with error: %v\n", err)
			return err
		}

		pathInfo, err := os.Stat(outputPath)
		if err != nil {
			fmt.Printf("Failed to get path info with error: %v\n", err)
			return err
		}
		if !pathInfo.IsDir() {
			fmt.Printf("Output path: %v is not a directory. The output path must be a directory\n", outputPath)
			return errors.New("invalid output path")
		}

		if fileInfo.IsDir() {
			if c.String("name") != "" {
				fmt.Printf("The 'name' flag will be ignored, as the output directory is already specified in 'output' flag\n")
			}

			err := filepath.WalkDir(inputPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Printf("Failed to walk directory at path: %v with error: %v\n", path, err)
					return err
				}

				if d.IsDir() {
					fmt.Printf("Skipping directory at path: %v. Sorry, not handling recursive walking yet\n", path)
				} else {
					err = internalConvert.Images(outputPath, format, path)
					if err != nil {
						fmt.Printf("Failed to convert pdf to image with error: %v\n", err)
						return err
					}
				}

				return nil
			})
			if err != nil {
				fmt.Printf("Failed to walk directory at path: %v with error: %v\n", inputPath, err)
				return err
			}

		} else {
			if filepath.Ext(inputPath) != ".pdf" {
				fmt.Printf("Invalid file extension: %v. Please specify a PDF file.\n", filepath.Ext(inputPath))
				return errors.New("invalid file extension")
			}

			err = internalConvert.Images(outputPath, format, inputPath)
			if err != nil {
				fmt.Printf("Failed to convert pdf to image with error: %v\n", err)
				return err
			}
		}

		return nil
	},
}
