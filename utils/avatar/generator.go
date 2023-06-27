package avatar

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"strconv"
)

var iconsConfig = map[string]map[string]int{
	"male": {
		"eye":        33,
		"hair":       37,
		"mouth":      26,
		"clothes":    66,
		"face":       4,
		"background": 5,
	},
	"female": {
		"eye":        53,
		"head":       34,
		"mouth":      17,
		"clothes":    59,
		"face":       4,
		"background": 5,
	},
}

type dataDict map[string]string

func GenerateAvatar(hash string, male bool) (image.Image, error) {
	dict := dataByHash(hash, male)
	baseImage := image.NewRGBA(image.Rect(0, 0, 400, 400))

	bgImage, err := loadImage(filePath(dict, "background"))
	if err != nil {
		return nil, err
	}
	faceImage, err := loadImage(filePath(dict, "face"))
	if err != nil {
		return nil, err
	}
	clothesImage, err := loadImage(filePath(dict, "clothes"))
	if err != nil {
		return nil, err
	}
	mouthImage, err := loadImage(filePath(dict, "mouth"))
	if err != nil {
		return nil, err
	}

	var hairImage image.Image
	if male {
		hairImage, err = loadImage(filePath(dict, "hair"))
		if err != nil {
			return nil, err
		}
	} else {
		hairImage, err = loadImage(filePath(dict, "head"))
		if err != nil {
			return nil, err
		}
	}

	eyeImage, err := loadImage(filePath(dict, "eye"))
	if err != nil {
		return nil, err
	}

	draw.Draw(baseImage, bgImage.Bounds(), bgImage, image.Point{}, draw.Src)
	draw.Draw(baseImage, faceImage.Bounds(), faceImage, image.Point{}, draw.Over)
	draw.Draw(baseImage, clothesImage.Bounds(), clothesImage, image.Point{}, draw.Over)
	draw.Draw(baseImage, mouthImage.Bounds(), mouthImage, image.Point{}, draw.Over)
	draw.Draw(baseImage, hairImage.Bounds(), hairImage, image.Point{}, draw.Over)
	draw.Draw(baseImage, eyeImage.Bounds(), eyeImage, image.Point{}, draw.Over)

	return baseImage, nil
}

func dataByHash(hash string, male bool) dataDict {
	var gender string
	md5bytes := md5.Sum([]byte(hash))

	if male {
		gender = "male"
	} else {
		gender = "female"
	}

	facePart := int(md5bytes[0])
	backgroundPart := (int(md5bytes[1]) << 16) + (int(md5bytes[2]) << 8) + int(md5bytes[3])
	eyePart := (int(md5bytes[4]) << 16) + (int(md5bytes[5]) << 8) + int(md5bytes[6])
	hairPart := (int(md5bytes[7]) << 16) + (int(md5bytes[8]) << 8) + int(md5bytes[9])
	mouthPart := (int(md5bytes[10]) << 16) + (int(md5bytes[11]) << 8) + int(md5bytes[12])
	clothesPart := (int(md5bytes[13]) << 16) + (int(md5bytes[14]) << 8) + int(md5bytes[15])

	face := 1 + facePart%iconsConfig[gender]["face"]
	eye := 1 + eyePart%iconsConfig[gender]["eye"]
	mouth := 1 + mouthPart%iconsConfig[gender]["mouth"]
	clothes := 1 + clothesPart%iconsConfig[gender]["clothes"]
	background := 1 + backgroundPart%iconsConfig[gender]["background"]

	dict := dataDict{
		"gender":     gender,
		"face":       strconv.Itoa(face),
		"eye":        strconv.Itoa(eye),
		"mouth":      strconv.Itoa(mouth),
		"clothes":    strconv.Itoa(clothes),
		"background": strconv.Itoa(background),
	}

	if male {
		dict["hair"] = strconv.Itoa(1 + hairPart%iconsConfig[gender]["hair"])
	} else {
		dict["head"] = strconv.Itoa(1 + hairPart%iconsConfig[gender]["head"])
	}

	return dict
}

func filePath(dict dataDict, key string) string {
	return fmt.Sprintf("./var/8biticon/%s/%s%s.png", dict["gender"], key, dict[key])
}

func loadImage(filePath string) (image.Image, error) {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()

	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}
