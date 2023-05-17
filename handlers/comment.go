package handlers

import (
	"net/http"
	"strconv"
	"xelbot.com/reprogl/backend"
	"xelbot.com/reprogl/container"
)

type addCommentResponse struct {
	Valid  bool                `json:"valid"`
	Errors []backend.FormError `json:"errors,omitempty"`
}

func AddCommentDummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Silence is gold"))
}

func AddComment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
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

		statusCode := http.StatusCreated

		violations, err := backend.SendComment(commentData)
		if err != nil {
			statusCode = http.StatusBadRequest
		}

		result := addCommentResponse{
			Valid:  true,
			Errors: violations,
		}

		if len(violations) > 0 {
			result.Valid = false
		}

		jsonResponse(w, statusCode, result)
	}
}
