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

    // Tests
    Dependencies.testsCore.forEach { testImplementation(it) }
}

application {
    mainClass.set("squaddashboard.AppKt")
}
