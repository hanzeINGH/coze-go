package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	filePath := os.Getenv("FILE_PATH")

	// upload file
	uploadResp, err := cozeCli.Files.Upload(ctx, coze.NewUploadFilesReqWithPath(filePath))
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	fileInfo := uploadResp.FileInfo

	// wait the server to process the file
	time.Sleep(time.Second)

	// retrieve file
	retrievedResp, err := cozeCli.Files.Retrieve(ctx, &coze.RetrieveFilesReq{
		FileID: fileInfo.ID,
	})
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		return
	}
	retrievedInfo := retrievedResp.FileInfo
	fmt.Println(retrievedInfo)
}
