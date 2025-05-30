package views

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
	"xelbot.com/reprogl/session"
)

const defaultPageSize = 64 * 1024

var (
	//go:embed templates markdown
	sources embed.FS

	templates map[string]*template.Template
)

func init() {
	templates = make(map[string]*template.Template)
}

func LoadViewSet() error {
	templatesMap := map[string][]string{
		"about.gohtml": {
			"templates/about.gohtml",
			"templates/partials/menu.gohtml",
			"templates/partials/sticky-header.gohtml",
			"templates/partials/big-header.gohtml",
			"templates/partials/footer.gohtml",
			"templates/partials/social-icons.gohtml",
			"templates/layout/svg-sprites.gohtml",
			"templates/layout/base.gohtml",
		},
		"article.gohtml": {
			"templates/article.gohtml",
			"templates/partials/author-info.gohtml",
			"templates/partials/comment-form.gohtml",
			"templates/partials/menu.gohtml",
			"templates/partials/sidebar.gohtml",
			"templates/partials/social-icons.gohtml",
			"templates/partials/powered-by.gohtml",
			"templates/layout/svg-auth.gohtml",
			"templates/layout/svg-sprites.gohtml",
			"templates/layout/base.gohtml",
		},
		"statistics.gohtml": {
			"templates/statistics.gohtml",
			"templates/partials/author-info.gohtml",
			"templates/partials/menu.gohtml",
			"templates/partials/sidebar.gohtml",
			"templates/partials/social-icons.gohtml",
			"templates/partials/powered-by.gohtml",
			"templates/layout/svg-sprites.gohtml",
			"templates/layout/base.gohtml",
		},
		"profile.gohtml": {
			"templates/profile.gohtml",
			"templates/partials/profile-form.gohtml",
			"templates/partials/menu.gohtml",
			"templates/partials/sidebar.gohtml",
			"templates/partials/social-icons.gohtml",
			"templates/layout/svg-sprites.gohtml",
			"templates/layout/base.gohtml",
		},
		"index.gohtml": {
			"templates/index.gohtml",
			"templates/partials/menu.gohtml",
			"templates/partials/sticky-header.gohtml",
			"templates/partials/big-header.gohtml",
			"templates/partials/footer.gohtml",
			"templates/partials/social-icons.gohtml",
			"templates/layout/svg-sprites.gohtml",
			"templates/layout/base.gohtml",
		},
		"categories.gohtml": {
			"templates/fragments/categories.gohtml",
		},
		"comments.gohtml": {
			"templates/fragments/comments.gohtml",
		},
		"recent-posts.gohtml": {
			"templates/fragments/recent-posts.gohtml",
		},
		"login.gohtml": {
			"templates/auth/login.gohtml",
		},
		"auth-navigation.gohtml": {
			"templates/fragments/auth-navigation.gohtml",
		},
		"menu-auth.gohtml": {
			"templates/fragments/menu-auth.gohtml",
		},
		"unsubscribe.gohtml": {
			"templates/unsubscribe.gohtml",
		},
		"oauth-pending.gohtml": {
			"templates/oauth/oauth-pending.gohtml",
		},
	}

	customFunctions := template.FuncMap{
		"raw":           rawHTML,
		"is_dev":        func() bool { return container.IsDevMode() },
		"path":          urlGenerator,
		"abs_path":      absUrlGenerator,
		"render_esi":    renderESI,
		"tags":          tags,
		"cdn":           cdnBase,
		"nl2br":         nl2br,
		"author_bio":    authorBio,
		"author_data":   authorDataPart,
		"author_adr":    authorLocation,
		"author_job":    authorJob,
		"author_avatar": authorAvatar,
		"substr":        subString,
		"time_tag":      timeTag,
		"asset":         assetTag,
		"go_version":    func() string { return container.GoVersionNumbers },
		"commit_hash":   func() string { return container.GitRevision },
		"cnt_comments":  commentsCountString,
		"cnt_times":     timesCountString,
		"emojiFlag":     emojiFlag,

		"articleStyles":    articleStyles,
		"statisticsStyles": statisticsStyles,
		"indexStyles":      indexStyles,
		"infoStyles":       infoStyles,
		"profileStyles":    profileStyles,

		"avatar_picture": avatarPictureTag,
	}

	for key, files := range templatesMap {
		tmpl, err := template.New(key).Funcs(customFunctions).ParseFS(sources, files...)
		if err != nil {
			return err
		}

		templates[key] = tmpl
	}

	return nil
}

func RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("the template %s does not exist", name)
	}

	var buf strings.Builder
	buf.Grow(defaultPageSize)

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func WriteTemplate(w http.ResponseWriter, name string, data any) error {
	content, err := RenderTemplate(name, data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Surrogate-Control", "content=\"ESI/1.0\"")
	_, err = w.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func WriteTemplateWithContext(ctx context.Context, w http.ResponseWriter, name string, data any) error {
	if flashObjectPart, ok := data.(DataWithFlashMessage); ok {
		if flashSuccessMessage, found := session.Pop[string](ctx, session.FlashSuccessKey); found {
			flashObjectPart.SetSuccessFlash(flashSuccessMessage)
		}
	}

	if identityPart, ok := data.(security.IdentityAware); ok {
		identity, _ := session.GetIdentity(ctx)
		identityPart.SetIdentity(identity)
	}

	return WriteTemplate(w, name, data)
}
