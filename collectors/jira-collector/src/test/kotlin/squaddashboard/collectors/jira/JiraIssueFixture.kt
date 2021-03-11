package squaddashboard.collectors.jira

import io.github.serpro69.kfaker.Faker
import squaddashboard.collectors.jira.client.jira.model.ChangeDetail
import squaddashboard.collectors.jira.client.jira.model.ChangeLog
import squaddashboard.collectors.jira.client.jira.model.ChangeLogs
import squaddashboard.collectors.jira.client.jira.model.JiraIssue
import squaddashboard.collectors.jira.client.jira.model.JiraIssueFields
import squaddashboard.collectors.jira.client.jira.model.JiraIssueStatus
import squaddashboard.collectors.jira.client.jira.model.JiraIssueType
import squaddashboard.collectors.jira.client.jira.model.JiraSearchResponse
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition
import kotlin.random.Random

object JiraFixtures {

    val faker = Faker()

    object JiraSearchResponseFixture {
        fun create(projectKey: String, startAt: Int, count: Int, total: Int) =
            JiraSearchResponse(
                startAt = startAt,
                maxResults = count,
                total = total,
                issues = (0..count).map {
                    JiraIssueFixture.create(projectKey, JiraWorkType.Story)
                }
            )
    }

    object JiraIssueFixture {

        fun create(projectKey: String, workType: JiraWorkType = JiraWorkType.Story, jiraChangeLogs: ChangeLogs = ChangeLogs(emptyList())): JiraIssue = JiraIssue(
            id = Random.nextLong().toString(),
            self = "a-fake-self-url",
            key = "$projectKey-${Random.nextInt(1, 500)}",
            changelog = jiraChangeLogs,
            fields = JiraIssueFields(
                summary = faker.michaelScott.quotes(),
                issueType = JiraIssueType(workType.typeName),
                created = Random.nextZonedDateTime(),
                updated = Random.nextZonedDateTime(),
                status = JiraIssueStatus("In Progress")
            ),
        )
    }

    object JiraChangeLogFixture {

        fun create(transitions: List<SquadDashboardJiraIssueTransition> = emptyList()): ChangeLogs =
            ChangeLogs(
                histories = transitions.map {
                    ChangeLog(
                        id = it.jiraId.toString(),
                        created = it.transitionAt,
                        items = listOf(
                            ChangeDetail(
                                field = "status",
                                toString = it.transitionTo,
                                fromString = it.transitionFrom,
                            )
                        )
                    )
                }
            )
    }
}