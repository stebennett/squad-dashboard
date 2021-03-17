package squaddashboard.collectors.jira.model

import java.time.Instant

data class SquadDashboardJiraIssue(
    val jiraId: Int,
    val jiraKey: String,
    val jiraWorkType: JiraWorkType,
    val jiraCreatedAt: Instant,
    val jiraWorkStartedAt: Instant? = null,
    val transitions: List<SquadDashboardJiraIssueTransition>,
    val jiraProjectKey: String,
)

data class SquadDashboardJiraIssueTransition(
    val jiraId: Int,
    val transitionTo: String,
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
