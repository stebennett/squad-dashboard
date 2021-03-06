package squaddashboard

import io.github.serpro69.kfaker.Faker
import kotlin.random.Random
import squaddashboard.client.jira.model.ChangeLogs
import squaddashboard.client.jira.model.JiraIssue
import squaddashboard.client.jira.model.JiraIssueFields
import squaddashboard.client.jira.model.JiraIssueStatus
import squaddashboard.client.jira.model.JiraIssueType
import squaddashboard.client.jira.model.JiraSearchResponse
import java.time.ZonedDateTime

object JiraFixtures {

    val faker = Faker()

    object JiraSearchResponseFixture {
        fun create(projectKey: String, startAt: Int, count: Int, total: Int) =
            JiraSearchResponse(
                startAt = startAt,
                maxResults = count,
                total = total,
                issues = (0..count).map {
                    JiraIssueFixure.create(projectKey)
                }
            )
    }

    object JiraIssueFixure {

        fun create(projectKey: String): JiraIssue = JiraIssue(
            id = faker.idNumber.toString(),
            self = "a-fake-self-url",
            key = "$projectKey-${Random.nextInt(1, 500)}",
            changelog = ChangeLogs(emptyList()),
            fields = JiraIssueFields(
                summary = faker.michaelScott.quotes(),
                issueType = JiraIssueType("Story"),
                created = Random.nextZonedDateTime(),
                updated = Random.nextZonedDateTime(),
                status = JiraIssueStatus("In Progress")
            ),
        )
    }
}

fun Random.nextZonedDateTime(): ZonedDateTime =
    ZonedDateTime.now().plusSeconds(nextLong((60 * 60 * 24 * 365) * -2, 0))
