plugins {
    kotlin("jvm")
    kotlin("kapt")
    `java-library`
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
    testImplementation(project(":common-test"))
}
