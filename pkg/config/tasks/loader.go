package tasks

import (
	"context"
	"fmt"

	"github.com/stebennett/squad-dashboard/pkg/config/models"
	"github.com/stebennett/squad-dashboard/pkg/config/repo"
)

func Load(ctx context.Context, repo repo.ConfigRepository, config models.Config) error {
	_, err := repo.SaveJiraToDoStates(context.Background(), config.JiraProject, config.JiraToDoStates)
	if err != nil {
		return fmt.Errorf("failed to save Jira Work ToDo States. %s", err)
	}

	_, err = repo.SaveJiraInProgressStates(context.Background(), config.JiraProject, config.JiraInProgressStates)
	if err != nil {
		return fmt.Errorf("failed to save Jira Work In Progress States. %s", err)
	}

	_, err = repo.SaveJiraDoneStates(context.Background(), config.JiraProject, config.JiraDoneStates)
	if err != nil {
		return fmt.Errorf("failed to save Jira Work Done States. %s", err)
	}

	_, err = repo.SaveNonWorkingDays(context.Background(), config.JiraProject, config.NonWorkingDays)
	if err != nil {
		return fmt.Errorf("failed to save non working days. %s", err)
	}

	return nil
}
