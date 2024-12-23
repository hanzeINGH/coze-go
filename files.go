package coze

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coze-dev/coze-go/internal"
)

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

// UploadFilesReq represents request for uploading file
type UploadFilesReq struct {
	// local file path
	filePath *string `json:"-"`

	// file byte array
	fileBytes []byte `json:"-"`

	// file name
	fileName *string `json:"-"`

	// file object
	file *os.File `json:"-"`
}

// NewUploadFilesReqWithBytes creates a new UploadFilesReq with file bytes
func NewUploadFilesReqWithBytes(fileName string, fileBytes []byte) *UploadFilesReq {
	return &UploadFilesReq{
		fileName:  internal.Ptr(fileName),
		fileBytes: fileBytes,
	}
}

// NewUploadFilesReqWithPath creates a new UploadFilesReq with file path
func NewUploadFilesReqWithPath(filePath string) *UploadFilesReq {
	return &UploadFilesReq{
		filePath: internal.Ptr(filePath),
	}
}

// NewUploadFilesReqWithFile creates a new UploadFilesReq with file path
func NewUploadFilesReqWithFile(file *os.File) *UploadFilesReq {
	return &UploadFilesReq{
		file: file,
	}
}

// RetrieveFilesReq represents request for retrieving file
type RetrieveFilesReq struct {
	FileID string `json:"file_id"`
}

// UploadFilesResp represents response for uploading file
type UploadFilesResp struct {
	internal.BaseResponse
	FileInfo *FileInfo `json:"data"`
}

// RetrieveFilesResp represents response for retrieving file
type RetrieveFilesResp struct {
	internal.BaseResponse
	FileInfo *FileInfo `json:"data"`
}

type files struct {
	client *internal.Client
}

func newFiles(client *internal.Client) *files {
	return &files{client: client}
}

func (r *files) Upload(ctx context.Context, req *UploadFilesReq) (*UploadFilesResp, error) {
	var fileSource any
	var filename string

	if req.filePath != nil {
		fileSource = &os.File{}
		filename = filepath.Base(*req.filePath)
		file, err := os.Open(*req.filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fileSource = file
	} else if len(req.fileBytes) > 0 {
		fileSource = req.fileBytes
		filename = internal.Value(req.fileName)
	} else if req.file != nil {
		fileSource = req.file
		filename = req.file.Name()
	} else {
		return nil, errors.New("file source is required")
	}

	return r.uploadFile(ctx, fileSource, filename)
}

// uploadFile 内部统一上传处理方法
func (r *files) uploadFile(ctx context.Context, fileSource any, filename string) (*UploadFilesResp, error) {
	path := "/v1/files/upload"
	var requestFile io.Reader

	switch v := fileSource.(type) {
	case *os.File:
		requestFile = v
	case []byte:
		requestFile = bytes.NewReader(v)
	}

	resp := &UploadFilesResp{}
	// 这里省略了 HTTP 请求的处理
	err := r.client.UploadFile(ctx, path, requestFile, filename, nil, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *files) Retrieve(ctx context.Context, req *RetrieveFilesReq) (*RetrieveFilesResp, error) {
	method := http.MethodPost
	uri := "/v1/files/retrieve"
	resp := &RetrieveFilesResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
