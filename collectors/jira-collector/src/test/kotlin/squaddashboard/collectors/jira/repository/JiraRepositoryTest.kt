package squaddashboard.collectors.jira.repository

import io.kotest.core.spec.style.FunSpec
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import squaddashboard.collectors.jira.client.JiraClient
import squaddashboard.collectors.jira.client.model.IssueSearchCommand
import squaddashboard.collectors.jira.client.model.JiraCommandFactory

class JiraRepositoryTest : FunSpec({

    val jiraClient = mockk<JiraClient>(relaxed = true)
    val jiraCommandFactory = mockk<JiraCommandFactory>()
    val jiraRepository = JiraRepository(jiraClient, jiraCommandFactory)

    test("should make request to client") {
        val projectKey = "ABC"
        val startAt = 10
        val maxResults = 100

        val expectedArgs = mockk<IssueSearchCommand>()
        every { jiraCommandFactory.makeSearchForAllIssuesForProjectCommand(projectKey, startAt, maxResults) } returns expectedArgs

        jiraRepository.fetchIssuesForProject(projectKey, startAt, maxResults)

        verify {
            jiraClient.issueSearch(expectedArgs)
        }
    }
})
