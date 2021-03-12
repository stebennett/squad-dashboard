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
        const val moshiVersion = "1.11.0"

        // Databases
        const val hikariCPVersion = "4.0.3"
        const val postgresDriverVersion = "42.2.16"

        // Testing
        const val junitVersion = "5.7.1"
        const val mockkVersion = "1.10.6"
        const val kotestVersion = "4.4.1"
        const val kotlinxCoroutinesTestVersion = "1.4.2"
        const val kotlinFakerVersion = "1.6.0"

        // Postgres container
        const val testContainersVersion = "1.15.2"

        // Flyway
        const val flywayVersion = "7.7.0"

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
    const val moshi = "com.squareup.moshi:moshi-kotlin:${Versions.moshiVersion}"
    const val converterMoshi = "com.squareup.retrofit2:converter-moshi:${Versions.retrofit2Version}"
    const val moshiCodeGen = "com.squareup.moshi:moshi-kotlin-codegen:${Versions.moshiVersion}"

    val httpClientCore = listOf(retrofit, okhttp, converterMoshi, moshi)

    // Database integration
    const val hikariCP = "com.zaxxer:HikariCP:${Versions.hikariCPVersion}"
    const val postgresDriver = "org.postgresql:postgresql:${Versions.postgresDriverVersion}"

    val dbCore = listOf(postgresDriver, hikariCP)

    // Testing libraries
    const val jupiter = "org.junit.jupiter:junit-jupiter:${Versions.junitVersion}"
    const val mockk = "io.mockk:mockk:${Versions.mockkVersion}"
    const val kotestRunner = "io.kotest:kotest-runner-junit5:${Versions.kotestVersion}"
    const val kotestProperties = "io.kotest:kotest-property:${Versions.kotestVersion}"
    const val kotestAssertions = "io.kotest:kotest-assertions-core-jvm:${Versions.kotestVersion}"
    const val kotestArrowAssertions = "io.kotest:kotest-assertions-arrow:${Versions.kotestVersion}"
    const val kotlinxCoroutinesTest = "org.jetbrains.kotlinx:kotlinx-coroutines-test:${Versions.kotlinxCoroutinesTestVersion}"
    const val kotlinFaker = "io.github.serpro69:kotlin-faker:${Versions.kotlinFakerVersion}"

    // Database Containers for Testing
    const val testContainers = "org.testcontainers:testcontainers:${Versions.testContainersVersion}"
    const val testContainersPostgres = "org.testcontainers:postgresql:${Versions.testContainersVersion}"
    const val kotestTestContainers = "io.kotest:kotest-extensions-testcontainers:${Versions.kotestVersion}"
    const val flyway = "org.flywaydb:flyway-core:${Versions.flywayVersion}"

    val testsCore = listOf(jupiter, mockk, kotestRunner, kotestProperties, kotestAssertions, kotestArrowAssertions, kotlinxCoroutinesTest, kotlinFaker)
    val testDbCore = listOf(testContainers, testContainersPostgres, kotestTestContainers, flyway)
}
