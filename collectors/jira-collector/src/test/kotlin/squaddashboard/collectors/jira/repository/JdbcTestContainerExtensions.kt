package squaddashboard.collectors.jira.repository

import io.mockk.stackTracesAlignmentValueOf
import org.testcontainers.containers.JdbcDatabaseContainer
import squaddashboard.collectors.jira.model.IngestionType
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssue

fun JdbcDatabaseContainer<*>.getJiraConfigCount(projectKey: String): Long {
    createConnection("").use { connection ->
        connection.prepareStatement("SELECT COUNT(project_key) FROM jira_config WHERE project_key=?").use { statement ->
            statement.setString(1, projectKey)
            val resultSet = statement.executeQuery()
            resultSet.next()
            return resultSet.getLong(1)
        }
    }
}

fun JdbcDatabaseContainer<*>.getJiraConfig(projectKey: String): JiraConfig {
    createConnection("").use { connection ->
        connection.prepareStatement("SELECT project_key, work_start_state, last_ingestion_type, last_ingestion_run_started, last_ingestion_run_completed FROM jira_config WHERE project_key = ?").use { statement ->
            statement.setString(1, projectKey)
            val resultSet = statement.executeQuery()
            resultSet.next()

            return JiraConfig(
                resultSet.getString(1),
                resultSet.getString(2),
                resultSet.getString(3)?.let { IngestionType.valueOf(it) },
                resultSet.getTimestamp(4)?.toInstant(),
                resultSet.getTimestamp(5)?.toInstant(),
            )
        }
    }
}

fun JdbcDatabaseContainer<*>.createJiraConfig(projectKey: String) {
    createConnection("").use { connection ->
        connection.prepareStatement("INSERT INTO jira_config (project_key, work_start_state) VALUES (?, ?)").use { statement ->
            statement.setString(1, projectKey)
            statement.setString(2, "In Progress")
            statement.executeUpdate()
        }
    }
}


fun JdbcDatabaseContainer<*>.getJiraIssueCountForProject(projectKey: String): Long {
    createConnection("").use { connection ->
        connection.prepareStatement("SELECT COUNT(_id) FROM jira_data WHERE jira_project_key = ?").use { statement ->
            statement.setString(1, projectKey)
            val resultSet = statement.executeQuery()
            resultSet.next()
            return resultSet.getLong(1)
        }
    }
}

fun JdbcDatabaseContainer<*>.getJiraIssue(jiraIssueId: Int): SquadDashboardJiraIssue {
    createConnection("").use { connection ->
        connection.prepareStatement("SELECT jira_id, jira_key, jira_project_key, jira_work_type, jira_created_at FROM jira_data WHERE jira_id = ?").use { statement ->
            statement.setInt(1, jiraIssueId)
            val resultSet = statement.executeQuery()
            resultSet.next()
            return SquadDashboardJiraIssue(
                jiraId = resultSet.getInt(1),
                jiraKey = resultSet.getString(2),
                jiraProjectKey = resultSet.getString(3),
                jiraWorkType = JiraWorkType.workTypeValueOf(resultSet.getString(4)),
                jiraCreatedAt = resultSet.getTimestamp(5).toInstant(),
                transitions = emptyList()
            )
        }
    }
}

