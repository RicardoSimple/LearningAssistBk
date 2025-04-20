package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ar-app-api/conf"
	"ar-app-api/util/log"
)

const (
	GitHubAPIURL     = "https://api.github.com/repos/%s/%s/contents/%s"
	ImageURLPrefix   = "https://raw.githubusercontent.com/%s/%s/main/%s"
	ContentTypeJSON  = "application/json"
	AuthorizationKey = "Authorization"
	AcceptKey        = "Accept"
)

var uploader *GitHubFileUploader

type GitHubFileUploader struct {
	UserName   string
	Repo       string
	Token      string
	Message    string
	httpClient *http.Client
}

func GetUploader() *GitHubFileUploader {
	return uploader
}

func InitGitHubFileUploader(ctx context.Context) {
	// 初始化github文件上传
	log.Info("[INIT] github uploader init")
	git := conf.GetConfig().Git
	uploader = &GitHubFileUploader{
		UserName:   git.UserName,
		Repo:       git.Repo,
		Token:      git.Token,
		Message:    git.Message,
		httpClient: &http.Client{},
	}
	return
}

func (uploader *GitHubFileUploader) UploadFile(fileName, content string) (string, error) {
	// GitHub API URL
	apiURL := fmt.Sprintf(GitHubAPIURL, uploader.UserName, uploader.Repo, fileName)

	// Request headers
	headers := map[string]string{
		AuthorizationKey: "token " + uploader.Token,
		AcceptKey:        "application/vnd.github+json",
	}

	// Request body as JSON
	requestBody := map[string]interface{}{
		"message": uploader.Message + fileName,
		"content": content,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Build the request
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute the request
	resp, err := uploader.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	// 返回文件链接
	return uploader.HandleFileURL(fileName), nil
}

func (uploader *GitHubFileUploader) HandleFileURL(fileName string) string {
	return fmt.Sprintf(ImageURLPrefix, uploader.UserName, uploader.Repo, fileName)
}
