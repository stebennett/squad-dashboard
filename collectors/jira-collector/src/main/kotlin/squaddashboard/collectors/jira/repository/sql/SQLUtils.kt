package squaddashboard.collectors.jira.repository.sql

import java.sql.Timestamp
import java.time.Instant

fun Instant.asTimestamp(): Timestamp =
    Timestamp.from(this)
