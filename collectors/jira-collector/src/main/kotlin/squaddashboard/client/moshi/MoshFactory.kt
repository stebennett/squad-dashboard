package squaddashboard.client.moshi

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import squaddashboard.client.moshi.adapter.ZonedDateTimeAdapter

class MoshFactory {

    fun make(): Moshi =
        Moshi.Builder()
            .add(ZonedDateTimeAdapter())
            .add(KotlinJsonAdapterFactory())
            .build()
}
