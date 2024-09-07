package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/xelbot/reverse"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/handlers"
	"xelbot.com/reprogl/models"
)

func getRoutes(app *container.Application) *chi.Mux {
	siteMux := chi.NewRouter()
	siteMux.Get(reverse.Add("article", "/article/{slug}"), handlers.PageAction(app))
	siteMux.Get(reverse.Add("home", "/"), handlers.IndexAction(app))
	siteMux.Get(reverse.Add("blog-page", "/{page:[0-9]+}"), handlers.IndexAction(app))
	siteMux.Get(reverse.Add("category-first", "/category/{slug}"), handlers.CategoryAction(app))
	siteMux.Get(reverse.Add("category", "/category/{slug}/{page:[0-9]+}"), handlers.CategoryAction(app))
	siteMux.Get(reverse.Add("tag-first", "/tag/{slug}"), handlers.TagAction(app))
	siteMux.Get(reverse.Add("tag", "/tag/{slug}/{page:[0-9]+}"), handlers.TagAction(app))
	siteMux.Get(reverse.Add("info-page", "/about"), handlers.InfoAction(app))
	siteMux.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/about", http.StatusMovedPermanently)
	})
	siteMux.Get(reverse.Add("statistics", "/statistika"), handlers.StatisticsAction(app))
	siteMux.Get("/robots.txt", handlers.RobotsTXTAction)
	siteMux.Get("/humans.txt", handlers.HumansTXTAction)
	siteMux.Get("/favicon.ico", handlers.FavIconAction)
	siteMux.Get("/headers", handlers.HeadersDebug)
	siteMux.Get("/sitemap.xml", handlers.SitemapAction(app))
	siteMux.Get(reverse.Add("feed-atom", "/feed/atom"), handlers.FeedAction(app, models.AtomFeedType))
	siteMux.Get(reverse.Add("feed-rss", "/feed/rss"), handlers.FeedAction(app, models.RssFeedType))
	siteMux.Post(reverse.Add("add-comment-dummy", "/add-comment"), handlers.AddCommentDummy)
	siteMux.Post(reverse.Add("add-comment", "/add-ajax-comment"), handlers.AddComment(app))
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
	siteMux.Mount("/oauth", oauthMux)

	oauthMux.Get(reverse.Add("oauth-authorize", "/authorize/{provider}"), handlers.OAuthLogin(app))
	oauthMux.Get(reverse.Add("oauth-verification", "/verification/{provider}"), handlers.OAuthCallback(app))
	oauthMux.Get("/check/{request_id}", handlers.OAuthCheckState(app))

	fragmentsMux := chi.NewRouter()
	siteMux.Mount("/_fragment", fragmentsMux)

	fragmentsMux.Get(reverse.Add("fragment-categories", "/categories"), handlers.CategoriesFragment(app))
	fragmentsMux.Get(reverse.Add("fragment-comments",
		"/comments/{article_id:[0-9]+}/{disabled_flag:(?:e|d)}/{last_time:[0-9]+}"),
		handlers.CommentsFragment(app),
	)
	fragmentsMux.Get(reverse.Add("fragment-recent-posts",
		"/recent-posts/{article_id:[0-9]+}"),
		handlers.RecentPostsFragment(app),
	)
	fragmentsMux.Get(reverse.Add("fragment-war", "/war-in-ukraine"), handlers.DaysOfWarCounter)
	fragmentsMux.Get(reverse.Add("auth-navigation", "/auth-navigation"), handlers.AuthNavigation(app))
	fragmentsMux.Get(reverse.Add("fragment-markdown", "/markdown/{filename}"), handlers.MarkdownAction(app))
	fragmentsMux.Get(reverse.Add("menu-auth", "/menu-auth"), handlers.MenuAuth(app))

	return siteMux
}

func getRoutesOld(app *container.Application) *mux.Router {
	siteMux := mux.NewRouter()
	siteMux.Get("test").URL()

	return siteMux
}
