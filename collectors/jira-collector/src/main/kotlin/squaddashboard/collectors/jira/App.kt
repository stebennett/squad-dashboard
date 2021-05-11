package squaddashboard.collectors.jira

import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import squaddashboard.collectors.common.client.ClientConfig
import squaddashboard.collectors.common.client.ClientFactory
import squaddashboard.collectors.jira.client.JiraClientFactory
import squaddashboard.collectors.jira.client.model.JiraCommandFactory
import squaddashboard.collectors.jira.client.moshi.MoshiFactory
import squaddashboard.collectors.jira.mapper.JiraIssueMapper
import squaddashboard.collectors.jira.repository.JiraRepository
import squaddashboard.collectors.jira.repository.SquadDashboardJiraIssueRepository
import squaddashboard.collectors.jira.service.JiraIssueService

@ExperimentalStdlibApi
fun main() {
    val clientConfig = ClientConfig(System.getenv("JIRA_BASE_URL"))
    val clientFactory = ClientFactory()
    val moshiFactory = MoshiFactory()

    val jiraClient = JiraClientFactory(clientConfig, clientFactory, moshiFactory).make()
    val jiraCommandFactory = JiraCommandFactory()
    val jiraRepository = JiraRepository(jiraClient, jiraCommandFactory)
    val jiraIssueService = JiraIssueService(jiraRepository)

    val hikariConfig = HikariConfig()
    hikariConfig.username = System.getenv("DATABASE_USERNAME")
    hikariConfig.password = System.getenv("DATABASE_PASSWORD")
    hikariConfig.jdbcUrl = System.getenv("DATABASE_JDBC_URL")
    val dataSource = HikariDataSource(hikariConfig)

    val squadDashboardJiraIssueRepository = SquadDashboardJiraIssueRepository(dataSource)
    val jiraIssueMapper = JiraIssueMapper()

    val backfillCommand = BackfillCommand(jiraIssueService, squadDashboardJiraIssueRepository, jiraIssueMapper)

    backfillCommand.run(System.getenv("JIRA_PROJECT"), System.getenv("JIRA_WORK_START_STATE"))
}
