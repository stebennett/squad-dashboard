package squaddashboard.client

data class ClientConfig(val baseUrl: String, val connectTimeout: Long = 30, val readTimeout: Long = 30)
