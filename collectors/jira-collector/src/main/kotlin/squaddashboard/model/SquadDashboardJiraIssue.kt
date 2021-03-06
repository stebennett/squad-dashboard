package squaddashboard.model

import java.time.ZonedDateTime

data class SquadDashboardJiraIssue(
    val jiraId: Long,
    val jiraKey: String,
    val jiraWorkType: JiraWorkType,
    val jiraCreatedAt: ZonedDateTime,
    val transitions: List<SquadDashboardJiraIssueTransition>,
)

data class SquadDashboardJiraIssueTransition(
    val jiraId: Long,
    val transitionTo: String,
    val transitionFrom: String,
    val transitionAt: ZonedDateTime,
)

enum class JiraWorkType(val typeName: String) {
    Story("story"),
    Task("task"),
    Bug("bug")
}
