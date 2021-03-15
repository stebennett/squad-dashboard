package squaddashboard.collectors.jira.repository.sql

import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition
import java.sql.Connection

object InsertTransitionsSQL {

    private const val insertStatement = """
        INSERT INTO jira_transitions
        (jira_data_id, jira_id, jira_transition_to, jira_transition_at)
        VALUES
        ((SELECT _id FROM jira_data WHERE jira_id = ?), ?, ?, ?)
        ON CONFLICT(jira_id)
        DO NOTHING
    """

    fun execute(issueId: Int, transitions: List<SquadDashboardJiraIssueTransition>, connection: Connection) {
        connection.prepareStatement(insertStatement).use { statement ->
            transitions.forEach {
                statement.setInt(1, issueId)
                statement.setString(2, it.jiraId.toString())
                statement.setString(3, it.transitionTo)
                statement.setTimestamp(4, it.transitionAt.asTimestamp())
            }
            statement.executeBatch()
        }
    }
}
