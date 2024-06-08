package models

import (
	"encoding/xml"
	"time"

	"xelbot.com/reprogl/container"
)

type Atom struct {
	XMLName     xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	ID          string      `xml:"id"`
	Title       string      `xml:"title"`
	Description string      `xml:"subtitle"`
	Links       []AtomLink  `xml:"link"`
	Updated     SitemapTime `xml:"updated"`
	Author      AtomPerson  `xml:"author"`
	Generator   string      `xml:"generator"`
	Entry       []AtomEntry `xml:"entry"`
}

type AtomLink struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Language string `xml:"hreflang,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Length   int    `xml:"length,attr,omitempty"`
	MimeType string `xml:"type,attr,omitempty"`
}

type AtomPerson struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type AtomEntry struct {
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	URLs    []AtomLink  `xml:"link"`
	Updated SitemapTime `xml:"updated"`
	Created SitemapTime `xml:"published"`
	Summary AtomSummary `xml:"summary"`
}

type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Text    string   `xml:",cdata"`
}

func (a *Atom) setChannelData(data FeedChannelData) {
	a.ID = "urn:uuid:" + data.GIUD()
	a.Title = data.Title
	a.Description = data.Description
	a.Generator = data.Generator
	a.Links = []AtomLink{
		{Rel: "alternate", Href: data.Link, Language: data.Language},
		{Rel: "self", Href: container.GenerateAbsoluteURL("feed-atom")},
	}
	a.Author = AtomPerson{Name: data.Author, Email: data.Email}

	for _, entry := range data.FeedItems {
		a.addFeedItem(entry)
	}
}

func (a *Atom) ContentType() string {
	return "application/atom+xml; charset=utf-8"
}

func (a *Atom) addFeedItem(entry *FeedItem) {
	links := make([]AtomLink, 0, 2)
	links = append(links, AtomLink{Rel: "alternate", Href: entry.URL})
	if enclosure := entry.GetRssEnclosure(); enclosure != nil {
		links = append(links, AtomLink{
			Rel:      "enclosure",
			Href:     enclosure.Url,
			Length:   enclosure.Length,
			MimeType: enclosure.MimeType,
		})
	}

	a.Entry = append(a.Entry, AtomEntry{
		ID:      "urn:uuid:" + entry.GIUD(),
		Title:   entry.Title,
		URLs:    links,
		Updated: SitemapTime(entry.UpdatedAt),
		Created: SitemapTime(entry.CreatedAt),
		Summary: AtomSummary{Text: entry.Text},
	})

	if entry.CreatedAt.After(time.Time(a.Updated)) {
		a.Updated = SitemapTime(entry.CreatedAt)
	}
}
