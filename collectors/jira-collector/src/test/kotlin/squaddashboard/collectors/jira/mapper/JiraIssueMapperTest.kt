package squaddashboard.collectors.jira.mapper

import io.kotest.core.spec.style.FunSpec
import io.kotest.inspectors.forAll
import io.kotest.matchers.collections.shouldContainExactlyInAnyOrder
import io.kotest.matchers.shouldBe
import squaddashboard.collectors.jira.JiraFixtures
import squaddashboard.collectors.jira.model.JiraWorkType
import squaddashboard.collectors.jira.model.SquadDashboardJiraIssueTransition
import squaddashboard.common.test.nextFromList
import squaddashboard.common.test.nextZonedDateTime
import kotlin.random.Random

@ExperimentalStdlibApi
class JiraIssueMapperTest : FunSpec({

    val jiraIssueMapper = JiraIssueMapper()
    val issueStates = listOf("to do", "in progress", "done", "verified")

    test("should map basic jira fields") {
        val projectKey = "ABC"
        val jiraIssue = JiraFixtures.JiraIssueFixture.create(projectKey)

        val result = jiraIssueMapper.map(jiraIssue, projectKey)

        result.jiraId shouldBe jiraIssue.id.toLong()
        result.jiraKey shouldBe jiraIssue.key
        result.jiraCreatedAt shouldBe jiraIssue.fields.created
        result.jiraProjectKey shouldBe projectKey
    }

    test("should map a Jira issue type correctly") {
        JiraWorkType.values().forAll {
            val jiraIssue = JiraFixtures.JiraIssueFixture.create("ABC", it)
            jiraIssueMapper.map(jiraIssue, "ABC").jiraWorkType shouldBe it
        }
    }

    test("should map a set of jira transitions") {
        val transitions = (1..10).map {
            SquadDashboardJiraIssueTransition(
                Random.nextInt(),
                Random.nextFromList(issueStates),
                Random.nextFromList(issueStates),
                Random.nextZonedDateTime().toInstant()
            )
        }
        val jiraChangeLogs = JiraFixtures.JiraChangeLogFixture.create(transitions)
        val jiraIssue = JiraFixtures.JiraIssueFixture.create("ABC", jiraChangeLogs = jiraChangeLogs)

        val result = jiraIssueMapper.map(jiraIssue, "ABC")

        result.transitions shouldContainExactlyInAnyOrder transitions
    }
})
