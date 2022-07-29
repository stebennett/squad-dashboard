package jirastatsservice

import "github.com/stebennett/squad-dashboard/pkg/jirarepository"

type JiraStatsService struct {
	JiraRepository jirarepository.JiraRepository
}

func (jss JiraStatsService) FetchThrougputDataForProject(project string) map[string]string {
	return map[string]string{
		"greeting": "hello",
	}
}
