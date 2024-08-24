package avatar

import (
	"errors"
	"image"
	"image/draw"
	"strconv"

	xdraw "golang.org/x/image/draw"
)

func tryUserSource(id, size int) (image.Image, error) {
	extensions := []string{".png", ".jpg"}
	for _, ext := range extensions {
		filename := "var/data/pictures/user." + strconv.Itoa(id) + ext
		img, err := loadImage(filename)
		if err != nil {
			continue
		}

		imageResult := image.NewRGBA(image.Rect(0, 0, size, size))
		xdraw.BiLinear.Scale(imageResult, imageResult.Bounds(), img, img.Bounds(), draw.Src, nil)

		return imageResult, nil
	}

	return nil, errors.New("")
}
