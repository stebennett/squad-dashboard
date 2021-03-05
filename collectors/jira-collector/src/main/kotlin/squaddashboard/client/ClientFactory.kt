package squaddashboard.client

import com.squareup.moshi.Moshi
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.util.concurrent.TimeUnit

class ClientFactory {

    inline fun <reified T> make(clientConfig: ClientConfig, moshi: Moshi): T {
        val okHttpClient = OkHttpClient().newBuilder()
            .readTimeout(clientConfig.readTimeout, TimeUnit.SECONDS)
            .connectTimeout(clientConfig.connectTimeout, TimeUnit.SECONDS)
            .build()

        val retrofit = Retrofit.Builder()
            .baseUrl(clientConfig.baseUrl)
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .client(okHttpClient)
            .build()

        return retrofit.create(T::class.java)
    }
}
