package backlog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var BACKLOG_BASE_URL string
var PROJECT_ID string
var API_KEY string
var CACHE_DIR string

func SetEnv(backlogBaseUrl string, projectId string, apiKey string, cacheDir string) {
	BACKLOG_BASE_URL = backlogBaseUrl
	PROJECT_ID = projectId
	API_KEY = apiKey
	CACHE_DIR = cacheDir
}

func GetIssuesCache() string {
	return CACHE_DIR + "/" + PROJECT_ID + "/issue"
}

func GetAllIssues(refreshAll bool) ([]interface{}, error) {
	cachePath := GetIssuesCache()
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil || refreshAll {
		var allIssues []interface{}
		offset := 0
		for {
			url := "https://" + BACKLOG_BASE_URL + "/api/v2/issues?projectId[]=" + PROJECT_ID + "&count=100&offset=" + fmt.Sprint(offset) + "&apiKey=" + API_KEY
			response, err := http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("URLが正しくありません")
			}
			byteArray, _ := ioutil.ReadAll(response.Body)
			var issue interface{}
			err = json.Unmarshal(byteArray, &issue)
			if err != nil {
				return nil, fmt.Errorf("想定外のJSON形式です")
			}
			partIssues := issue.([]interface{})
			if len(partIssues) > 0 {
				allIssues = append(allIssues, partIssues...)
				offset += len(partIssues)
			} else {
				break
			}
		}
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(allIssues)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var issues interface{}
	_ = json.NewDecoder(file).Decode(&issues)
	return issues.([]interface{}), nil
}

func UpdateIssueStatus(issueId string, targetStatus string) error {
	targetStatusId := map[string]string{
		"MITAIOU":   "1",
		"TAIOUCHUU": "2",
		"SYORIZUMI": "3",
		"KANRYOU":   "4&resolutionId=0",
	}[targetStatus]
	url := "https://" + BACKLOG_BASE_URL + "/api/v2/issues/" + issueId + "?statusId=" + targetStatusId + "&apiKey=" + API_KEY
	req, _ := http.NewRequest(http.MethodPatch, url, nil)
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ステータスの更新に失敗しました")
	}
	return nil
}