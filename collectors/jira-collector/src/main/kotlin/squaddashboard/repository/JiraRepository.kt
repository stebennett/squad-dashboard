package squaddashboard.repository

import squaddashboard.client.jira.JiraClient
import squaddashboard.client.jira.JiraCommandFactory
import squaddashboard.client.jira.model.JiraSearchResponse

class JiraRepository(private val jiraClient: JiraClient, private val jiraCommandFactory: JiraCommandFactory) {

    fun fetchAllIssuesForProject(project: String, startAt: Int = 0, maxResults: Int): JiraSearchResponse {
        val issueSearchCommand = jiraCommandFactory.makeProjectIssuesCommand(project, startAt, maxResults)
        return jiraClient.issueSearch(issueSearchCommand)
    }
}
