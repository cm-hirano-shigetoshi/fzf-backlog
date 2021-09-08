package backlog4fzf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func getRepositories() (map[string]string, error) {
	url := "https://" + BACKLOG_BASE_URL + "/api/v2/projects/" + PROJECT_ID + "/git/repositories?apiKey=" + API_KEY
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("URLが正しくありません")
	}
	byteArray, _ := ioutil.ReadAll(response.Body)
	var repositories interface{}
	err = json.Unmarshal(byteArray, &repositories)
	if err != nil {
		return nil, fmt.Errorf("想定外のJSON形式です")
	}
	repos := map[string]string{}
	for _, repo := range repositories.([]interface{}) {
		r := repo.(map[string]interface{})
		repos[strconv.Itoa(int(r["id"].(float64)))] = r["name"].(string)
	}
	return repos, nil
}

func getAllPullrequests(refreshAll bool) ([]interface{}, error) {
	cachePath := CACHE_DIR + "/" + PROJECT_ID + "/pullrequest"
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil || refreshAll {
		var allPullrequests []interface{}
		repositories, _ := getRepositories()
		for repoId, repoName := range repositories {
			var repoPullrequests []interface{}
			offset := 0
			for {
				url := "https://" + BACKLOG_BASE_URL + "/api/v2/projects/" + PROJECT_ID + "/git/repositories/" + repoId + "/pullRequests?&count=100&offset=" + strconv.Itoa(offset) + "&apiKey=" + API_KEY
				response, err := http.Get(url)
				if err != nil {
					return nil, fmt.Errorf("URLが正しくありません")
				}
				byteArray, _ := ioutil.ReadAll(response.Body)
				var pullrequest interface{}
				err = json.Unmarshal(byteArray, &pullrequest)
				if err != nil {
					return nil, fmt.Errorf("想定外のJSON形式です")
				}
				partPullrequests := pullrequest.([]interface{})
				if len(partPullrequests) > 0 {
					repoPullrequests = append(repoPullrequests, partPullrequests...)
					offset += len(partPullrequests)
				} else {
					break
				}
			}
			jsonObj := map[string]interface{}{}
			jsonObj["repositoryId"] = repoId
			jsonObj["repositoryName"] = repoName
			jsonObj["pullRequests"] = repoPullrequests
			allPullrequests = append(allPullrequests, jsonObj)
		}
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(allPullrequests)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var pullrequests interface{}
	_ = json.NewDecoder(file).Decode(&pullrequests)
	if pullrequests == nil {
		return nil, nil
	} else {
		return pullrequests.([]interface{}), nil
	}
}
