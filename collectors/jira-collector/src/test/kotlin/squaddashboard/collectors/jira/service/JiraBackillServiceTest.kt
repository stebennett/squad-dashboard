package squaddashboard.collectors.jira.service

import io.kotest.core.spec.style.FunSpec
import io.kotest.matchers.collections.shouldContainExactlyInAnyOrder
import io.mockk.every
import io.mockk.mockk
import squaddashboard.JiraFixtures
import squaddashboard.client.jira.model.JiraIssue
import squaddashboard.repository.JiraRepository

class JiraBackillServiceTest : FunSpec({

    val jiraRepository = mockk<JiraRepository>()
    val jiraBackfillService = JiraIssueService(jiraRepository)

    test("should load issues from repository with single page of results") {
        val projectKey = "DEF"
        val response = JiraFixtures.JiraSearchResponseFixture.create(projectKey, 0, 10, 10)

        val captures = mutableListOf<JiraIssue>()

        every { jiraRepository.fetchIssuesForProject(projectKey, any(), any()) } returns response

        jiraBackfillService.loadIssues(projectKey) {
            captures.add(it)
        }

        captures shouldContainExactlyInAnyOrder response.issues
    }

    test("should load issues from repository across multiple pages") {
        val projectKey = "GHI"

        val page1Response = JiraFixtures.JiraSearchResponseFixture.create(projectKey, 0, 10, 30)
        val page2Response = JiraFixtures.JiraSearchResponseFixture.create(projectKey, 10, 10, 30)
        val page3Response = JiraFixtures.JiraSearchResponseFixture.create(projectKey, 20, 10, 30)

        every { jiraRepository.fetchIssuesForProject(projectKey, 0, 10) } returns page1Response
        every { jiraRepository.fetchIssuesForProject(projectKey, 10, 10) } returns page2Response
        every { jiraRepository.fetchIssuesForProject(projectKey, 20, 10) } returns page3Response

        val captures = mutableListOf<JiraIssue>()

        jiraBackfillService.loadIssues(projectKey, batchCount = 10) {
            captures.add(it)
        }

        captures shouldContainExactlyInAnyOrder page1Response.issues + page2Response.issues + page3Response.issues
    }
})
