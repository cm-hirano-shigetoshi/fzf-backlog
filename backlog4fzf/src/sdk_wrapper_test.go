package backlog4fzf

import (
	"testing"
)

func Test_getAllIssues(t *testing.T) {
	issues, _ := getAllIssues("41471")
	if len(issues) != 52 {
		t.Fatal(len(issues))
	}
}
