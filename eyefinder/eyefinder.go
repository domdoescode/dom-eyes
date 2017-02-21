package eyefinder

import (
	"image"

	"github.com/lazywei/go-opencv/opencv"
)

var faceCascade *opencv.HaarCascade

type Finder struct {
	cascade *opencv.HaarCascade
}

func NewFinder(xml string) *Finder {
	return &Finder{
		cascade: opencv.LoadHaarClassifierCascade(xml),
	}
}

func (f *Finder) Detect(i image.Image) []image.Rectangle {
	var output []image.Rectangle

	objects := f.cascade.DetectObjects(opencv.FromImage(i))
	for _, object := range objects {
		output = append(output, image.Rectangle{
			image.Point{object.X(), object.Y()},
			image.Point{object.X() + object.Width(), object.Y() + object.Height()},
		})
	}

	return output
}
