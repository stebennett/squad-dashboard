package squaddashboard.repository

import squaddashboard.client.jira.JiraClient

class JiraRepository(private val jiraClient: JiraClient) {

    fun fetchAllIssuesForProject(project: String, startAt: Int = 0, count: Int) {
        // TODO
        throw NotImplementedError("Not yet implemented")
    }
}
