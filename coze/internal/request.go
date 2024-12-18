package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coze/coze/auth_error"
)

// Doer 是一个执行 HTTP 请求的接口
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client HTTP 客户端封装
type Client struct {
	doer    Doer
	baseURL string
}

// NewClient 创建新的 HTTP 客户端
func NewClient(doer Doer, baseURL string) *Client {
	if doer == nil {
		doer = &http.Client{}
	}
	return &Client{
		doer:    doer,
		baseURL: baseURL,
	}
}

// RequestOption 请求选项函数类型
type RequestOption func(*http.Request) error

// WithHeader 添加请求头
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// WithQuery 添加查询参数
func WithQuery(key, value string) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

// Request 发送请求
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, instance interface{}, opts ...RequestOption) error {
	urlInfo := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlInfo, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// 设置默认请求头
	req.Header.Set("Content-Type", "application/json")

	// 应用请求选项
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return fmt.Errorf("apply option: %w", err)
		}
	}

	resp, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	// 鉴权的情况，需要解析
	if resp.StatusCode != http.StatusOK {
		errorInfo := auth_error.CozeError{}
		err = json.Unmarshal(bodyBytes, &errorInfo)
		if err != nil {
			// todo log
			return fmt.Errorf("unmarshal error response: %w", err)
		}
		return auth_error.NewCozeAuthExceptionWithoutParent(&errorInfo, resp.StatusCode, getLogID(resp.Header))
	}

	return json.Unmarshal(bodyBytes, instance)
}

//// Stream 处理流式响应
//func (c *Client) Stream(ctx context.Context, method, path string, body interface{}, handler func([]byte) error, opts ...RequestOption) error {
//	resp, err := c.Request(ctx, method, path, nil, body, opts...)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	reader := bufio.NewReader(resp.Body)
//	for {
//		line, err := reader.ReadBytes('\n')
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			return fmt.Errorf("read stream: %w", err)
//		}
//
//		if err := handler(line); err != nil {
//			return fmt.Errorf("handle stream data: %w", err)
//		}
//	}
//
//	return nil
//}

// UploadFile 上传文件
func (c *Client) UploadFile(ctx context.Context, path string, files map[string]string, fields map[string]string, opts ...RequestOption) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	for fieldName, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("open file %s: %w", filePath, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			return nil, fmt.Errorf("create form file: %w", err)
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("copy file content: %w", err)
		}
	}

	// 添加其他字段
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("write field %s: %w", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, path), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 应用请求选项
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

// 使用示例：
type MockDoer struct {
	Response *http.Response
	Error    error
}

func (m *MockDoer) Do(*http.Request) (*http.Response, error) {
	return m.Response, m.Error
}
