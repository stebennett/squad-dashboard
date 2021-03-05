package squaddashboard.client.moshi.adapter

import com.squareup.moshi.FromJson
import com.squareup.moshi.ToJson
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

class ZonedDateTimeAdapter {

    @FromJson
    fun fromJson(dateTimeString: String): ZonedDateTime {
        return ZonedDateTime.parse(dateTimeString, DateTimeFormatter.ISO_ZONED_DATE_TIME)
    }

    @ToJson
    fun toJson(zonedDateTime: ZonedDateTime): String {
        return zonedDateTime.format(DateTimeFormatter.ISO_ZONED_DATE_TIME)
    }
}
