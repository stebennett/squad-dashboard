package squaddashboard.collectors.jira.client.jira

import squaddashboard.collectors.jira.client.jira.model.IssueSearchCommand

class JiraCommandFactory {

    fun makeProjectIssuesCommand(projectKey: String, startAt: Int, maxResults: Int): IssueSearchCommand =
        IssueSearchCommand(
            jql = "project = $projectKey",
            startAt = startAt,
            maxResults = maxResults,
            fields = listOf("status", "issuetype", "created", "updated", "summary"),
            expand = listOf("changelog")
        )
}