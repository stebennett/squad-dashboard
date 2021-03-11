package squaddashboard.collectors.jira.client.moshi

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import squaddashboard.collectors.common.moshi.adapter.ZonedDateTimeAdapter

class MoshiFactory {

    fun make(): Moshi =
        Moshi.Builder()
            .add(ZonedDateTimeAdapter())
            .add(KotlinJsonAdapterFactory())
            .build()
}
