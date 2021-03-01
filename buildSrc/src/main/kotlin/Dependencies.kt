object Dependencies {

    object Versions {

        // Kotlin
        const val kotlinVersion = "1.4.20"

        // Guava
        const val guavaVersion = "29.0-jre"

        // Testing
        const val junitVersion = "5.7.1"
        const val mockkVersion = "1.10.6"
        const val kotestVersion = "4.0.7"
        const val kotlinxCoroutinesTestVersion = "1.4.2"
    }

    // Kotlin Core
    const val kotlinStdLibJdk8 = "org.jetbrains.kotlin:kotlin-stdlib-jdk8"
    const val kotlinBom = "org.jetbrains.kotlin:kotlin-bom"
    const val guava = "com.google.guava:guava:${Versions.guavaVersion}"

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
