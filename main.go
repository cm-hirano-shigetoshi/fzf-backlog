package main

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	flag "github.com/spf13/pflag"
	"log"
	"main/src"
	"os"
	"strings"
)

func getProfileList(profilePath string) ([]string, error) {
	tree, err := toml.LoadFile(profilePath)
	if err != nil {
		return []string{}, fmt.Errorf("profileファイルの読み込みに失敗")
	}
	profileList := []string{}
	for _, a := range tree.Keys() {
		profileList = append(profileList, a)
	}
	return profileList, nil
}

func getBacklogProfile(profile string, profilePath string) (string, string, string, error) {
	tree, err := toml.LoadFile(profilePath)
	if err != nil {
		return "", "", "", fmt.Errorf("profileファイルの読み込みに失敗")
	}
	if !tree.Has(profile) {
		return "", "", "", fmt.Errorf("profile: " + profile + " がありません")
	}
	backlogBaseUrl := tree.GetDefault(profile+".base_url", "").(string)
	projectId := tree.GetDefault(profile+".project_id", "").(string)
	apiKey := tree.GetDefault(profile+".api_key", "").(string)
	if len(backlogBaseUrl)*len(projectId)*len(apiKey) == 0 {
		return "", "", "", fmt.Errorf("profileに必要な設定がありません")
	}
	return backlogBaseUrl, projectId, apiKey, nil
}

func printIssues(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	backlog.SetEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	allIssues, err := backlog.GetAllIssues(refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["issues"] = allIssues
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func printPullrequests(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	backlog.SetEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	allPullrequests, err := backlog.GetAllPullrequests(refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["repositories"] = allPullrequests
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func printWikis(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	backlog.SetEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	allWikis, err := backlog.GetAllWikis(refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["wikis"] = allWikis
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func printWiki(profile string, wikiId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	backlog.SetEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	wikiContent, err := backlog.GetWikiContent(wikiId, refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["content"] = wikiContent
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func deleteWikiContentCache(profile string, wikiId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	cachePath := cacheDir + "/" + projectId + "/wiki-contents/" + wikiId
	fmt.Println("Remove " + cachePath)
	if err := os.Remove(cachePath); err != nil {
		fmt.Println(err)
	}
}

func updateIssueStatus(profile string, issueId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool, targetStatus string) {
	backlog.SetEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	err := backlog.UpdateIssueStatus(issueId, targetStatus)
	if err != nil {
		log.Fatal("ステータスの更新に失敗しました")
	}
}

func main() {
	profilesPtr := flag.String("profiles", "", "comma separated profiles")
	refreshAllPtrPtr := flag.Bool("refresh-all", false, "reload forcely")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("引数がありません")
	}

	cacheDir := ""
	if val, ok := os.LookupEnv("FZF_BACKLOG_CACHE_DIR"); ok {
		cacheDir = val
	} else {
		cacheDir = os.Getenv("HOME") + "/.backlog/cache"
	}

	if flag.Args()[0] == "issues" {
		profiles := *profilesPtr
		if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
			profiles = val
		}

		if len(profiles) == 0 {
			return
		}
		for _, profile := range strings.Split(profiles, ",") {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			printIssues(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
		}
	} else if flag.Args()[0] == "issues-across-projects" {
		profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
		if err != nil {
			log.Fatal(err)
		}
		for _, profile := range profileList {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			printIssues(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.Args()[0] == "issue-urls" {
		urls := []string{}
		for _, issue := range flag.Args()[1:] {
			sp := strings.Split(issue, ":")
			profile := sp[0]
			issueId := sp[1]
			backlogBaseUrl, _, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			urls = append(urls, "https://"+backlogBaseUrl+"/view/"+issueId)
		}
		fmt.Println(strings.Join(urls, " "))
	} else if flag.Args()[0] == "pullrequests" {
		profiles := *profilesPtr
		if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
			profiles = val
		}

		if len(profiles) == 0 {
			return
		}
		for _, profile := range strings.Split(profiles, ",") {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			printPullrequests(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
		}
	} else if flag.Args()[0] == "pullrequests-across-projects" {
		profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
		if err != nil {
			log.Fatal(err)
		}
		for _, profile := range profileList {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			printPullrequests(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.Args()[0] == "pullrequest-urls" {
		urls := []string{}
		for _, issue := range flag.Args()[1:] {
			sp := strings.Split(issue, ":")
			profile := sp[0]
			repository := sp[1]
			number := sp[2][1:]
			backlogBaseUrl, projectId, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			urls = append(urls, "https://"+backlogBaseUrl+"/git/"+projectId+"/"+repository+"/pullRequests/"+number)
		}
		fmt.Println(strings.Join(urls, " "))
	} else if flag.Args()[0] == "wikis" {
		profiles := *profilesPtr
		if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
			profiles = val
		}

		if len(profiles) == 0 {
			return
		}
		for _, profile := range strings.Split(profiles, ",") {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			printWikis(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
		}
	} else if flag.Args()[0] == "wikis-across-projects" {
		profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
		if err != nil {
			log.Fatal(err)
		}
		for _, profile := range profileList {
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			printWikis(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.Args()[0] == "wiki-content" {
		sp := strings.Split(flag.Args()[1], ":")
		profile := sp[0]
		wikiId := sp[1]
		backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
		if err != nil {
			log.Fatal(err)
		}
		printWiki(profile, wikiId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
	} else if flag.Args()[0] == "wiki-urls" {
		urls := []string{}
		for _, wiki := range flag.Args()[1:] {
			sp := strings.Split(wiki, ":")
			profile := sp[0]
			id := sp[1]
			backlogBaseUrl, _, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			urls = append(urls, "https://"+backlogBaseUrl+"/alias/wiki/"+id)
		}
		fmt.Println(strings.Join(urls, " "))
	} else if flag.Args()[0] == "delte-wiki-content-cache" {
		sp := strings.Split(flag.Args()[1], ":")
		profile := sp[0]
		wikiId := sp[1]
		backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
		if err != nil {
			log.Fatal(err)
		}
		deleteWikiContentCache(profile, wikiId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
	} else if flag.Args()[0] == "update-issue-status" {
		targetStatus := flag.Args()[1]
		for _, issue := range flag.Args()[2:] {
			sp := strings.Split(issue, ":")
			profile := sp[0]
			issueId := sp[1]
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			updateIssueStatus(profile, issueId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr, targetStatus)
		}
	}
}
