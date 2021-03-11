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

    // common utils
    implementation(project(":collectors:collector-common"))
}

application {
    mainClass.set("squaddashboard.AppKt")
}
