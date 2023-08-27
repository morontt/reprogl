package models

import (
	"encoding/xml"
)

type SitemapItem struct {
	Slug      string      `xml:"-"`
	URL       string      `xml:"loc"`
	UpdatedAt SitemapTime `xml:"lastmod"`
}

type SitemapItemList []*SitemapItem

type SitemapURLSet struct {
	XMLName xml.Name         `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Items   *SitemapItemList `xml:"url"`
}
