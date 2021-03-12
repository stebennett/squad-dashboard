package squaddashboard.collectors.jira.repository.sql

import squaddashboard.collectors.jira.model.IngestionType
import java.sql.Connection
import java.time.ZonedDateTime

object StartIngestionSQL {

    private const val updateStatement = """
        UPDATE jira_config
        SET last_ingestion_run_started=?, last_ingestion_type=?
        WHERE project_key=?
    """

    fun execute(projectKey: String, ingestionType: IngestionType, startedAt: ZonedDateTime, connection: Connection) {
        connection.prepareStatement(updateStatement).use { statement ->
            statement.setTimestamp(1, startedAt.asTimestamp())
            statement.setObject(2, ingestionType)
            statement.setString(3, projectKey)
            statement.executeUpdate()
        }
    }
}
