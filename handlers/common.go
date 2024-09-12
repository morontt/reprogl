package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/security"
	"xelbot.com/reprogl/session"
)

func pageOrRedirect(pageString string) (int, bool) {
	var page int

	if pageString == "1" {
		return 1, true
	} else if pageString == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageString)
	}

	return page, false
}

func doESI(w http.ResponseWriter) {
	w.Header().Set("Surrogate-Control", "content=\"ESI/1.0\"")
}

func cacheControl(w http.ResponseWriter, age int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", age))
}

func setExpires(w http.ResponseWriter, t time.Time) {
	t = t.In(time.UTC)
	w.Header().Set("Expires", fmt.Sprintf("%s GMT", t.Format("Mon, 02 Jan 2006 15:04:05")))
}

func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	jsonResult, _ := json.Marshal(data)
	w.Write(jsonResult)
}

func generateRandomToken() string {
	nonce := make([]byte, 18)
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(nonce)
}

func authSuccess(user *models.LoggedUser, app *container.Application, ip string, ctx context.Context) {
	session.SetIdentity(ctx, security.CreateIdentity(user))

	repo := repositories.UserRepository{DB: app.DB}
	if err := repo.SaveLoginEvent(user.ID, ip); err != nil {
		app.LogError(err)
	}
}

func saveLoginReferer(w http.ResponseWriter, r *http.Request) {
	if _, errNoCookie := r.Cookie(session.RefererCookie); errNoCookie == nil {
		return
	}

	host := strings.ReplaceAll(container.GetConfig().Host, ".", "\\.")
	matches := regexp.MustCompile(`^https?:\/\/` + host + `(.*)$`).FindStringSubmatch(r.Referer())
	if matches != nil && matches[1] != "/login" && !strings.HasPrefix(matches[1], "/oauth") {
		session.WriteSessionCookie(w, session.RefererCookie, matches[1], "/")
	}
}

func popLoginReferer(w http.ResponseWriter, r *http.Request) (redirectUrl string, found bool) {
	if cookie, errNoCookie := r.Cookie(session.RefererCookie); errNoCookie == nil {
		redirectUrl = cookie.Value
		found = true

		deleteRefererCookie(w)
	}

	return
}

func deleteRefererCookie(w http.ResponseWriter) {
	session.DeleteCookie(w, session.RefererCookie, "/")
}
