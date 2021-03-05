package squaddashboard.client

import squaddashboard.client.jira.JiraClient
import squaddashboard.client.moshi.MoshFactory

class JiraClientFactory(
    private val clientConfig: ClientConfig,
    private val clientFactory: ClientFactory,
    private val moshiFactory: MoshFactory,
) {

    fun make(): JiraClient {
        val moshi = moshiFactory.make()

        return clientFactory.make(clientConfig, moshi)
    }
}
