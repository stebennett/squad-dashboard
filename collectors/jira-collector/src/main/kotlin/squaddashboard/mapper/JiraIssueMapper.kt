package squaddashboard.mapper

import squaddashboard.client.jira.model.JiraIssue
import squaddashboard.model.JiraWorkType
import squaddashboard.model.SquadDashboardJiraIssue

class JiraIssueMapper {

    @ExperimentalStdlibApi
    fun map(jiraIssue: JiraIssue): SquadDashboardJiraIssue =
        SquadDashboardJiraIssue(
            jiraId = jiraIssue.id.toLong(),
            jiraKey = jiraIssue.key,
            jiraCreatedAt = jiraIssue.fields.created,
            jiraWorkType = JiraWorkType.workTypeValueOf(jiraIssue.fields.issueType.name.lowercase()),
            transitions = emptyList()
        )
}
