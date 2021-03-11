package squaddashboard.collectors.jira.service

import squaddashboard.client.jira.model.JiraIssue
import squaddashboard.repository.JiraRepository

class JiraIssueService(private val jiraRepository: JiraRepository) {

    fun loadIssues(projectKey: String, startAt: Int = 0, batchCount: Int = 50, issueProcessor: (issue: JiraIssue) -> Unit) {
        val response = jiraRepository.fetchIssuesForProject(projectKey, startAt, batchCount)
        response.issues.forEach {
            issueProcessor(it)
        }

        if (response.total > (response.startAt + response.maxResults)) {
            loadIssues(projectKey, startAt + batchCount, batchCount, issueProcessor)
        }
    }
}
