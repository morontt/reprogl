package models

import (
	"encoding/xml"
	"time"
)

type SitemapTime time.Time

type SitemapItem struct {
	XMLName   xml.Name    `xml:"url"`
	Slug      string      `xml:"-"`
	URL       string      `xml:"loc"`
	UpdatedAt SitemapTime `xml:"lastmod"`
}

type SitemapItemList []*SitemapItem

type SitemapURLSet struct {
	XMLName xml.Name         `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Items   *SitemapItemList `xml:"loc"`
}

func (ct SitemapTime) MarshalText() ([]byte, error) {
	t := time.Time(ct)

	return []byte(t.Format(time.RFC3339)), nil
}
