package miro

import "context"

type Client interface {
	GetFileInfo(ctx context.Context, req GetFileInfoRequest) (*FileInfoResponse, error)
	GetFilesInfo(ctx context.Context, req GetFilesInfoRequest) (*FilesInfoResponse, error)
	GetFilePublicURL(ctx context.Context, req GetFilePublicURLRequest) (*FileLocationResponse, error)
	GetBoardMember(ctx context.Context, req GetBoardMemberRequest) (*BoardMemberResponse, error)
	CreateFile(ctx context.Context, req CreateFileRequest) (*FileCreatedResponse, error)
	UploadFile(ctx context.Context, req UploadFileRequest) (*FileLocationResponse, error)
}
