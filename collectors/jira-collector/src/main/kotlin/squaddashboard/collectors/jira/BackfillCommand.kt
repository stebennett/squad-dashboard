package squaddashboard.collectors.jira

import squaddashboard.collectors.jira.mapper.JiraIssueMapper
import squaddashboard.collectors.jira.model.IngestionType
import squaddashboard.collectors.jira.repository.SquadDashboardJiraIssueRepository
import squaddashboard.collectors.jira.service.JiraIssueService
import java.time.Instant

class BackfillCommand(private val jiraIssueService: JiraIssueService,
                      private val squadDashboardJiraIssueRepository: SquadDashboardJiraIssueRepository,
                      private val jiraIssueMapper: JiraIssueMapper) {

    @ExperimentalStdlibApi
    fun run(projectKey: String, workStartState: String) {
        // create the project in the database - if it already exists, then fail as we can use this as our signal that the backfill has already happened
        squadDashboardJiraIssueRepository.createProjectConfig(projectKey, workStartState)
        squadDashboardJiraIssueRepository.startIngestion(projectKey, IngestionType.Backfill, Instant.now())

        // load all the issues into the db
        jiraIssueService.loadIssues(projectKey) { jiraIssue ->
            val mappedIssue = jiraIssueMapper.map(jiraIssue, projectKey)
            squadDashboardJiraIssueRepository.saveIssue(mappedIssue)
        }

        // persist the last run timestamp to the project table, with type backfill
        squadDashboardJiraIssueRepository.completeIngestion(projectKey, Instant.now())
    }
}
