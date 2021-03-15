package squaddashboard.collectors.jira.repository

import com.zaxxer.hikari.HikariDataSource
import squaddashboard.collectors.jira.model.IngestionType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue
import squaddashboard.collectors.jira.repository.sql.CompleteIngestionSQL
import squaddashboard.collectors.jira.repository.sql.CreateProjectConfigSQL
import squaddashboard.collectors.jira.repository.sql.InsertIssueSQL
import squaddashboard.collectors.jira.repository.sql.InsertTransitionsSQL
import squaddashboard.collectors.jira.repository.sql.StartIngestionSQL
import java.time.Instant

class SquadDashboardJiraIssueRepository(private val dataSource: HikariDataSource) {

    fun saveIssue(issue: SquadDashboardJiraIssue) {
        dataSource.connection.use { connection ->
            InsertIssueSQL.execute(issue, connection)
            InsertTransitionsSQL.execute(issue.jiraId, issue.transitions, connection)
        }
    }

    fun createProjectConfig(projectKey: String, workStartState: String) {
        dataSource.connection.use { connection ->
            CreateProjectConfigSQL.execute(projectKey, workStartState, connection)
        }
    }

    @ExperimentalStdlibApi
    fun startIngestion(projectKey: String, ingestionType: IngestionType, startedAt: Instant) {
        dataSource.connection.use { connection ->
            StartIngestionSQL.execute(projectKey, ingestionType, startedAt, connection)
        }
    }

    fun completeIngestion(projectKey: String, completedAt: Instant) {
        dataSource.connection.use { connection ->
            CompleteIngestionSQL.execute(projectKey, completedAt, connection)
        }
    }
}
