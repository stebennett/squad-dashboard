package squaddashboard.collectors.jira.repository

import squaddashboard.collectors.jira.client.JiraClient
import squaddashboard.collectors.jira.client.model.JiraCommandFactory
import squaddashboard.collectors.jira.client.model.JiraSearchResponse

class JiraRepository(private val jiraClient: JiraClient, private val jiraCommandFactory: JiraCommandFactory) {

    fun fetchIssuesForProject(projectKey: String, startAt: Int = 0, maxResults: Int): JiraSearchResponse {
        val issueSearchCommand = jiraCommandFactory.makeProjectIssuesCommand(projectKey, startAt, maxResults)
        return jiraClient.issueSearch(issueSearchCommand)
    }
}
