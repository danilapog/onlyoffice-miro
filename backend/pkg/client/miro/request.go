package miro

import (
	"fmt"
	"strings"
)

type FileUploadRequestData struct {
	URL string `json:"url"`
}

type FileUploadRequest struct {
	Data FileUploadRequestData `json:"data"`
}

type GetFileInfoRequest struct {
	BoardID string `json:"board_id"`
	ItemID  string `json:"item_id"`
	Token   string `json:"token"`
}

func (r *GetFileInfoRequest) Validate() error {
	if strings.TrimSpace(r.BoardID) == "" {
		return fmt.Errorf("boardID is required")
	}

	if strings.TrimSpace(r.ItemID) == "" {
		return fmt.Errorf("itemID is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

type GetFilesInfoRequest struct {
	Cursor  string `json:"cursor"`
	BoardID string `json:"board_id"`
	Token   string `json:"token"`
}

func (r *GetFilesInfoRequest) Validate() error {
	if strings.TrimSpace(r.BoardID) == "" {
		return fmt.Errorf("boardID is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

type GetFilePublicURLRequest struct {
	URL   string
	Token string
}

func (r *GetFilePublicURLRequest) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

type GetBoardMemberRequest struct {
	BoardID  string
	MemberID string
	Token    string
}

func (r *GetBoardMemberRequest) Validate() error {
	if strings.TrimSpace(r.BoardID) == "" {
		return fmt.Errorf("boardID is required")
	}

	if strings.TrimSpace(r.MemberID) == "" {
		return fmt.Errorf("memberID is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

type UploadFileRequest struct {
	BoardID string
	ItemID  string
	FileURL string
	Token   string
}

func (r *UploadFileRequest) Validate() error {
	if strings.TrimSpace(r.BoardID) == "" {
		return fmt.Errorf("boardID is required")
	}

	if strings.TrimSpace(r.ItemID) == "" {
		return fmt.Errorf("itemID is required")
	}

	if strings.TrimSpace(r.FileURL) == "" {
		return fmt.Errorf("fileURL is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

type CreateFileRequest struct {
	BoardID  string
	Name     string
	Type     DocumentType
	Language string
	Token    string
}

func (r *CreateFileRequest) Validate() error {
	if strings.TrimSpace(string(r.Type)) == "" {
		return fmt.Errorf("documentType is required")
	}

	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if strings.TrimSpace(r.BoardID) == "" {
		return fmt.Errorf("boardID is required")
	}

	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token is required")
	}

	if strings.TrimSpace(r.Language) == "" {
		r.Language = "en-US"
	}

	return nil
}
