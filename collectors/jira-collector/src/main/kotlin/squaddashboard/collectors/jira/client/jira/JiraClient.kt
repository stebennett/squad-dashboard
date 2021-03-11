package squaddashboard.collectors.jira.client.jira

import retrofit2.http.Body
import retrofit2.http.POST
import squaddashboard.collectors.jira.client.jira.model.IssueSearchCommand
import squaddashboard.collectors.jira.client.jira.model.JiraSearchResponse

interface JiraClient {

    @POST("/rest/api/2/search")
    fun issueSearch(@Body issueSearchCommand: IssueSearchCommand): JiraSearchResponse
}
