package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/api/telegram"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/security"
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
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

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

		var commentator *backend.CommentatorDTO
		var user *backend.UserDTO
		if identity, hasIdentity := session.GetIdentity(r.Context()); hasIdentity {
			user = &backend.UserDTO{
				ID: identity.ID,
			}
		} else {
			commentator = &backend.CommentatorDTO{
				Name:    r.PostForm.Get("name"),
				Email:   r.PostForm.Get("mail"),
				Website: r.PostForm.Get("website"),
			}
		}

		commentData := backend.CommentDTO{
			Commentator: commentator,
			User:        user,
			Text:        r.PostForm.Get("comment_text"),
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

func CommentsFragment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		articleId, err := strconv.Atoi(vars["article_id"])
		if err != nil {
			app.ServerError(w, err)

			return
		}

		repo := repositories.CommentRepository{DB: app.DB}

		var comments *models.CommentList
		var hasIdentity bool
		var identity security.Identity
		if identity, hasIdentity = session.GetIdentity(r.Context()); hasIdentity {
			comments, err = repo.GetCollectionForUsersByArticleId(articleId)
		} else {
			comments, err = repo.GetCollectionByArticleId(articleId)
		}
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := &views.FragmentCommentsData{
			Comments:        comments,
			EnabledComments: vars["disabled_flag"] == models.EnabledComments,
			HasIdentity:     hasIdentity,
			Identity:        identity,
		}

		cacheControl(w, container.DefaultEsiTTL)
		err = views.WriteTemplate(w, "comments.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
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
