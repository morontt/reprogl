package views

import "xelbot.com/reprogl/models"

type ArticlePageData struct {
	Article *models.Article
}

type IndexPageData struct {
	PageNumber int
	Articles   models.ArticleList
}
