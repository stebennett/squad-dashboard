package squaddashboard.collectors.jira.client

import retrofit2.http.Body
import retrofit2.http.POST
import squaddashboard.collectors.jira.client.model.IssueSearchCommand
import squaddashboard.collectors.jira.client.model.JiraSearchResponse

interface JiraClient {

    @POST("/rest/api/2/search")
    fun issueSearch(@Body issueSearchCommand: IssueSearchCommand): JiraSearchResponse
}
