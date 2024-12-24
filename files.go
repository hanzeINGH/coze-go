package coze

import (
	"context"
	"io"
	"net/http"
)

type files struct {
	client *httpClient
}

func newFiles(client *httpClient) *files {
	return &files{client: client}
}

func (r *files) Upload(ctx context.Context, req fileInterface) (*UploadFilesResp, error) {
	path := "/v1/files/upload"
	resp := &uploadFilesResp{}
	err := r.client.UploadFile(ctx, path, req, req.Name(), nil, resp)
	if err != nil {
		return nil, err
	}

	resp.FileInfo.SetLogID(resp.LogID)
	return resp.FileInfo, nil
}

func (r *files) Retrieve(ctx context.Context, req *RetrieveFilesReq) (*RetrieveFilesResp, error) {
	method := http.MethodPost
	uri := "/v1/files/retrieve"
	resp := &retrieveFilesResp{}
	err := r.client.Request(ctx, method, uri, nil, resp, withHTTPQuery("file_id", req.FileID))
	if err != nil {
		return nil, err
	}
	resp.FileInfo.SetLogID(resp.LogID)
	return resp.FileInfo, nil
}

// FileInfo represents information about a file
type FileInfo struct {
	// The ID of the uploaded file.
	ID string `json:"id"`

	// The total byte size of the file.
	Bytes int `json:"bytes"`

	// The upload time of the file, in the format of a 10-digit Unix timestamp in seconds (s).
	CreatedAt int `json:"created_at"`

	// The name of the file.
	FileName string `json:"file_name"`
}

type fileInterface interface {
	io.Reader
	Name() string
}

type UploadFilesReq struct {
	io.Reader
	fileName string
}

func (r *UploadFilesReq) Name() string {
	return r.fileName
}

func NewUploadFileReq(reader io.Reader, fileName string) *UploadFilesReq {
	return &UploadFilesReq{
		fileName: fileName,
		Reader:   reader,
	}
}

// RetrieveFilesReq represents request for retrieving file
type RetrieveFilesReq struct {
	FileID string `json:"file_id"`
}

// uploadFilesResp represents response for uploading file
type uploadFilesResp struct {
	baseResponse
	FileInfo *UploadFilesResp `json:"data"`
}

// UploadFilesResp represents response for uploading file
type UploadFilesResp struct {
	baseModel
	FileInfo
}

// retrieveFilesResp represents response for retrieving file
type retrieveFilesResp struct {
	baseResponse
	FileInfo *RetrieveFilesResp `json:"data"`
}

// RetrieveFilesResp represents response for retrieving file
type RetrieveFilesResp struct {
	baseModel
	FileInfo
}
