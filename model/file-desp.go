package model

type FileDesp struct {
	Name  string `json:"name"`
	Last  string `json:"last"`
	Size  string `json:"size"`
	IsDir bool   `json:"is_dir"`
}
