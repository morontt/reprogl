package models

import (
	"database/sql"
	"encoding/json"
)

type FeaturedImage struct {
	ImagePath  sql.NullString
	Width      sql.NullInt32
	PictureTag sql.NullString
	SrcSet     sql.NullString
	Alt        sql.NullString
}

type SrcImage struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SrcSetItem struct {
	Items []SrcImage `json:"items"`
	Type  string     `json:"type"`
}

func (i *FeaturedImage) HasImage() bool {
	return i.PictureTag.Valid || i.ImagePath.Valid
}

func (i *FeaturedImage) HasWebp() bool {
	if !i.SrcSet.Valid {
		return false
	}

	srcSet := i.DecodeSrcSet()
	_, found := srcSet["webp"]

	return found
}

func (i *FeaturedImage) HasAvif() bool {
	if !i.SrcSet.Valid {
		return false
	}

	srcSet := i.DecodeSrcSet()
	_, found := srcSet["avif"]

	return found
}

func (i *FeaturedImage) DecodeSrcSet() map[string]SrcSetItem {
	data := make(map[string]SrcSetItem)
	if i.SrcSet.Valid {
		raw := []byte(i.SrcSet.String)
		if json.Valid(raw) {
			err := json.Unmarshal(raw, &data)
			if err != nil {
				panic(err)
			}
		}
	}

	return data
}

func (i *FeaturedImage) SrcImageForOpenGraph() *SrcImage {
	srcSet := i.DecodeSrcSet()
	if srcSetItem, found := srcSet["origin"]; found {
		return &srcSetItem.Items[0]
	}

	return nil
}
