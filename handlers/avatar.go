package handlers

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/utils/avatar"
	"xelbot.com/reprogl/utils/hashid"
)

func AvatarGenerator(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		writeAvatar(w, app, hash, 80)
	}
}

func AvatarGeneratorWithSize(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]
		sizeStr := vars["size"]

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			app.NotFound(w)
			return
		}

		writeAvatar(w, app, hash, size)
	}
}

func writeAvatar(w http.ResponseWriter, app *container.Application, hash string, size int) {
	hashModel, err := hashid.Decode(hash, true)
	if err != nil {
		app.NotFound(w)
		return
	}

	img, err := avatar.GenerateAvatar(hashModel, app, size)
	if err != nil {
		if errors.Is(err, models.RecordNotFound) || errors.Is(err, avatar.InvalidSize) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}

		return
	}

	expires := time.Now().Add(container.AvatarTTL * time.Second)

	cacheControl(w, container.AvatarTTL)
	setExpires(w, expires)
	w.Header().Set("Content-Type", "image/png")

	rawImage, err := pngquantPipe(img)
	if err != nil {
		panic(err)
	}

	saveToFile(hash, size, rawImage)

	_, err = w.Write(rawImage)
	if err != nil {
		app.ServerError(w, err)
	}
}

func pngquantPipe(avatarImage image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, avatarImage)
	if err != nil {
		return []byte{}, err
	}

	quantBuf := new(bytes.Buffer)

	cmd := exec.Command("/usr/bin/pngquant", "-s1", "--quality=60-80", "-")
	cmd.Stdin = bufio.NewReader(buf)
	cmd.Stderr = os.Stderr
	cmd.Stdout = quantBuf

	err = cmd.Run()
	if err != nil {
		return []byte{}, err
	}

	return quantBuf.Bytes(), nil
}

func saveToFile(hash string, size int, rawImage []byte) {
	var postfix string
	if size != 80 {
		postfix = ".w" + strconv.Itoa(size)
	}

	filename := "public/images/avatar/" + hash + postfix + ".png"
	f, err := os.Open(filename)
	if errors.Is(err, fs.ErrNotExist) {
		_ = os.WriteFile(filename, rawImage, 0644)
	}
	_ = f.Close()
}
