package squaddashboard.collectors.jira.repository

import squaddashboard.client.jira.JiraClient
import squaddashboard.client.jira.JiraCommandFactory
import squaddashboard.client.jira.model.JiraSearchResponse

class JiraRepository(private val jiraClient: JiraClient, private val jiraCommandFactory: JiraCommandFactory) {

    fun fetchIssuesForProject(projectKey: String, startAt: Int = 0, maxResults: Int): JiraSearchResponse {
        val issueSearchCommand = jiraCommandFactory.makeProjectIssuesCommand(projectKey, startAt, maxResults)
        return jiraClient.issueSearch(issueSearchCommand)
    }
}
