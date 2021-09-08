package backlog4fzf

import (
	backlog "github.com/kenzo0107/backlog"
	"strconv"
)

func getAllIssuesSDK(projectIdStr string) ([]interface{}, error) {
	projectId, _ := strconv.Atoi(projectIdStr)
	client := backlog.New(API_KEY, "https://"+BACKLOG_BASE_URL)
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
