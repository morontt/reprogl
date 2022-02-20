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
	XMLName   xml.Name         `xml:"urlset"`
	Namespace string           `xml:"xmlns,attr"`
	Items     *SitemapItemList `xml:"loc"`
}

func (ct SitemapTime) MarshalText() ([]byte, error) {
	t := time.Time(ct)

	return []byte(t.Format(time.RFC3339)), nil
}
