package backlog4fzf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func printWikiList(withDesc bool, output string) (int, error) {
	for _, profile := range strings.Split(getProfiles(), ",") {
		backlogProfile, err := getBacklogProfile(profile, *appProfileConfig)
		if err != nil {
			return 1, err
		}
		cachePath := getWikiListCache(backlogProfile, *appCacheDir)
		wikis, err := getAllWikis(backlogProfile, cachePath)
		if err != nil {
			return 1, err
		}
		if output == "oneline" {
			for _, wiki := range wikis {
				fmt.Println(toOneLineWiki(profile, wiki.(map[string]interface{}), withDesc))
			}
		} else if output == "json" {
			_ = json.NewEncoder(os.Stdout).Encode(wikis)
		}
	}
	return 0, nil
}

func printWikiUrls(profile_wikis []string) (int, error) {
	urls := []string{}
	for _, profile_wiki := range profile_wikis {
		sp := strings.Split(profile_wiki, ":")
		backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
		if err != nil {
			return 1, err
		}
		urls = append(urls, "https://"+backlogProfile.baseUrl+"/alias/wiki/"+sp[1])
	}
	fmt.Println(strings.Join(urls, " "))
	return 0, nil
}

func getWikiListCache(profile BacklogProfile, cacheDir string) string {
	return cacheDir + "/" + profile.projectId + "/wiki-list"
}

func getAllWikis(profile BacklogProfile, cachePath string) ([]interface{}, error) {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil {
		allWikis, _ := getAllWikisSDK(profile)
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(allWikis)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var wikis interface{}
	_ = json.NewDecoder(file).Decode(&wikis)
	return wikis.([]interface{}), nil
}

func toOneLineWiki(profile string, wiki map[string]interface{}, withDesc bool) string {
	elem := []string{}
	elem = append(elem, profile)
	elem = append(elem, strconv.Itoa(int(wiki["id"].(float64))))
	elem = append(elem, wiki["name"].(string))
	return strings.Join(elem, ":")
}

func printWikiContent(profile_wiki string) (int, error) {
	sp := strings.Split(profile_wiki, ":")
	backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
	if err != nil {
		return 1, err
	}
	cachePath := getWikiContentCache(backlogProfile, *appCacheDir, sp[1])
	content, err := getWikiContent(backlogProfile, sp[1], cachePath)
	if err != nil {
		return 1, err
	}
	fmt.Println(content["content"])
	return 0, nil
}

func getWikiContent(profile BacklogProfile, wikiId string, cachePath string) (map[string]interface{}, error) {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil {
		allContent, _ := getWikiContentSDK(profile, wikiId)
		file, _ := os.Create(cachePath)
		defer file.Close()
		if allContent == nil {
			allContent = []interface{}{}
		}
		_ = json.NewEncoder(file).Encode(allContent)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var content interface{}
	_ = json.NewDecoder(file).Decode(&content)
	return content.(map[string]interface{}), nil
}

func getWikiContentCache(profile BacklogProfile, cacheDir string, wikiId string) string {
	return cacheDir + "/" + profile.projectId + "/wiki-content/" + wikiId
}

/*
func getAllWikis(refreshAll bool) ([]interface{}, error) {
	cachePath := getWikisCache()
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

func getWikiContent(wikiId string, refreshAll bool) (map[string]interface{}, error) {
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
*/
