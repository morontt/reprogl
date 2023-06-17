package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/api/telegram"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
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

		repo := repositories.ArticleRepository{DB: app.DB}
		article, err := repo.GetByIdForComment(topicId)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

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

		var responseData any
		statusCode := http.StatusCreated

		apiResponse, err := backend.SendComment(commentData)
		if err != nil {
			if errors.Is(err, backend.NotAllowedComment) {
				statusCode = http.StatusOK
				responseData = addCommentResponse{
					Valid: false,
					Errors: []backend.FormError{
						{
							Path:    "comment_text",
							Message: "Добавление комментариев тут отключено 😐",
						},
					},
				}
			} else {
				app.LogError(err)
				statusCode = http.StatusBadRequest
			}
		}

		if apiResponse != nil {
			result := addCommentResponse{
				Valid:  true,
				Errors: apiResponse.Violations,
			}

			if apiResponse.Violations != nil && len(apiResponse.Violations) > 0 {
				result.Valid = false
			} else {
				if apiResponse.Comment != nil {
					go afterCommentHook(app, apiResponse.Comment, article)
				}
			}

			responseData = result
		}

		jsonResponse(w, statusCode, responseData)
	}
}

func afterCommentHook(
	app *container.Application,
	comment *backend.CreatedCommentDTO,
	article *models.ArticleForComment,
) {
	var updatedComment *backend.CreatedCommentDTO

	backend.PingGeolocation()
	refreshedComment, err := backend.RefreshComment(comment.ID)
	if err != nil {
		updatedComment = comment
	} else {
		updatedComment = refreshedComment
	}

	telegram.SendNotification(app, updatedComment, article)
}
