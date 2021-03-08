package squaddashboard.repository

import com.zaxxer.hikari.HikariDataSource
import squaddashboard.model.SquadDashboardJiraIssue
import squaddashboard.model.SquadDashboardJiraIssueTransition
import java.sql.Connection
import java.sql.Timestamp
import java.time.ZonedDateTime

class SquadDashboardJiraIssueRepository(private val dataSource: HikariDataSource) {

    fun saveIssue(issue: SquadDashboardJiraIssue) {
        dataSource.connection.use { connection ->
            InsertIssueSQL.execute(issue, connection)
            InsertTransitionSQL.execute(issue.jiraId, issue.transitions, connection)
        }
    }
}

object InsertIssueSQL {

    private const val insertStatement = """
        INSERT INTO jira_data
        (jira_id, jira_key, jira_work_type, jira_created_at)
        VALUES
        (?, ?, ?, ?)
        ON CONFLICT (jira_id, jira_key) 
        DO UPDATE SET jira_work_type=? 
    """

    fun execute(issue: SquadDashboardJiraIssue, connection: Connection) {
        connection.prepareStatement(insertStatement).use { statement ->
            statement.setLong(1, issue.jiraId)
            statement.setString(2, issue.jiraKey)
            statement.setObject(3, issue.jiraWorkType)
            statement.setTimestamp(4, issue.jiraCreatedAt.asTimestamp())
            statement.setObject(4, issue.jiraWorkType)
            statement.executeUpdate()
        }
    }
}

object InsertTransitionSQL {

    private const val insertStatement = """
        INSERT INTO jira_transitions
        (jira_data_id, jira_id, jira_transition_to, jira_transition_at)
        VALUES
        ((SELECT _id FROM jira_data WHERE jira_id = ?), ?, ?, ?)
        ON CONFLICT(jira_id)
        DO NOTHING
    """

    fun execute(issueId: Long, transitions: List<SquadDashboardJiraIssueTransition>, connection: Connection) {
        connection.prepareStatement(insertStatement).use { statement ->
            transitions.forEach {
                statement.setLong(1, issueId)
                statement.setString(2, it.jiraId.toString())
                statement.setString(3, it.transitionTo)
                statement.setTimestamp(4, it.transitionAt.asTimestamp())
            }
            statement.executeBatch()
        }
    }
}

private fun ZonedDateTime.asTimestamp(): Timestamp =
    Timestamp.valueOf(this.toLocalDateTime())
