object Dependencies {

    object Versions {

        val jvmTargetVersion = "1.8"

        // Kotlin
        const val kotlinVersion = "1.4.20"

        // Kapt
        const val kaptPluginVersion = "1.4.31"

        // Guava
        const val guavaVersion = "29.0-jre"

        // OkHttp, Retrofit, Moshi
        const val retrofit2Version = "2.9.0"
        const val okhttpVersion = "4.9.1"
        const val moshiVersion = ""

        // Testing
        const val junitVersion = "5.7.1"
        const val mockkVersion = "1.10.6"
        const val kotestVersion = "4.4.1"
        const val kotlinxCoroutinesTestVersion = "1.4.2"

        // Code style
        const val ktlintVersion = "10.0.0"
    }

    // Kotlin Core
    const val kotlinStdLibJdk8 = "org.jetbrains.kotlin:kotlin-stdlib-jdk8"
    const val kotlinBom = "org.jetbrains.kotlin:kotlin-bom"
    const val guava = "com.google.guava:guava:${Versions.guavaVersion}"

    // OKHttp, Retrofit, Moshi
    const val retrofit = "com.squareup.retrofit2:retrofit:${Versions.retrofit2Version}"
    const val okhttp = "com.squareup.okhttp3:okhttp:${Versions.okhttpVersion}"
    const val moshi = "com.squareup.moshi:moshi:${Versions.moshiVersion}"
    const val moshiCodeGen = "com.squareup.moshi:moshi-kotlin-codegen:${Versions.moshiVersion}"

    val httpClientCore = listOf(retrofit, okhttp, moshi)

    // Testing libraries
    val jupiter = "org.junit.jupiter:junit-jupiter:${Versions.junitVersion}"
    val mockk = "io.mockk:mockk:${Versions.mockkVersion}"
    val kotestRunner = "io.kotest:kotest-runner-junit5:${Versions.kotestVersion}"
    val kotestProperties = "io.kotest:kotest-property:${Versions.kotestVersion}"
    val kotestAssertions = "io.kotest:kotest-assertions-core-jvm:${Versions.kotestVersion}"
    val kotestArrowAssertions = "io.kotest:kotest-assertions-arrow:${Versions.kotestVersion}"
    val kotlinxCoroutinesTest = "org.jetbrains.kotlinx:kotlinx-coroutines-test:${Versions.kotlinxCoroutinesTestVersion}"

    val testsCore = listOf(jupiter, mockk, kotestRunner, kotestProperties, kotestAssertions, kotestArrowAssertions, kotlinxCoroutinesTest)


}
