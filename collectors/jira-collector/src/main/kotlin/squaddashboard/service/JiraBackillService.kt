package squaddashboard

import squaddashboard.repository.JiraRepository
import java.awt.KeyEventPostProcessor

class JiraBackillService(private val jiraRepository: JiraRepository) {

    fun loadIssues(projectKey: String, issueProcessor: () -> Unit) {
        throw NotImplementedError("Not yet implemented")
    }
}
