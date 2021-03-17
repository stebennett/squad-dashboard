package squaddashboard.collectors.jira.client.model

class JiraCommandFactory {

    fun makeSearchForAllIssuesForProjectCommand(projectKey: String, startAt: Int, maxResults: Int): IssueSearchCommand =
        IssueSearchCommand(
            jql = "project = $projectKey order by created asc",
            startAt = startAt,
            maxResults = maxResults,
            fields = listOf("status", "issuetype", "created", "updated", "summary", "resolutiondate"),
            expand = listOf("changelog")
        )
}
