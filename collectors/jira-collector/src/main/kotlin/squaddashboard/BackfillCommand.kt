package squaddashboard

import squaddashboard.repository.SquadDashboardJiraIssueRepository
import squaddashboard.service.JiraIssueService

class BackfillCommand(private val jiraIssueService: JiraIssueService, private val squadDashboardJiraIssueRepository: SquadDashboardJiraIssueRepository) {

    fun run(projectKey: String) {
    }
}
