package squaddashboard.collectors.jira.mapper

import squaddashboard.collectors.jira.client.jira.model.ChangeLogs
import squaddashboard.collectors.jira.client.jira.model.JiraIssue
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition

class JiraIssueMapper {

    @ExperimentalStdlibApi
    fun map(jiraIssue: JiraIssue): SquadDashboardJiraIssue =
        SquadDashboardJiraIssue(
            jiraId = jiraIssue.id.toLong(),
            jiraKey = jiraIssue.key,
            jiraCreatedAt = jiraIssue.fields.created,
            jiraWorkType = JiraWorkType.workTypeValueOf(jiraIssue.fields.issueType.name.lowercase()),
            transitions = mapTransitions(jiraIssue.changelog)
        )

    private fun mapTransitions(changeLogs: ChangeLogs): List<SquadDashboardJiraIssueTransition> =
        // select only the history items that have a status change somewhere in the log
        changeLogs.statusChanges().mapNotNull { changeLog ->

            // find the status change - there can only be one!
            changeLog.statusChange()?.let {
                SquadDashboardJiraIssueTransition(
                    jiraId = changeLog.id.toLong(),
                    transitionFrom = it.fromString,
                    transitionTo = it.toString,
                    transitionAt = changeLog.created
                )
            }
        }
}