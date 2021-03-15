package squaddashboard.collectors.jira.repository.sql

import squaddashboard.collectors.jira.model.IngestionType
import java.sql.Connection
import java.time.Instant

object StartIngestionSQL {

    private const val updateStatement = """
        UPDATE jira_config
        SET last_ingestion_run_started=?, last_ingestion_type=CAST(? AS INGESTION_TYPE)
        WHERE project_key=?
    """

    @ExperimentalStdlibApi
    fun execute(projectKey: String, ingestionType: IngestionType, startedAt: Instant, connection: Connection) {
        connection.prepareStatement(updateStatement).use { statement ->
            statement.setTimestamp(1, startedAt.asTimestamp())
            statement.setString(2, ingestionType.name)
            statement.setString(3, projectKey)
            statement.executeUpdate()
        }
    }
}
