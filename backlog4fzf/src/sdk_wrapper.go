package backlog4fzf

import (
	//"fmt"
	backlog "github.com/kenzo0107/backlog"
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
