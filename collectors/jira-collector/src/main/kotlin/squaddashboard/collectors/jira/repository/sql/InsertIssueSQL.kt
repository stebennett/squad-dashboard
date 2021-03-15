package squaddashboard.collectors.jira.repository.sql

import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue
import java.sql.Connection

object InsertIssueSQL {

    private const val insertStatement = """
        INSERT INTO jira_data
        (jira_id, jira_key, jira_work_type, jira_created_at, jira_project_key)
        VALUES
        (?, ?, CAST(? AS WORK_TYPE), ?, ?)
        ON CONFLICT (jira_id) 
        DO UPDATE SET jira_work_type= CAST(? AS WORK_TYPE) 
    """

    fun execute(issue: SquadDashboardJiraIssue, connection: Connection) {
        connection.prepareStatement(insertStatement).use { statement ->
            statement.setInt(1, issue.jiraId)
            statement.setString(2, issue.jiraKey)
            statement.setString(3, issue.jiraWorkType.typeName)
            statement.setTimestamp(4, issue.jiraCreatedAt.asTimestamp())
            statement.setString(5, issue.jiraProjectKey)
            statement.setString(6, issue.jiraWorkType.typeName)
            statement.executeUpdate()
        }
    }
}
