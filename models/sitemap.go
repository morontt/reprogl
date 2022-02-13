package models

import "time"

type SitemapItem struct {
	Slug      string
	UpdatedAt time.Time
}

type SitemapItemList []*SitemapItem
