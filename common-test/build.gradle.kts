plugins {
    kotlin("jvm")
    `java-library`
}

dependencies {
    // Kotlin
    implementation(platform(Dependencies.kotlinBom))
    implementation(Dependencies.kotlinStdLibJdk8)
    implementation(Dependencies.guava)

    // Tests
    Dependencies.testsCore.forEach { implementation(it) }
}
