package squaddashboard.collectors.jira.repository.sql

import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue
import java.sql.Connection

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
