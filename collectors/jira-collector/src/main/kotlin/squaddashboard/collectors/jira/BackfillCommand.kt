package squaddashboard.collectors.jira

import squaddashboard.collectors.jira.repository.SquadDashboardJiraIssueRepository
import squaddashboard.collectors.jira.service.JiraIssueService

class BackfillCommand(private val jiraIssueService: JiraIssueService, private val squadDashboardJiraIssueRepository: SquadDashboardJiraIssueRepository) {

    fun run(projectKey: String) {
    }
}
