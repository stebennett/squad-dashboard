package squaddashboard.collectors.jira.repository.sql

import java.sql.Connection
import java.time.Instant
object CompleteIngestionSQL {

    private const val updateStatement = """
        UPDATE jira_config
        SET last_ingestion_run_completed=?
        WHERE project_key=?
    """

    fun execute(projectKey: String, completedAt: Instant, connection: Connection) {
        connection.prepareStatement(updateStatement).use { statement ->
            statement.setTimestamp(1, completedAt.asTimestamp())
            statement.setString(2, projectKey)
            statement.executeUpdate()
        }
    }
}
