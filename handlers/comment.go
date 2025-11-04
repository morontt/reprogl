package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/api/telegram"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
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
		topicId, err := strconv.Atoi(r.PostFormValue("topicId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		parentId, err := strconv.Atoi(r.PostFormValue("parentId"))
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		article, err := repo.GetByIdForComment(topicId >> 7)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		var commentator *backend.CommentatorDTO
		var user *backend.CommentUserDTO
		if identity, hasIdentity := session.GetIdentity(r.Context()); hasIdentity {
			user = &backend.CommentUserDTO{
				ID: identity.ID,
			}
		} else {
			commentator = &backend.CommentatorDTO{
				Name:    r.PostFormValue("name"),
				Email:   r.PostFormValue("mail"),
				Website: r.PostFormValue("website"),
			}
		}

		commentData := backend.CommentDTO{
			Commentator: commentator,
			User:        user,
			Text:        r.PostFormValue("comment_text"),
			TopicID:     topicId,
			ParentID:    parentId,
			UserAgent:   r.UserAgent(),
			IP:          container.RealRemoteAddress(r),
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

			if len(apiResponse.Violations) > 0 {
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

func CommentsFragment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		articleId, err := strconv.Atoi(chi.URLParam(r, "article_id"))
		if err != nil {
			app.ServerError(w, err)

			return
		}

		repo := repositories.CommentRepository{DB: app.DB}

		var comments models.CommentList
		identity, found := session.GetIdentity(r.Context())
		if found && identity.IsAdmin() {
			comments, err = repo.GetCollectionWithExtraDataByArticleId(articleId)
		} else {
			comments, err = repo.GetCollectionByArticleId(articleId)
		}
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := &views.FragmentCommentsData{
			Comments:        comments,
			EnabledComments: chi.URLParam(r, "disabled_flag") == models.EnabledComments,
		}

		cacheControl(w, container.DefaultEsiTTL)
		err, wh := views.WriteTemplateWithContext(r.Context(), w, "comments.gohtml", templateData)
		if err != nil {
			if wh {
				app.LogError(err)
			} else {
				app.ServerError(w, err)
			}
		}
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
