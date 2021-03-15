package squaddashboard.collectors.jira.repository

import squaddashboard.collectors.jira.model.IngestionType
import java.time.Instant

data class JiraConfig(
    val projectKey: String,
    val workStartState: String?,
    val lastIngestionType: IngestionType?,
    val lastIngestionStarted: Instant?,
    val lastIngestionCompleted: Instant?,
)
