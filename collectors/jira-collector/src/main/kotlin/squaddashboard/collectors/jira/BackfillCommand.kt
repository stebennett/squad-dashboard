package squaddashboard.collectors.jira

import squaddashboard.collectors.jira.repository.SquadDashboardJiraIssueRepository
import squaddashboard.collectors.jira.service.JiraIssueService

class BackfillCommand(private val jiraIssueService: JiraIssueService, private val squadDashboardJiraIssueRepository: SquadDashboardJiraIssueRepository) {

    fun run(projectKey: String) {
        // create the project in the database - if it already exists, then fail as we can use this as our signal that the backfill has already happened

        // load all the issues into the db

        // persist the last run timestamp to the project table, with type backfill
    }
}
