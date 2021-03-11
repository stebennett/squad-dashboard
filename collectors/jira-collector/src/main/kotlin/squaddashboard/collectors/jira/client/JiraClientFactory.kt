package squaddashboard.collectors.jira.client

import squaddashboard.collectors.common.client.ClientConfig
import squaddashboard.collectors.common.client.ClientFactory
import squaddashboard.collectors.jira.client.moshi.MoshiFactory

class JiraClientFactory(
    private val clientConfig: ClientConfig,
    private val clientFactory: ClientFactory,
    private val moshiFactory: MoshiFactory,
) {

    fun make(): JiraClient {
        val moshi = moshiFactory.make()

        return clientFactory.make(clientConfig, moshi)
    }
}
