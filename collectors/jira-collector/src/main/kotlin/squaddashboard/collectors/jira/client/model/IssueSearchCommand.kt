package squaddashboard.collectors.jira.client.model

data class IssueSearchCommand(
    val jql: String,
    val fields: List<String>,
    val expand: List<String>,
    val maxResults: Int,
    val startAt: Int,
)
