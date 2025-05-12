package file

import (
	"github.com/golang-jwt/jwt/v5"
)

type convertClaims struct {
	Async      bool   `json:"async,omitempty"`
	FileType   string `json:"filetype"`
	Key        string `json:"key"`
	OutputType string `json:"outputtype"`
	Title      string `json:"title,omitempty"`
	URL        string `json:"url"`
	jwt.RegisteredClaims
}

type createBody struct {
	BoardId  string `json:"board_id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileLang string `json:"file_lang"`
}
