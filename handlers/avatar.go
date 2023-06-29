package handlers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"
	"net/http"
	"os"
	"os/exec"

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

		hashModel, err := hashid.Decode(hash)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		img, err := avatar.GenerateAvatar(hashModel, app)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		cacheControl(w, container.AvatarTTL)
		w.Header().Set("Content-Type", "image/png")

		err = pngquantPipe(w, img, r.Header.Values("If-None-Match"))
		if err != nil {
			panic(err)
		}
	}
}

func pngquantPipe(w http.ResponseWriter, avatarImage image.Image, etags []string) error {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, avatarImage)
	if err != nil {
		return err
	}

	crc32cs := crc32.Checksum(buf.Bytes(), crc32.MakeTable(crc32.Castagnoli))
	chechSum := fmt.Sprintf("%08X", crc32cs)

	w.Header().Set("Etag", chechSum)
	for _, item := range etags {
		if chechSum == item {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	quantBuf := new(bytes.Buffer)

	cmd := exec.Command("/usr/bin/pngquant", "-s1", "--quality=60-80", "-")
	cmd.Stdin = bufio.NewReader(buf)
	cmd.Stderr = os.Stderr
	cmd.Stdout = quantBuf

	err = cmd.Run()
	if err != nil {
		return err
	}

	_, err = w.Write(quantBuf.Bytes())
	return err
}
