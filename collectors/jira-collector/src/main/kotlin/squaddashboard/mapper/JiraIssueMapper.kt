package squaddashboard.mapper

import squaddashboard.client.jira.model.ChangeLog
import squaddashboard.client.jira.model.JiraIssue
import squaddashboard.model.JiraWorkType
import squaddashboard.model.SquadDashboardJiraIssue
import squaddashboard.model.SquadDashboardJiraIssueTransition

class JiraIssueMapper {

    @ExperimentalStdlibApi
    fun map(jiraIssue: JiraIssue): SquadDashboardJiraIssue =
        SquadDashboardJiraIssue(
            jiraId = jiraIssue.id.toLong(),
            jiraKey = jiraIssue.key,
            jiraCreatedAt = jiraIssue.fields.created,
            jiraWorkType = JiraWorkType.workTypeValueOf(jiraIssue.fields.issueType.name.lowercase()),
            transitions = mapTransitions(jiraIssue.changelog.histories)
        )

    private fun mapTransitions(histories: List<ChangeLog>): List<SquadDashboardJiraIssueTransition> =
        // select only the history items that have a status change somewher in the log
        histories.filter {
            it.items.any {  changeDetail ->
                changeDetail.field == "status"
            }
        }.map {
            // find the status change - there can only be one!
            val statusChange = it.items.first { changeDetail ->
                changeDetail.field == "status"
            }

            SquadDashboardJiraIssueTransition(
                jiraId = it.id.toLong(),
                transitionFrom = statusChange.fromString,
                transitionTo = statusChange.toString,
                transitionAt = it.created
            )
        }
}
