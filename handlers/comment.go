package handlers

import (
	"net/http"
	"strconv"
	"xelbot.com/reprogl/backend"
	"xelbot.com/reprogl/container"
)

func AddCommentDummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Silence is gold"))
}

func AddComment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		//err := r.ParseMultipartForm(512 * 1024)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		commentText := r.PostForm.Get("comment_text")
		nickname := r.PostForm.Get("name")
		email := r.PostForm.Get("mail")
		website := r.PostForm.Get("website")

		topicId, err := strconv.Atoi(r.PostForm.Get("topicId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		parentId, err := strconv.Atoi(r.PostForm.Get("parentId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		commentData := backend.CommentDTO{
			Commentator: backend.CommentatorDTO{
				Name:    nickname,
				Email:   email,
				Website: website,
			},
			Text:      commentText,
			TopicID:   topicId,
			ParentID:  parentId,
			UserAgent: r.UserAgent(),
			IP:        container.RealRemoteAddress(r),
		}
		err = backend.SendComment(commentData)

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Comment created: " + commentText + "\n"))
		w.Write([]byte("Name: " + nickname + "\n"))
		w.Write([]byte("Email: " + email + "\n"))
		w.Write([]byte("Website: " + website + "\n"))

		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "\n"))
		}
	}
}
