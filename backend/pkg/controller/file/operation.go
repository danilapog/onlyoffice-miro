package file

import (
	"context"
	"errors"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/controller/base"
	"github.com/labstack/echo/v4"
)

func toDocumentType(ftype string) miro.DocumentType {
	switch ftype {
	case string(miro.DOCX):
		return miro.DOCX
	case string(miro.PPTX):
		return miro.PPTX
	case string(miro.XLSX):
		return miro.XLSX
	default:
		return miro.DOCX
	}
}

func PrepareRequest(
	ctx echo.Context,
	tctx context.Context,
	c *base.BaseController,
) (*boardAuthenticationResponse, error) {
	token, err := c.ExtractUserToken(ctx)
	if err != nil {
		return nil, c.HandleError(ctx, err, http.StatusForbidden, "failed to extract authentication parameters")
	}

	bid, err := c.GetQueryParam(ctx, "bid")
	if err != nil {
		return nil, c.HandleError(ctx, err, http.StatusBadRequest, "board id parameter is missing")
	}

	_, auth, err := c.FetchAuthenticationWithSettings(tctx, token.User, token.Team, bid)
	if err != nil {
		if errors.Is(err, base.ErrMissingAuthentication) {
			return nil, c.HandleError(ctx, err, http.StatusUnauthorized, "could not retrieve authentication")
		}

		if errors.Is(err, base.ErrSettingsNotConfigured) {
			return nil, c.HandleError(ctx, err, http.StatusConflict, "could not retrieve document editor settigns")
		}

		return nil, c.HandleError(ctx, err, http.StatusBadRequest, "could not retrieve required data")
	}

	return &boardAuthenticationResponse{
		BoardID:        bid,
		Authentication: auth,
	}, nil
}

func CreateFile(
	ctx echo.Context,
	tctx context.Context,
	c *base.BaseController,
	req miro.CreateFileRequest,
) (*miro.FileLocationResponse, error) {
	response, err := c.MiroClient.CreateFile(tctx, req)
	if err != nil {
		return nil, c.HandleError(ctx, err, http.StatusInternalServerError, "failed to create a file")
	}

	return response, nil
}

func GetFileInfo(
	ctx echo.Context,
	tctx context.Context,
	c *base.BaseController,
	boardID string,
	fileID string,
	accessToken string,
) (*miro.FileInfoResponse, error) {
	file, err := c.MiroClient.GetFileInfo(tctx, miro.GetFileInfoRequest{
		BoardID: boardID,
		ItemID:  fileID,
		Token:   accessToken,
	})

	if err != nil {
		return nil, c.HandleError(ctx, err, http.StatusBadRequest, "failed to fetch miro file")
	}

	return file, nil
}

func GetFilesInfo(
	ctx echo.Context,
	tctx context.Context,
	c *base.BaseController,
	boardID string,
	cursor string,
	accessToken string,
) (*miro.FilesInfoResponse, error) {
	files, err := c.MiroClient.GetFilesInfo(tctx, miro.GetFilesInfoRequest{
		Cursor:  cursor,
		BoardID: boardID,
		Token:   accessToken,
	})

	if err != nil {
		return nil, c.HandleError(ctx, err, http.StatusBadRequest, "failed to fetch miro files")
	}

	return files, nil
}
