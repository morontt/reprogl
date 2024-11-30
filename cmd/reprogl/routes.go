package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xelbot/reverse"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/handlers"
	"xelbot.com/reprogl/models"
)

func getRoutes(app *container.Application) *chi.Mux {
	siteMux := chi.NewRouter()
	siteMux.Get(reverse.Add("article", "/article/{slug}"),
		app.Metrics.Duration("GET /article/:slug", handlers.PageAction(app)))
	siteMux.Get(reverse.Add("home", "/"),
		app.Metrics.Duration("GET /:page", handlers.IndexAction(app)))
	siteMux.Get(reverse.Add("blog-page", "/{page:[0-9]+}"),
		app.Metrics.Duration("GET /:page", handlers.IndexAction(app)))
	siteMux.Get(reverse.Add("category-first", "/category/{slug}"),
		app.Metrics.Duration("GET /category/:slug/:page", handlers.CategoryAction(app)))
	siteMux.Get(reverse.Add("category", "/category/{slug}/{page:[0-9]+}"),
		app.Metrics.Duration("GET /category/:slug/:page", handlers.CategoryAction(app)))
	siteMux.Get(reverse.Add("tag-first", "/tag/{slug}"),
		app.Metrics.Duration("GET /tag/:slug/:page", handlers.TagAction(app)))
	siteMux.Get(reverse.Add("tag", "/tag/{slug}/{page:[0-9]+}"),
		app.Metrics.Duration("GET /tag/:slug/:page", handlers.TagAction(app)))
	siteMux.Get(reverse.Add("info-page", "/about"),
		app.Metrics.Duration("GET /about", handlers.InfoAction(app)))
	siteMux.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/about", http.StatusMovedPermanently)
	})
	siteMux.Get(reverse.Add("statistics", "/statistika"),
		app.Metrics.Duration("GET /statistika", handlers.StatisticsAction(app)))
	siteMux.Get("/robots.txt", handlers.RobotsTXTAction)
	siteMux.Get("/humans.txt", handlers.HumansTXTAction)
	siteMux.Get("/favicon.ico", handlers.FavIconAction)
	siteMux.Get("/headers", handlers.HeadersDebug)
	siteMux.Get("/sitemap.xml", handlers.SitemapAction(app))
	siteMux.Get(reverse.Add("feed-atom", "/feed/atom"),
		app.Metrics.Duration("GET /feed/atom", handlers.FeedAction(app, models.AtomFeedType)))
	siteMux.Get(reverse.Add("feed-rss", "/feed/rss"),
		app.Metrics.Duration("GET /feed/rss", handlers.FeedAction(app, models.RssFeedType)))
	siteMux.Post(reverse.Add("add-comment-dummy", "/add-comment"), handlers.AddCommentDummy)
	siteMux.Post(reverse.Add("add-comment", "/add-ajax-comment"),
		app.Metrics.Duration("POST /add-ajax-comment", handlers.AddComment(app)))
	siteMux.Post("/purge-cache", handlers.PurgeCache(app))
	siteMux.Get("/images/avatar/{hash:[0-9A-Z]+}.png", handlers.AvatarGenerator(app))
	siteMux.Get("/images/avatar/{hash:[0-9A-Z]+}.w{size:[0-9]+}.png", handlers.AvatarGeneratorWithSize(app))
	siteMux.Get(reverse.Add("profile", "/profile"), handlers.ProfileAction(app))
	siteMux.Post("/profile", handlers.UpdateProfile(app))
	siteMux.Get("/email-unsubscribe/{hash:[0-9A-Z]+}", handlers.EmailUnsubscribe(app))
	siteMux.Post("/email-unsubscribe/{hash:[0-9A-Z]+}", handlers.EmailUnsubscribePost(app))
	siteMux.Get(reverse.Add("login", "/login"), handlers.LoginAction(app))
	siteMux.Post("/login", handlers.LoginCheck(app))
	siteMux.Get(reverse.Add("logout", "/logout"), handlers.LogoutAction)

	oauthMux := chi.NewRouter()
	siteMux.Mount(reverse.Group("oauth", "/oauth"), oauthMux)

	oauthMux.Get(reverse.AddGr("oauth-authorize", "oauth", "/authorize/{provider}"), handlers.OAuthLogin(app))
	oauthMux.Get(reverse.AddGr("oauth-verification", "oauth", "/verification/{provider}"), handlers.OAuthCallback(app))
	oauthMux.Get(reverse.AddGr("oauth-check", "oauth", "/check/{request_id}"), handlers.OAuthCheckState(app))

	fragmentsMux := chi.NewRouter()
	siteMux.Mount(reverse.Group("fragments", "/_fragment"), fragmentsMux)

	fragmentsMux.Get(reverse.AddGr("fragment-categories", "fragments", "/categories"), handlers.CategoriesFragment(app))
	fragmentsMux.Get(reverse.AddGr("fragment-comments",
		"fragments",
		"/comments/{article_id:[0-9]+}/{disabled_flag:(?:e|d)}/{last_time:[0-9]+}"),
		app.Metrics.Duration("GET /_fragment/comments/:article", handlers.CommentsFragment(app)),
	)
	fragmentsMux.Get(reverse.AddGr("fragment-recent-posts",
		"fragments",
		"/recent-posts/{article_id:[0-9]+}"),
		handlers.RecentPostsFragment(app),
	)
	fragmentsMux.Get(reverse.AddGr("fragment-war", "fragments", "/war-in-ukraine"), handlers.DaysOfWarCounter)
	fragmentsMux.Get(reverse.AddGr("auth-navigation", "fragments", "/auth-navigation"), handlers.AuthNavigation(app))
	fragmentsMux.Get(reverse.AddGr("fragment-markdown", "fragments", "/markdown/{filename}"), handlers.MarkdownAction(app))
	fragmentsMux.Get(reverse.AddGr("menu-auth", "fragments", "/menu-auth"), handlers.MenuAuth(app))

	return siteMux
}
