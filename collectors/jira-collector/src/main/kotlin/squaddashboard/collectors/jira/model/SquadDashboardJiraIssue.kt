package squaddashboard.collectors.jira.model

import java.time.Instant
import java.time.ZonedDateTime

data class SquadDashboardJiraIssue(
    val jiraId: Long,
    val jiraKey: String,
    val jiraWorkType: JiraWorkType,
    val jiraCreatedAt: Instant,
    val transitions: List<SquadDashboardJiraIssueTransition>,
)

data class SquadDashboardJiraIssueTransition(
    val jiraId: Long,
    val transitionTo: String,
    val transitionFrom: String,
    val transitionAt: Instant,
)

enum class JiraWorkType(val typeName: String) {
    Story("story"),
    Task("task"),
    Bug("bug");

    companion object {
        fun workTypeValueOf(typeName: String): JiraWorkType =
            enumValues<JiraWorkType>().firstOrNull {
                it.typeName.equals(typeName, true)
            } ?: throw IllegalArgumentException("No enum constant for $typeName")
    }
}
