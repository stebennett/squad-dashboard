plugins {
    application
    kotlin("jvm")
}

dependencies {
    implementation(platform(Dependencies.kotlinBom))
    implementation(Dependencies.kotlinStdLibJdk8)
    implementation(Dependencies.guava)

    // Tests
    Dependencies.testsCore.forEach { testImplementation(it) }
}

application {
    mainClass.set("squaddashboard.AppKt")
}
