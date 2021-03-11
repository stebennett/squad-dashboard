package squaddashboard.collectors.jira.mapper

import io.kotest.core.spec.style.FunSpec
import io.kotest.inspectors.forAll
import io.kotest.matchers.collections.shouldContainExactlyInAnyOrder
import io.kotest.matchers.shouldBe
import squaddashboard.collectors.jira.JiraFixtures
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition
import squaddashboard.collectors.jira.nextFromList
import squaddashboard.collectors.jira.nextZonedDateTime
import kotlin.random.Random

@ExperimentalStdlibApi
class JiraIssueMapperTest : FunSpec({

    val jiraIssueMapper = JiraIssueMapper()
    val issueStates = listOf("to do", "in progress", "done", "verified")

    test("should map basic jira fields") {
        val jiraIssue = JiraFixtures.JiraIssueFixture.create("ABC")

        val result = jiraIssueMapper.map(jiraIssue)

        result.jiraId shouldBe jiraIssue.id.toLong()
        result.jiraKey shouldBe jiraIssue.key
        result.jiraCreatedAt shouldBe jiraIssue.fields.created
    }

    test("should map a Jira issue type correctly") {
        JiraWorkType.values().forAll {
            val jiraIssue = JiraFixtures.JiraIssueFixture.create("ABC", it)
            jiraIssueMapper.map(jiraIssue).jiraWorkType shouldBe it
        }
    }

    test("should map a set of jira transitions") {
        val transitions = (1..10).map {
            SquadDashboardJiraIssueTransition(
                Random.nextLong(),
                Random.nextFromList(issueStates),
                Random.nextFromList(issueStates),
                Random.nextZonedDateTime()
            )
        }
        val jiraChangeLogs = JiraFixtures.JiraChangeLogFixture.create(transitions)
        val jiraIssue = JiraFixtures.JiraIssueFixture.create("ABC", jiraChangeLogs = jiraChangeLogs)

        val result = jiraIssueMapper.map(jiraIssue)

        result.transitions shouldContainExactlyInAnyOrder transitions
    }
})
