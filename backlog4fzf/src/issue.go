package backlog4fzf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printIssueList(withDesc bool, output string) (int, error) {
	for _, profile := range strings.Split(getProfiles(), ",") {
		backlogProfile, err := getBacklogProfile(profile, *appProfileConfig)
		if err != nil {
			return 1, err
		}
		cachePath := getIssueListCache(backlogProfile, *appCacheDir)
		issues, err := getAllIssues(backlogProfile, cachePath)
		if err != nil {
			return 1, err
		}
		if output == "oneline" {
			for _, issue := range issues {
				fmt.Println(toOneLineIssue(profile, issue.(map[string]interface{}), withDesc))
			}
		} else if output == "json" {
			_ = json.NewEncoder(os.Stdout).Encode(issues)
		}
	}
	return 0, nil
}

func printIssueUrls(profile_issues []string) (int, error) {
	urls := []string{}
	for _, profile_issue := range profile_issues {
		sp := strings.Split(profile_issue, ":")
		backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
		if err != nil {
			return 1, err
		}
		urls = append(urls, "https://"+backlogProfile.baseUrl+"/view/"+sp[1])
	}
	fmt.Println(strings.Join(urls, " "))
	return 0, nil
}

func deleteIssueCache(profiles []string) (int, error) {
	uniq := map[string]bool{}
	for _, profile := range profiles {
		uniq[profile] = true
	}
	for profile, _ := range uniq {
		backlogProfile, err := getBacklogProfile(profile, *appProfileConfig)
		cachePath := getIssueListCache(backlogProfile, *appCacheDir)
		if err != nil {
			return 1, err
		}
		os.Remove(cachePath)
	}
	return 0, nil
}

func updateIssueStatus(target string, profile_issues []string) (int, error) {
	for _, profile_issue := range profile_issues {
		sp := strings.Split(profile_issue, ":")
		backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
		if err != nil {
			return 1, err
		}
		updateIssueStatusSDK(backlogProfile, sp[1], target)
	}
	return 0, nil
}

func printIssueDescription(profile_issue string) (int, error) {
	sp := strings.Split(profile_issue, ":")
	backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
	if err != nil {
		return 1, err
	}
	cachePath := getIssueListCache(backlogProfile, *appCacheDir)
	desc, err := getIssueDescription(backlogProfile, sp[1], cachePath)
	if err != nil {
		return 1, err
	}
	fmt.Println(desc)
	return 0, nil
}

func printIssueComments(profile_issue string) (int, error) {
	sp := strings.Split(profile_issue, ":")
	backlogProfile, err := getBacklogProfile(sp[0], *appProfileConfig)
	if err != nil {
		return 1, err
	}
	cachePath := getIssueCommentsCache(backlogProfile, *appCacheDir, sp[1])
	comments, err := getIssueComments(backlogProfile, sp[1], cachePath)
	if err != nil {
		return 1, err
	}
	contents := []string{}
	for _, comment := range comments {
		if content, ok := comment.(map[string]interface{})["content"]; ok && content != nil {
			contents = append(contents, content.(string))
		}
	}
	fmt.Println(strings.Join(contents, "\n--\n"))
	return 0, nil
}

func printIssueDescriptionAndComments(profile_issue string) (int, error) {
	exit, err := printIssueDescription(profile_issue)
	if err != nil {
		return exit, err
	}
	fmt.Println("==\n")
	exit, err = printIssueComments(profile_issue)
	if err != nil {
		return exit, err
	}
	return 0, nil
}

func colorIssueStatus(status string) string {
	if status == "未対応" {
		return "\033[31m" + status + "\033[0m"
	} else if status == "処理中" {
		return "\033[34m" + status + "\033[0m"
	} else if status == "処理済み" {
		return "\033[32m" + status + "\033[0m"
	} else if status == "完了" {
		return "\033[33m" + status + "\033[0m"
	}
	return status
}

func toOneLineIssue(profile string, issue map[string]interface{}, withDesc bool) string {
	elem := []string{}
	elem = append(elem, profile)
	elem = append(elem, issue["issueKey"].(string))
	elem = append(elem, colorIssueStatus(issue["status"].(map[string]interface{})["name"].(string)))
	if assignee := issue["assignee"]; assignee != nil {
		elem = append(elem, assignee.(map[string]interface{})["name"].(string))
	} else {
		elem = append(elem, "")
	}
	elem = append(elem, issue["summary"].(string))
	if withDesc {
		space := ""
		for i := 0; i <= 100; i++ {
			space += "          "
		}
		elem[len(elem)-1] += space + "\\n"
		desc := issue["description"].(string)
		elem = append(elem, strings.Replace(desc, "\n", "\\n", -1))
	}
	return strings.Join(elem, ":")
}

func getAllIssues(profile BacklogProfile, cachePath string) ([]interface{}, error) {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil {
		allIssues, _ := getAllIssuesSDK(profile)
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

func getIssueDescription(profile BacklogProfile, issueId string, cachePath string) (string, error) {
	issues, err := getAllIssues(profile, cachePath)
	if err != nil {
		return "", err
	}
	for _, issue := range issues {
		if issue.(map[string]interface{})["issueKey"].(string) == issueId {
			return issue.(map[string]interface{})["description"].(string), nil
		}
	}
	return "", nil
}

func getIssueComments(profile BacklogProfile, issueId string, cachePath string) ([]interface{}, error) {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	if _, err := os.Stat(cachePath); err != nil {
		allComments, _ := getIssueCommentsSDK(profile, issueId)
		file, _ := os.Create(cachePath)
		defer file.Close()
		if allComments == nil {
			allComments = []interface{}{}
		}
		_ = json.NewEncoder(file).Encode(allComments)
	}
	file, _ := os.Open(cachePath)
	defer file.Close()
	var comments interface{}
	_ = json.NewDecoder(file).Decode(&comments)
	return comments.([]interface{}), nil
}

func getIssueListCache(profile BacklogProfile, cacheDir string) string {
	return cacheDir + "/" + profile.projectId + "/issue-list"
}

func getIssueCommentsCache(profile BacklogProfile, cacheDir string, issueId string) string {
	return cacheDir + "/" + profile.projectId + "/issue-comments/" + issueId
}
