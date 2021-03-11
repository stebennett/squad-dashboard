plugins {
    application
    kotlin("jvm")
    kotlin("kapt")
}

dependencies {
    // Kotlin
    implementation(platform(Dependencies.kotlinBom))
    implementation(Dependencies.kotlinStdLibJdk8)
    implementation(Dependencies.guava)

    // HttpClient
    Dependencies.httpClientCore.forEach { implementation(it) }
    kapt(Dependencies.moshiCodeGen)

    // Databases
    Dependencies.dbCore.forEach { implementation(it) }

    // Tests
    Dependencies.testsCore.forEach { testImplementation(it) }
    testImplementation(project(":common-test"))
}

application {
    mainClass.set("squaddashboard.collectors.datadog.AppKt")
}
