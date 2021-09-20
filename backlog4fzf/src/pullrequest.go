package backlog4fzf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printPullrequestList(withDesc bool) (int, error) {
	for _, profile := range strings.Split(getProfiles(), ",") {
		backlogProfile, err := getBacklogProfile(profile, *appProfileConfig)
		if err != nil {
			return 1, err
		}
		cachePath := getPullrequestListCache(backlogProfile, *appCacheDir)
		var pullrequests []interface{}
		pullrequests, err = getAllPullrequests(backlogProfile, cachePath)
		if err != nil {
			return 1, err
		}
		for _, pullrequest := range pullrequests {
			oneLine := toOneLinePullrequest(profile, pullrequest.(map[string]interface{}), withDesc)
			if len(oneLine) > 0 {
				fmt.Println(oneLine)
			}
		}
	}
	return 0, nil
}

func printPullrequestUrls(profile_pullrequests []string) (int, error) {
	urls := []string{}
	for _, profile_pullrequest := range profile_pullrequests {
		sp := strings.Split(profile_pullrequest, ":")
		backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
		if err != nil {
			return 1, err
		}
		urls = append(urls, "https://"+backlogProfile.baseUrl+"/git/"+backlogProfile.projectId+"/"+sp[1]+"/pullRequests/"+sp[2][1:])
	}
	fmt.Println(strings.Join(urls, " "))
	return 0, nil
}

func deletePullrequestCache(profiles []string) (int, error) {
	uniq := map[string]bool{}
	for _, profile := range profiles {
		uniq[profile] = true
	}
	for profile, _ := range uniq {
		backlogProfile, err := getBacklogProfile(profile, *appProfileConfig)
		cachePath := getPullrequestListCache(backlogProfile, *appCacheDir)
		if err != nil {
			return 1, err
		}
		os.Remove(cachePath)
	}
	return 0, nil
}

func getAllPullrequests(profile BacklogProfile, cachePath string) ([]interface{}, error) {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil {
		var pullrequests []map[string]interface{}
		pullrequests, _ = getAllPullrequestsSDK(profile)
		file, _ := os.Create(cachePath)
		defer file.Close()
		_ = json.NewEncoder(file).Encode(pullrequests)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var pullrequests []interface{}
	_ = json.NewDecoder(file).Decode(&pullrequests)
	return pullrequests, nil
}

func getPullrequestListCache(profile BacklogProfile, cacheDir string) string {
	return cacheDir + "/" + profile.projectId + "/pullrequest-list"
}

func colorPullrequestStatus(status string) string {
	if status == "Open" {
		return "\033[31m" + status + "\033[0m"
	} else if status == "Merged" {
		return "\033[32m" + status + "\033[0m"
	} else if status == "Closed" {
		return "\033[90m" + status + "\033[0m"
	}
	return status
}

func toOneLinePullrequest(profile string, pullrequest map[string]interface{}, withDesc bool) string {
	if x, ok := pullrequest["pullRequests"]; !ok || x == nil {
		return ""
	}
	lines := []string{}
	repositoryName := pullrequest["repositoryName"].(string)
	for _, p := range pullrequest["pullRequests"].([]interface{}) {
		pullreq := p.(map[string]interface{})
		elem := []string{}
		elem = append(elem, profile)
		elem = append(elem, repositoryName)
		number := pullreq["number"]
		elem = append(elem, "#"+fmt.Sprintf("%v", number))
		elem = append(elem, colorPullrequestStatus(pullreq["status"].(map[string]interface{})["name"].(string)))
		if assignee := pullreq["assignee"]; assignee != nil {
			elem = append(elem, assignee.(map[string]interface{})["name"].(string))
		} else {
			elem = append(elem, "")
		}
		elem = append(elem, pullreq["summary"].(string))
		if withDesc {
			space := ""
			for i := 0; i <= 100; i++ {
				space += "          "
			}
			elem[len(elem)-1] += space + "\\n"
			desc := pullreq["description"].(string)
			elem = append(elem, strings.Replace(desc, "\n", "\\n", -1))
		}
		line := strings.Join(elem, ":")
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
