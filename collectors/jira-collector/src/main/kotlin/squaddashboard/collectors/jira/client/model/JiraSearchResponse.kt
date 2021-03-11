package squaddashboard.collectors.jira.client.model

data class JiraSearchResponse(
    val startAt: Int,
    val maxResults: Int,
    val total: Int,

    val issues: List<JiraIssue>,
)