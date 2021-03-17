package squaddashboard.collectors.jira

import io.kotest.core.spec.style.FunSpec
import io.mockk.mockk
import io.mockk.verifyOrder
import squaddashboard.collectors.jira.mapper.JiraIssueMapper
import squaddashboard.collectors.jira.model.IngestionType
import squaddashboard.collectors.jira.repository.SquadDashboardJiraIssueRepository
import squaddashboard.collectors.jira.service.JiraIssueService

@ExperimentalStdlibApi
class BackfillCommandTest : FunSpec({

    val jiraIssueService = mockk<JiraIssueService>(relaxed = true)
    val squadDashboardJiraIssueRepository = mockk<SquadDashboardJiraIssueRepository>(relaxed = true)
    val jiraIssueMapper = mockk<JiraIssueMapper>(relaxed = true)

    val backfillCommand = BackfillCommand(jiraIssueService, squadDashboardJiraIssueRepository, jiraIssueMapper)

    test("should create a project and enter start and stop for backfill") {

        val projectKey = "ABC"
        val workStartState = "workStartState"

        backfillCommand.run(projectKey, workStartState)

        verifyOrder {
            squadDashboardJiraIssueRepository.createProjectConfig(projectKey, workStartState)
            squadDashboardJiraIssueRepository.startIngestion(projectKey, IngestionType.Backfill, any())
            squadDashboardJiraIssueRepository.completeIngestion(projectKey, any())
        }
    }
})
