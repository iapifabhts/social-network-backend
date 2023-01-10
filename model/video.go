package model

import (
	"github.com/iapifabhts/social-network-backend/validator"
)

type Video struct {
	Content
	Description string `json:"description"`
	PosterPath  string `json:"posterPath"`
	File        File   `json:"file"`
}

type VideoCreation struct {
	Title       string `json:"title"`
	CreatorID   string `json:"creatorID"`
	Description string `json:"description"`
	PosterPath  string `json:"posterPath"`
	FilePath    string `json:"filePath"`
}

func (v VideoCreation) Valid() error {
	val := validator.New()
	val.Verify(len(v.Title) == 0, "")
	val.Verify(len(v.PosterPath) == 0, "")
	val.Verify(len(v.FilePath) == 0, "")
	val.Verify(len(v.CreatorID) == 0, "")
	return val.Verdict()
}
