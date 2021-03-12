package squaddashboard.collectors.jira.repository

import org.testcontainers.containers.JdbcDatabaseContainer

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

fun JdbcDatabaseContainer<*>.getJiraConfigWorkStartState(projectKey: String): String {
    createConnection("").use { connection ->
        connection.prepareStatement("SELECT work_start_state FROM jira_config WHERE project_key=?").use { statement ->
            statement.setString(1, projectKey)
            val resultSet = statement.executeQuery()
            resultSet.next()
            return resultSet.getString(1)
        }
    }
}
