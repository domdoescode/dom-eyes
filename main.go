package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/domudall/dom-eyes/eyefinder"
)

var haarCascade = flag.String("haar", "haarcascade_eye.xml", "The location of the Haar Cascade XML configuration to be provided to OpenCV.")

func main() {
	flag.Parse()

	reader, err := os.Open("dom-eye.png")
	if err != nil {
		log.Fatal(err)
	}

	domEye, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	domEye = imaging.Resize(domEye, 38*5, 24*5, imaging.Lanczos)

	file := flag.Arg(0)

	finder := eyefinder.NewFinder(*haarCascade)
	baseImage := loadImage(file)

	eyes := finder.Detect(baseImage)

	bounds := baseImage.Bounds()

	canvas := canvasFromImage(baseImage)

	for _, eye := range eyes {
		rect := rectMargin(-15.0, eye)

		chrisFace := imaging.Fit(domEye, rect.Dx(), rect.Dy(), imaging.Lanczos)

		draw.Draw(
			canvas,
			rect,
			chrisFace,
			bounds.Min,
			draw.Over,
		)
	}

	jpeg.Encode(os.Stdout, canvas, &jpeg.Options{Quality: jpeg.DefaultQuality})
}

func rectMargin(pct float64, rect image.Rectangle) image.Rectangle {
	width := float64(rect.Max.X - rect.Min.X)
	height := float64(rect.Max.Y - rect.Min.Y)

	paddingWidth := int(pct * (width / 100) / 2)
	paddingHeight := int(pct * (height / 100) / 2)

	return image.Rect(
		rect.Min.X-paddingWidth,
		rect.Min.Y-paddingHeight*3,
		rect.Max.X+paddingWidth,
		rect.Max.Y+paddingHeight,
	)
}

func loadImage(file string) image.Image {
	reader, err := os.Open(file)
	if err != nil {
		log.Fatalf("error loading %s: %s", file, err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalf("error loading %s: %s", file, err)
	}
	return img
}

func canvasFromImage(i image.Image) *image.RGBA {
	bounds := i.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, i, bounds.Min, draw.Src)

	return canvas
}
