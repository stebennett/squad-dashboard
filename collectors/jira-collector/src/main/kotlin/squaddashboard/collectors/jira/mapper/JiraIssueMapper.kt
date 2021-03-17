package squaddashboard.collectors.jira.mapper

import squaddashboard.collectors.jira.client.model.ChangeLogs
import squaddashboard.collectors.jira.client.model.JiraIssue
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition

class JiraIssueMapper {

    @ExperimentalStdlibApi
    fun map(jiraIssue: JiraIssue, jiraProjectKey: String): SquadDashboardJiraIssue =
        SquadDashboardJiraIssue(
            jiraId = jiraIssue.id.toInt(),
            jiraKey = jiraIssue.key,
            jiraProjectKey = jiraProjectKey,
            jiraCreatedAt = jiraIssue.fields.created,
            jiraCompletedAt = jiraIssue.fields.resolutiondate,
            jiraWorkType = JiraWorkType.workTypeValueOf(jiraIssue.fields.issueType.name.lowercase()),
            transitions = mapTransitions(jiraIssue.changelog),
        )

    private fun mapTransitions(changeLogs: ChangeLogs): List<SquadDashboardJiraIssueTransition> =
        // select only the history items that have a status change somewhere in the log
        changeLogs.statusChanges().mapNotNull { changeLog ->

            // find the status change - there can only be one!
            changeLog.statusChange()?.let {
                SquadDashboardJiraIssueTransition(
                    jiraId = changeLog.id.toInt(),
                    transitionTo = it.toString,
                    transitionAt = changeLog.created
                )
            }
        }
}
