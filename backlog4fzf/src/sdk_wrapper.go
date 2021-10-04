package backlog4fzf

import (
	"encoding/json"
	"fmt"
	backlog "github.com/kenzo0107/backlog"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getAllIssuesSDK(profile BacklogProfile) ([]interface{}, error) {
	projectId, _ := strconv.Atoi(profile.projectId)
	client := backlog.New(profile.apiKey, "https://"+profile.baseUrl)
	var allIssues []interface{}
	offset := 0
	for {
		opts := &backlog.GetIssuesOptions{
			ProjectIDs: []int{projectId},
			Count:      backlog.Int(100),
			Offset:     backlog.Int(offset),
		}
		issues, err := client.GetIssues(opts)
		if err != nil {
			return nil, err
		}
		if len(issues) == 0 {
			break
		}
		offset += len(issues)
		bbb := []interface{}{}
		for _, issue := range issues {
			bbb = append(bbb, *issue)
		}
		allIssues = append(allIssues, bbb...)
	}
	return allIssues, nil
}

func updateIssueStatusSDK(profile BacklogProfile, issueId string, target string) error {
	targetId := map[string][]int{
		"MITAIOU":   []int{1},
		"TAIOUCHUU": []int{2},
		"SYORIZUMI": []int{3},
		"KANRYOU":   []int{4, 0},
	}[target]

	client := backlog.New(profile.apiKey, "https://"+profile.baseUrl)
	updateIssueInput := &backlog.UpdateIssueInput{
		StatusID: backlog.Int(targetId[0]),
	}
	if len(targetId) > 1 {
		updateIssueInput.ResolutionID = backlog.Int(targetId[1])
	}
	client.UpdateIssue(issueId, updateIssueInput)
	return nil
}

func getIssueCommentsSDK(profile BacklogProfile, issueId string) ([]interface{}, error) {
	client := backlog.New(profile.apiKey, "https://"+profile.baseUrl)
	var allComments []interface{}
	comments, err := client.GetIssueComments(issueId, &backlog.GetIssueCommentsOptions{
		Count: backlog.Int(100),
	})
	if err != nil {
		return nil, err
	}
	bbb := []interface{}{}
	for _, comment := range comments {
		bbb = append(bbb, *comment)
	}
	allComments = append(allComments, bbb...)
	return allComments, nil
}

func getRepositories(profile BacklogProfile) (map[string]string, error) {
	url := "https://" + profile.baseUrl + "/api/v2/projects/" + profile.projectId + "/git/repositories?apiKey=" + profile.apiKey
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

func getAllPullrequestsOfRepository(profile BacklogProfile, repoId string) ([]interface{}, error) {
	var pullrequests []interface{}
	offset := 0
	for {
		url := "https://" + profile.baseUrl + "/api/v2/projects/" + profile.projectId + "/git/repositories/" + repoId + "/pullRequests?&count=100&offset=" + strconv.Itoa(offset) + "&apiKey=" + profile.apiKey
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
			pullrequests = append(pullrequests, partPullrequests...)
			offset += len(partPullrequests)
		} else {
			break
		}
	}
	return pullrequests, nil
}

func getAllPullrequestsSDK(profile BacklogProfile) ([]map[string]interface{}, error) {
	allPullrequests := []map[string]interface{}{}
	repositories, _ := getRepositories(profile)
	for repoId, repoName := range repositories {
		var pullrequests []interface{}
		pullrequests, err := getAllPullrequestsOfRepository(profile, repoId)
		if err != nil {
			return nil, err
		}
		jsonObj := map[string]interface{}{}
		jsonObj["repositoryId"] = repoId
		jsonObj["repositoryName"] = repoName
		jsonObj["pullRequests"] = pullrequests
		allPullrequests = append(allPullrequests, jsonObj)
	}
	return allPullrequests, nil
}

func getAllWikisSDK(profile BacklogProfile) (interface{}, error) {
	projectId, _ := strconv.Atoi(profile.projectId)
	client := backlog.New(profile.apiKey, "https://"+profile.baseUrl)
	opts := &backlog.GetWikisOptions{
		ProjectIDOrKey: projectId,
	}
	wikis, err := client.GetWikis(opts)
	if err != nil {
		return nil, err
	}
	return wikis, nil
}

func getWikiContentSDK(profile BacklogProfile, wikiId string) (interface{}, error) {
	client := backlog.New(profile.apiKey, "https://"+profile.baseUrl)
	id, _ := strconv.Atoi(wikiId)
	content, err := client.GetWiki(id)
	if err != nil {
		return nil, err
	}
	return content, nil
}
