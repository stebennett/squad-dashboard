package squaddashboard.collectors.jira.client.jira.model

import java.time.ZonedDateTime

data class JiraIssue(
    val id: String,
    val self: String,
    val key: String,
    val changelog: ChangeLogs,
    val fields: JiraIssueFields,
)

data class ChangeLogs(
    val histories: List<ChangeLog>,
) {

    fun statusChanges(): List<ChangeLog> =
        histories.filter {
            it.items.any {
                it.field == "status"
            }
        }
}

data class ChangeLog(
    val id: String,
    val created: ZonedDateTime,
    val items: List<ChangeDetail>,
) {
    fun statusChange(): ChangeDetail? =
        items.firstOrNull {
            it.field == "status"
        }
}

data class ChangeDetail(
    val field: String,
    val fromString: String,
    val toString: String,
)

data class JiraIssueFields(
    val summary: String,
    val issueType: JiraIssueType,
    val created: ZonedDateTime,
    val updated: ZonedDateTime,
    val status: JiraIssueStatus,
)

data class JiraIssueType(
    val name: String,
)

data class JiraIssueStatus(
    val name: String,
)
