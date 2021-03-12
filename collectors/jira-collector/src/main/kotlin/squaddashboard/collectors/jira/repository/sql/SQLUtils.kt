package squaddashboard.collectors.jira.repository.sql

import java.sql.Timestamp
import java.time.ZonedDateTime

fun ZonedDateTime.asTimestamp(): Timestamp =
    Timestamp.valueOf(this.toLocalDateTime())

