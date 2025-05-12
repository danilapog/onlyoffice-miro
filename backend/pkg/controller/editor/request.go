package editor

import "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"

type editorRequestParams struct {
	uid  string
	tid  string
	bid  string
	fid  string
	lang string
}

type builderRequest struct {
	Board string
	File  miro.FileInfoResponse
}

func (r builderRequest) ID() string {
	return r.File.ID
}

func (r builderRequest) FolderID() string {
	return r.Board
}

func (r builderRequest) Title() string {
	return r.File.Data.Title
}

func (r builderRequest) URL() string {
	return r.File.Data.DocumentURL
}

func (r builderRequest) ModifiedAt() string {
	return r.File.ModifiedAt
}
