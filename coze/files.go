package coze

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
	FilePath string `json:"-"`

	// file byte array
	FileBytes []byte `json:"-"`

	// file name
	FileName string `json:"-"`

	// file object
	File any `json:"-"` // In Go we'll handle this differently
}

// NewUploadFileReqWithBytes creates a new UploadFileReq with file bytes
func NewUploadFileReqWithBytes(fileName string, fileBytes []byte) *UploadFileReq {
	return &UploadFileReq{
		FileName:  fileName,
		FileBytes: fileBytes,
	}
}

// NewUploadFileReqWithPath creates a new UploadFileReq with file path
func NewUploadFileReqWithPath(filePath string) *UploadFileReq {
	return &UploadFileReq{
		FilePath: filePath,
	}
}

// RetrieveFileReq represents request for retrieving file
type RetrieveFileReq struct {
	FileID string `json:"file_id"`
}

// NewRetrieveFileReq creates a new RetrieveFileReq
func NewRetrieveFileReq(fileID string) *RetrieveFileReq {
	return &RetrieveFileReq{
		FileID: fileID,
	}
}

// UploadFileResp represents response for uploading file
type UploadFileResp struct {
	FileInfo *FileInfo `json:"file_info"`
}

// RetrieveFileResp represents response for retrieving file
type RetrieveFileResp struct {
	FileInfo *FileInfo `json:"file_info"`
}

type files struct {
}
