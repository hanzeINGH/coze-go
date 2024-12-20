package coze

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/coze-dev/coze-go/coze/internal"
)

type streamable interface {
	ChatEvent | WorkflowEvent
}

// type eventProcessor[T streamable] interface {
//	ProcessLine(line []byte, reader *bufio.Reader, logID string) (T, bool, error)
// }

type eventProcessor[T streamable] func(line []byte, reader *bufio.Reader, logID string) (*T, bool, error)

type streamReader[T streamable] struct {
	isFinished bool

	reader    *bufio.Reader
	response  *http.Response
	logID     string
	processor eventProcessor[T]
}

func (s *streamReader[T]) Recv() (response *T, err error) {
	return s.processLines()
}

//nolint:gocognit
func (s *streamReader[T]) processLines() (*T, error) {
	err := s.checkRespErr()
	if err != nil {
		return nil, err
	}
	for {
		line, _, readErr := s.reader.ReadLine()
		if readErr != nil {
			return nil, readErr
		}

		if line == nil {
			s.isFinished = true
			break
		}
		if len(line) == 0 {
			continue
		}
		event, isDone, err := s.processor(line, s.reader, s.logID)
		if err != nil {
			return nil, err
		}
		s.isFinished = isDone
		if event == nil {
			continue
		}
		return event, nil
	}
	return nil, io.EOF
}

func (s *streamReader[T]) checkRespErr() error {
	contentType := s.response.Header.Get("Content-Type")
	if contentType != "" && strings.Contains(contentType, "application/json") {
		respStr, err := ioutil.ReadAll(s.response.Body)
		if err != nil {
			// logger.Warn("Error reading response body: ", err)
			return err
		}
		var baseResp internal.BaseResponse
		if err := json.Unmarshal(respStr, &baseResp); err != nil {
			// logger.Warn("Error unmarshalling response: ", err)
			return err
		}
		if baseResp.Code != 0 {
			// logger.Warn("API error: %d %s", baseResp.Code, baseResp.Msg)
			// todo
			return errors.New(fmt.Sprintf("API error: %d %s", baseResp.Code, baseResp.Msg))
		}
		return nil
	}
	return nil
}

func (s *streamReader[T]) Close() error {
	return s.response.Body.Close()
}
