package handlers

import (
	"net/http"
	"strconv"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

func AddCommentDummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Silence is gold"))
}

func AddComment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//err := r.ParseForm()
		err := r.ParseMultipartForm(512 * 1024)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		commentText := r.PostForm.Get("comment_text")
		nickname := r.PostForm.Get("name")
		email := r.PostForm.Get("mail")
		website := r.PostForm.Get("website")

		_, err = strconv.Atoi(r.PostForm.Get("topicId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		_, err = strconv.Atoi(r.PostForm.Get("parentId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		_, wsse := security.GetWSSEHeader()

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Comment created: " + commentText + "\n"))
		w.Write([]byte("Name: " + nickname + "\n"))
		w.Write([]byte("Email: " + email + "\n"))
		w.Write([]byte("Website: " + website + "\n"))
		w.Write([]byte(wsse + "\n"))
	}
}
