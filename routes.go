package main

import (
	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/handlers"
	"xelbot.com/reprogl/models"
)

func getRoutes(app *container.Application) *mux.Router {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/article/{slug}", handlers.PageAction(app)).Name("article")
	siteMux.HandleFunc("/{page:[0-9]*}", handlers.IndexAction(app)).Name("blog-page")
	siteMux.HandleFunc("/category/{slug}", handlers.CategoryAction(app)).Name("category-first")
	siteMux.HandleFunc("/category/{slug}/{page:[0-9]+}", handlers.CategoryAction(app)).Name("category")
	siteMux.HandleFunc("/tag/{slug}", handlers.TagAction(app)).Name("tag-first")
	siteMux.HandleFunc("/tag/{slug}/{page:[0-9]+}", handlers.TagAction(app)).Name("tag")
	siteMux.HandleFunc("/info", handlers.InfoAction(app)).Name("info-page")
	siteMux.HandleFunc("/statistika", handlers.StatisticsAction(app)).Name("statistics")
	siteMux.HandleFunc("/robots.txt", handlers.RobotsTXTAction)
	siteMux.HandleFunc("/humans.txt", handlers.HumansTXTAction)
	siteMux.HandleFunc("/favicon.ico", handlers.FavIconAction)
	siteMux.HandleFunc("/sitemap.xml", handlers.SitemapAction(app))
	siteMux.HandleFunc("/feed/atom", handlers.FeedAction(app, models.AtomFeedType)).Name("feed-atom")
	siteMux.HandleFunc("/feed/rss", handlers.FeedAction(app, models.RssFeedType)).Name("feed-rss")
	siteMux.HandleFunc("/add-comment", handlers.AddCommentDummy).Methods("POST").Name("add-comment-dummy")
	siteMux.HandleFunc("/add-ajax-comment", handlers.AddComment(app)).Methods("POST").Name("add-comment")
	siteMux.HandleFunc("/purge-cache", handlers.PurgeCache(app)).Methods("POST")
	siteMux.HandleFunc("/images/avatar/{hash:[0-9A-Z]+}.png", handlers.AvatarGenerator(app))
	siteMux.HandleFunc("/login", handlers.LoginAction(app)).Name("login")

	fragmentsMux := siteMux.PathPrefix("/_fragment").Subrouter()
	fragmentsMux.HandleFunc("/categories", handlers.CategoriesFragment(app)).Name("fragment-categories")
	fragmentsMux.HandleFunc(
		"/comments/{article_id:[0-9]+}/{disabled_flag:(?:e|d)}/{last_time:[0-9]+}",
		handlers.CommentsFragment(app),
	).Name("fragment-comments")
	fragmentsMux.HandleFunc(
		"/recent-posts/{article_id:[0-9]+}",
		handlers.RecentPostsFragment(app),
	).Name("fragment-recent-posts")
	fragmentsMux.HandleFunc("/war-in-ukraine", handlers.DaysOfWarCounter).Name("fragment-war")
	fragmentsMux.HandleFunc("/login-logout", handlers.LoginLogoutLinks(app)).Name("login-logout")

	return siteMux
}
