package configloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/stebennett/squad-dashboard/pkg/configrepository"
)

type Config struct {
	JiraProject          string
	JiraToDoStates       []string
	JiraInProgressStates []string
	JiraDoneStates       []string
	NonWorkingDays       []string
}

func Load(ctx context.Context, repo configrepository.ConfigRepository, config Config) error {
	_, err := repo.SaveJiraToDoStates(context.Background(), config.JiraProject, config.JiraToDoStates)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save Jira Work ToDo States. %s", err))
	}

	_, err = repo.SaveJiraInProgressStates(context.Background(), config.JiraProject, config.JiraInProgressStates)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save Jira Work In Progress States. %s", err))
	}

	_, err = repo.SaveJiraDoneStates(context.Background(), config.JiraProject, config.JiraDoneStates)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to save Jira Work Done States. %s", err))
	}

	_, err = repo.SaveNonWorkingDays(context.Background(), config.JiraProject, config.NonWorkingDays)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to save non working days. %s", err))
	}

	return nil
}
