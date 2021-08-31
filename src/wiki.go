package backlog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GetAllWikis(refreshAll bool) ([]interface{}, error) {
	cachePath := CACHE_DIR + "/" + PROJECT_ID + "/wiki"
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil || refreshAll {
		url := "https://" + BACKLOG_BASE_URL + "/api/v2/wikis?projectIdOrKey=" + PROJECT_ID + "&apiKey=" + API_KEY
		response, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("URLが正しくありません")
		}
		byteArray, _ := ioutil.ReadAll(response.Body)
		var wiki interface{}
		err = json.Unmarshal(byteArray, &wiki)
		if err != nil {
			return nil, fmt.Errorf("想定外のJSON形式です")
		}
		wikis := wiki.([]interface{})
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(wikis)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var wikis interface{}
	_ = json.NewDecoder(file).Decode(&wikis)
	return wikis.([]interface{}), nil
}

func GetWikiContent(wikiId string, refreshAll bool) (map[string]interface{}, error) {
	cachePath := CACHE_DIR + "/" + PROJECT_ID + "/wiki-contents/" + wikiId
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil || refreshAll {
		url := "https://" + BACKLOG_BASE_URL + "/api/v2/wikis/" + wikiId + "?apiKey=" + API_KEY
		response, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("URLが正しくありません")
		}
		byteArray, _ := ioutil.ReadAll(response.Body)
		var wiki interface{}
		err = json.Unmarshal(byteArray, &wiki)
		if err != nil {
			return nil, fmt.Errorf("想定外のJSON形式です")
		}
		wikiContent := wiki.(map[string]interface{})
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(wikiContent)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var wikiContent interface{}
	_ = json.NewDecoder(file).Decode(&wikiContent)
	return wikiContent.(map[string]interface{}), nil
}
