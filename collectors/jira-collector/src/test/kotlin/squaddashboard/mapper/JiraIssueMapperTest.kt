package squaddashboard.mapper

import io.kotest.core.spec.style.FunSpec
import io.kotest.inspectors.forAll
import io.kotest.matchers.shouldBe
import squaddashboard.JiraFixtures
import squaddashboard.model.JiraWorkType

@ExperimentalStdlibApi
class JiraIssueMapperTest : FunSpec({

    val jiraIssueMapper = JiraIssueMapper()

    test("should map basic jira fields") {
        val jiraIssue = JiraFixtures.JiraIssueFixure.create("ABC")

        val result = jiraIssueMapper.map(jiraIssue)

        result.jiraId shouldBe jiraIssue.id.toLong()
        result.jiraKey shouldBe jiraIssue.key
        result.jiraCreatedAt shouldBe jiraIssue.fields.created
    }

    test("should map a Jira issue type correctly") {
        JiraWorkType.values().forAll {
            val jiraIssue = JiraFixtures.JiraIssueFixure.create("ABC", it)
            jiraIssueMapper.map(jiraIssue).jiraWorkType shouldBe it
        }
    }
})
