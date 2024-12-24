package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const Version = "1.0.0"

type userAgentInfo struct {
	Version     string `json:"version"`
	Lang        string `json:"lang"`
	LangVersion string `json:"lang_version"`
	OsName      string `json:"os_name"`
	OsVersion   string `json:"os_version"`
}

func setUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", getUserAgent())

	// 添加 X-Coze-Client-User-Agent 头
	clientUA, err := getCozeClientUserAgent()
	if err == nil {
		req.Header.Set("X-Coze-Client-User-Agent", clientUA)
	}
}

func getOsVersion() string {
	return runtime.GOOS + "/" + os.Getenv("OSVERSION")
}

func getUserAgent() string {
	return fmt.Sprintf(
		"cozego/%s go/%s %s",
		Version,
		strings.TrimPrefix(runtime.Version(), "go"),
		getOsVersion(),
	)
}

func getCozeClientUserAgent() (string, error) {
	ua := userAgentInfo{
		Version:     Version,
		Lang:        "go",
		LangVersion: strings.TrimPrefix(runtime.Version(), "go"),
		OsName:      runtime.GOOS,
		OsVersion:   os.Getenv("OSVERSION"),
	}

	data, err := json.Marshal(ua)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
