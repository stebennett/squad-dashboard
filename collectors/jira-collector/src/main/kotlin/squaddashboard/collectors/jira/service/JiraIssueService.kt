package squaddashboard.collectors.jira.service

import squaddashboard.collectors.jira.client.model.JiraIssue
import squaddashboard.collectors.jira.repository.JiraRepository

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
