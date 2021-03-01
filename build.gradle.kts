plugins {
    kotlin("jvm") version Dependencies.Versions.kotlinVersion
    id("org.jlleitschuh.gradle.ktlint") version Dependencies.Versions.ktlintVersion
}

allprojects {
    repositories {
        mavenCentral()
    }

    tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile>().configureEach {
        kotlinOptions.jvmTarget = Dependencies.Versions.jvmTargetVersion
        kotlinOptions.freeCompilerArgs += "-Xexperimental=io.ktor.util.KtorExperimentalAPI"
        kotlinOptions.freeCompilerArgs += "-Xexperimental=io.ktor.locations.KtorExperimentalLocationsAPI"
    }

    tasks.withType<Test>() {
        useJUnitPlatform()
        testLogging {
            events("skipped", "failed")
        }
    }
}

subprojects {
    apply(plugin = "org.jlleitschuh.gradle.ktlint")
}
