package miro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/assets"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
)

type client struct {
	baseUrl    string
	httpClient *http.Client
	errors     *Errors
}

func NewMiroClient(config *config.MiroConfig) Client {
	return &client{
		baseUrl: config.BaseURL,
		errors:  NewErrors(),
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: common.NewRetryableTransport(&http.Transport{
				MaxIdleConnsPerHost:    100,
				IdleConnTimeout:        90 * time.Second,
				MaxResponseHeaderBytes: 1 << 20,
				DisableCompression:     false,
				ForceAttemptHTTP2:      true,
			},
			),
		},
	}
}

func (c *client) buildURL(paths ...string) string {
	var sb strings.Builder
	sb.WriteString(c.baseUrl)
	sb.WriteString("/")
	sb.WriteString(strings.Join(paths, "/"))
	return sb.String()
}

func (c *client) sendRequest(
	ctx context.Context,
	method, url, token string,
	body io.Reader,
	headers map[string]string,
	result any,
) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return c.errors.FailedToCreateRequest(err)
	}

	req.Header.Set("Authorization", common.Concat("Bearer ", token))
	req.Header.Set("Accept", "application/json")
	if body != nil && headers == nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return c.errors.FailedToSendRequest(err)
	}

	defer res.Body.Close()
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return c.errors.RequestFailed(res.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			return c.errors.FailedToDecodeResponse(err)
		}
	}

	return nil
}

func (c *client) GetFileInfo(ctx context.Context, req GetFileInfoRequest) (*FileInfoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var response FileInfoResponse
	url := c.buildURL("boards", req.BoardID, "items", req.ItemID)
	if err := c.sendRequest(ctx, http.MethodGet, url, req.Token, nil, nil, &response); err != nil {
		return nil, c.errors.FailedToGetFileInfo(err)
	}

	return &response, nil
}

func (c *client) GetFilesInfo(ctx context.Context, req GetFilesInfoRequest) (*FilesInfoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var response FilesInfoResponse
	url := c.buildURL("boards", req.BoardID, "items?type=document&limit=50")
	if req.Cursor != "" {
		var sb strings.Builder
		sb.WriteString(url)
		sb.WriteString("&cursor=")
		sb.WriteString(req.Cursor)
		url = sb.String()
	}

	if err := c.sendRequest(ctx, http.MethodGet, url, req.Token, nil, nil, &response); err != nil {
		return nil, c.errors.FailedToGetFileInfo(err)
	}

	return &response, nil
}

func (c *client) GetFilePublicURL(ctx context.Context, req GetFilePublicURLRequest) (*FileLocationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var response FileLocationResponse
	if err := c.sendRequest(ctx, http.MethodGet, req.URL, req.Token, nil, nil, &response); err != nil {
		return nil, c.errors.FailedToGetFileURL(err)
	}

	return &response, nil
}

func (c *client) GetBoardMember(ctx context.Context, req GetBoardMemberRequest) (*BoardMemberResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var response BoardMemberResponse
	url := c.buildURL("boards", req.BoardID, "members", req.MemberID)
	if err := c.sendRequest(ctx, http.MethodGet, url, req.Token, nil, nil, &response); err != nil {
		return nil, c.errors.FailedToGetBoardMember(err)
	}

	return &response, nil
}

func (c *client) UploadFile(ctx context.Context, req UploadFileRequest) (*FileLocationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	body := FileUploadRequest{
		Data: FileUploadRequestData{
			URL: req.FileURL,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, c.errors.FailedToMarshalRequest(err)
	}

	url := c.buildURL("boards", req.BoardID, "documents", req.ItemID)
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var response FileLocationResponse
	if err := c.sendRequest(ctx, http.MethodPatch, url, req.Token, bytes.NewBuffer(payload), headers, &response); err != nil {
		return nil, c.errors.FailedToUploadFile(err)
	}

	return &response, nil
}

func (c *client) createMultipartForm(data map[string]any, filePath string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, "", c.errors.FailedToMarshalRequest(err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, "data"))
	h.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, "", c.errors.FailedToCreateFormFile(err)
	}

	if _, err := part.Write(payload); err != nil {
		return nil, "", c.errors.FailedToWriteFileData(err)
	}

	fileData, err := assets.Templates.ReadFile(filePath)
	if err != nil {
		return nil, "", c.errors.FailedToReadFile(err)
	}

	fileWriter, err := writer.CreateFormFile("resource", filepath.Base(filePath))
	if err != nil {
		return nil, "", c.errors.FailedToCreateFormFile(err)
	}

	if _, err := fileWriter.Write(fileData); err != nil {
		return nil, "", c.errors.FailedToWriteFileData(err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", c.errors.FailedToCloseWriter(err)
	}

	return body, writer.FormDataContentType(), nil
}

func (c *client) CreateFile(ctx context.Context, req CreateFileRequest) (*FileCreatedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	assetPath := fmt.Sprintf("new.%s", string(req.Type))
	templatePath := filepath.Join("templates", req.Language, assetPath)
	_, err := assets.Templates.ReadFile(templatePath)
	if err != nil {
		templatePath = filepath.Join("templates", "default", assetPath)
	}

	data := map[string]any{
		"title": fmt.Sprintf("%s.%s", req.Name, string(req.Type)),
	}

	body, contentType, err := c.createMultipartForm(data, templatePath)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("boards", req.BoardID, "documents")
	headers := map[string]string{
		"Content-Type": contentType,
	}

	var response FileCreatedResponse
	if err := c.sendRequest(ctx, http.MethodPost, url, req.Token, body, headers, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
