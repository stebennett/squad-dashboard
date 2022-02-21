package adapters

import (
	"testing"

	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stretchr/testify/assert"
)

func TestAdaptIssue(t *testing.T) {
	inputIssue := jiraservice.JiraIssue{}

	outputIssue := AdaptIssue(inputIssue)

	assert.Equal(t, inputIssue.Key, outputIssue.Key)
	assert.Equal(t, inputIssue.Fields.IssueType.Name, outputIssue.IssueType)
}
