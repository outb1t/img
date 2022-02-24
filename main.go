package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func main() {

	files, err := ioutil.ReadDir("./")
	check(err)

	extensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if !extensions[ext] {
			continue
		}

		fmt.Println("Processing ", f.Name())
		bytes := processImage(f.Name())

		f, err := os.Create("output/" + f.Name() + ".txt")
		check(err)
		defer f.Close()

		f.Write(bytes.Bytes())
		f.Sync()
	}
}

func processImage(path string) bytes.Buffer {
	file, err := os.Open(path)
	check(err)

	defer file.Close()

	pixels, err := getPixels(file)
	check(err)

	os.Mkdir("output", os.ModePerm)

	var output bytes.Buffer
	for y := 0; y < len(pixels)-1; y++ {
		for x := 0; x < len(pixels[y])-1; x++ {
			percent := float64(pixels[y][x].G) / 255.0 * 100.0
			pixelData := fmt.Sprintf("x:%v y:%v %v\n", x, y, percent)
			output.WriteString(pixelData)
		}
	}
	return output
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	img, _, err := image.Decode(file)
	check(err)

	img = imaging.Resize(img, 64, 64, imaging.Lanczos)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
