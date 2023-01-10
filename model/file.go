package model

type File struct {
	Path      string `json:"path"`
	Type      string `json:"type"`
	Published bool   `json:"published"`
}
