package squaddashboard.collectors.jira.repository.sql

import java.sql.Connection

object CreateProjectConfigSQL {

    private const val insertStatement = """INSERT INTO jira_config (project_key, work_start_state) VALUES (?, ?);"""

    fun execute(projectKey: String, workStartState: String, connection: Connection) {
        connection.prepareStatement(insertStatement).use { statement ->
            statement.setString(1, projectKey)
            statement.setString(2, workStartState)
            statement.executeUpdate()
        }
    }
}
