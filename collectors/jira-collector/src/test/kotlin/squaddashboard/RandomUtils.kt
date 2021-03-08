package squaddashboard

import kotlin.random.Random
import java.time.ZonedDateTime

fun Random.nextZonedDateTime(): ZonedDateTime =
    ZonedDateTime.now().plusSeconds(nextLong((60 * 60 * 24 * 365) * -2, 0))

fun <T> Random.nextFromList(list: List<T>): T =
    list[Random.nextInt(0, list.size)]
