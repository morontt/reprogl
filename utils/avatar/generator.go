package avatar

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"

	xdraw "golang.org/x/image/draw"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/utils/hashid"
)

var (
	iconsConfig = map[string]map[string]int{
		"male": {
			"eye":        33,
			"hair":       36,
			"mouth":      26,
			"clothes":    64,
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
	validSizes = map[int]bool{
		32:  true,
		80:  true,
		160: true,
		200: true,
	}

	InvalidSize = errors.New("avatar: invalid size")
)

type dataDict map[string]string

func GenerateAvatar(hashData hashid.HashData, app *container.Application, size int) (image.Image, error) {
	if _, ok := validSizes[size]; !ok {
		return nil, InvalidSize
	}

	app.InfoLog.Printf("[IMG] avatar generation by %+v\n", hashData)
	var sourcePrefix string
	if hashData.IsUser() {
		sourcePrefix = "user"
	} else {
		sourcePrefix = "commentator"
	}

	var err error
	imageSrc, err := tryAvatarSource(sourcePrefix, hashData.ID, size)
	if err == nil && imageSrc != nil {
		app.InfoLog.Printf("[IMG] avatar %s found on public/data/pictures\n", hashData.Hash)

		return imageSrc, nil
	}

	var object MaybeGravatar
	if !hashData.IsUser() {
		repository := repositories.CommentRepository{DB: app.DB}
		object, err = repository.FindForGravatar(hashData.ID)
		if err != nil {
			return nil, err
		}
	} else {
		repository := repositories.UserRepository{DB: app.DB}
		object, err = repository.Find(hashData.ID)
		if err != nil {
			return nil, err
		}
	}

	gravatar, err := tryGravatar(object, size, app.InfoLog)
	if err == nil {
		app.InfoLog.Printf("[IMG] avatar %s found on gravatar.com\n", hashData.Hash)

		return gravatar, nil
	}

	return generate8BitIconAvatar(hashData.Hash, hashData.IsMale(), size, app.InfoLog)
}

func generate8BitIconAvatar(hash string, male bool, size int, logger *log.Logger) (image.Image, error) {
	dict := dataByHash(hash, male)
	logger.Printf("[IMG] details for %s: %+v\n", hash, dict)

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

	imageResult := image.NewRGBA(image.Rect(0, 0, size, size))
	xdraw.NearestNeighbor.Scale(imageResult, imageResult.Bounds(), baseImage, baseImage.Bounds(), draw.Src, nil)

	return imageResult, nil
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
