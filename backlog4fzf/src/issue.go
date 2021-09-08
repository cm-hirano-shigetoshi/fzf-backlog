package backlog4fzf

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func getIssuesCache() string {
	return CACHE_DIR + "/" + PROJECT_ID + "/issue"
}

func getAllIssues(refreshAll bool) ([]interface{}, error) {
	cachePath := getIssuesCache()
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil || refreshAll {
		allIssues, _ := getAllIssuesSDK(PROJECT_ID)
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

func updateIssueStatus(profile string, issueId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool, targetStatus string) error {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
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
