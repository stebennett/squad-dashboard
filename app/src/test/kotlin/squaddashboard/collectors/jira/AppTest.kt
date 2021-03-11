package squaddashboard.collectors.jira

import io.kotest.core.spec.style.FunSpec
import io.kotest.matchers.shouldNotBe

class AppTest : FunSpec({

    test("App has a greeting") {
        App().greeting shouldNotBe null
    }
})
