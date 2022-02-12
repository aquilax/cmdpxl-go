package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	_ "image/png"
)

func main() {
	fileName := flag.String("f", "", "Path for the file you want to open")
	res := flag.String("res", "", "Image height and width separated by a comma, e.g. 20,10 for a 20x10 image. Note that no spaces can be used.")

	flag.Parse()

	var m image.Image
	var err error

	if *fileName != "" {
		m, err = loadImage(*fileName)
		if err != nil {
			log.Fatal(err)
		}
	} else if *res != "" {
		m, err = createImage(*res)
		if err != nil {
			log.Fatal(err)
		}
		fn := "empty.png"
		fileName = &fn
	} else {
		log.Fatal("need to set either filename or resolution")
	}

	log.Fatal(NewCmdPxl(*fileName, m).Run())
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
