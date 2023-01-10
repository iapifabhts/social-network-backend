package model

import (
	"github.com/iapifabhts/social-network-backend/validator"
)

type Article struct {
	Content
	Description string `json:"description"`
	Files       []File `json:"files"`
}

type ArticleCreation struct {
	Title       string   `json:"title"`
	CreatorID   string   `json:"creatorID"`
	Description string   `json:"description"`
	FilePaths   []string `json:"filePaths"`
}

func (a ArticleCreation) Valid() error {
	val := validator.New()
	val.Verify(len(a.Title) == 0, "")
	val.Verify(len(a.CreatorID) == 0, "")
	val.Verify(len(a.FilePaths) == 0, "")
	return val.Verdict()
}
