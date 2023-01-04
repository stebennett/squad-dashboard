package dashboard

import (
	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

func GenerateEscapedDefects(weekCount int, repo jirarepository.JiraRepository) ([]models.EscapedDefectCount, error) {
	return []models.EscapedDefectCount{}, nil
}
