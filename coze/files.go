package coze

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chyroc/go-ptr"
	"github.com/coze/coze/internal"
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

// UploadFileReq represents request for uploading file
type UploadFileReq struct {
	// local file path
	filePath *string `json:"-"`

	// file byte array
	fileBytes []byte `json:"-"`

	// file name
	fileName *string `json:"-"`

	// file object
	file *os.File `json:"-"`
}

// NewUploadFileReqWithBytes creates a new UploadFileReq with file bytes
func NewUploadFileReqWithBytes(fileName string, fileBytes []byte) *UploadFileReq {
	return &UploadFileReq{
		fileName:  ptr.Ptr(fileName),
		fileBytes: fileBytes,
	}
}

// NewUploadFileReqWithPath creates a new UploadFileReq with file path
func NewUploadFileReqWithPath(filePath string) *UploadFileReq {
	return &UploadFileReq{
		filePath: ptr.Ptr(filePath),
	}
}

// NewUploadFileReqWithFile creates a new UploadFileReq with file path
func NewUploadFileReqWithFile(file *os.File) *UploadFileReq {
	return &UploadFileReq{
		file: file,
	}
}

// RetrieveFileReq represents request for retrieving file
type RetrieveFileReq struct {
	FileID string `json:"file_id"`
}

// UploadFileResp represents response for uploading file
type UploadFileResp struct {
	internal.BaseResponse
	FileInfo *FileInfo `json:"data"`
}

// RetrieveFileResp represents response for retrieving file
type RetrieveFileResp struct {
	internal.BaseResponse
	FileInfo *FileInfo `json:"data"`
}

type files struct {
	client *internal.Client
}

func newFiles(client *internal.Client) *files {
	return &files{client: client}
}

func (r *files) Upload(ctx context.Context, req UploadFileReq) (*UploadFileResp, error) {
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
		filename = ptr.Value(req.fileName)
	} else if req.file != nil {
		fileSource = req.file
		filename = req.file.Name()
	} else {
		return nil, errors.New("file source is required")
	}

	return r.uploadFile(ctx, fileSource, filename)
}

// uploadFile 内部统一上传处理方法
func (r *files) uploadFile(ctx context.Context, fileSource any, filename string) (*UploadFileResp, error) {
	path := "/v1/files/upload"
	var requestFile io.Reader

	switch v := fileSource.(type) {
	case *os.File:
		requestFile = v
	case []byte:
		requestFile = bytes.NewReader(v)
	}

	resp := &UploadFileResp{}
	// 这里省略了 HTTP 请求的处理
	err := r.client.UploadFile(ctx, path, requestFile, filename, nil, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *files) Retrieve(ctx context.Context, req RetrieveFileReq) (*RetrieveFileResp, error) {
	method := http.MethodPost
	uri := "/v1/files/retrieve"
	resp := &RetrieveFileResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
