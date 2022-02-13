package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"image/png"
	_ "image/png"
)

func main() {
	fileName := flag.String("f", "", "Path for the file you want to open")
	res := flag.String("res", "", "Image height and width separated by a comma, e.g. 20,10 for a 20x10 image. Note that no spaces can be used.")

	flag.Parse()

	var m image.Image
	var err error

	isExistingFile := fileExists(*fileName)

	if *fileName != "" {
		if isExistingFile {
			m, err = loadImage(*fileName)
			if err != nil {
				log.Fatal(err)
			}
		} else if *res != "" && !isExistingFile {
			m, err = createImage(*res)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := NewCmdPxl(*fileName, m, saveImage).Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("need to set either existing filename or resolution and new filename")
	}
}

func loadImage(fileName string) (image.Image, error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func createImage(res string) (image.Image, error) {
	resArr := strings.Split(res, ",")
	if len(resArr) != 2 {
		return nil, fmt.Errorf("invalid resolution %s", res)
	}
	h, err := strconv.Atoi(resArr[0])
	if err != nil {
		return nil, fmt.Errorf("invalid image height %s", resArr[0])
	}
	w, err := strconv.Atoi(resArr[1])
	if err != nil {
		return nil, fmt.Errorf("invalid image width %s", resArr[0])
	}
	return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}}), nil
}

func saveImage(fileName string, m image.Image) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	return png.Encode(outFile, m)
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
