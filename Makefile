.PHONY: default
default: build-all

.PHONY: build-all
build-all: 
	bazel build //...

.PHONY: deps
deps:
	bazel run //:gazelle

.PHONY: test
test:
	bazel test //...


.PHONY: metabase
metabase: 
	docker-compose up -d metabasedb metabase

.PHONY: db
db: 
	docker-compose up -d db

.PHONY: jiracollector
jiracollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
	docker-compose up jiracollector

.PHONY: githubprcollector
githubprcollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/githubprcollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/githubprcollector:image
	docker-compose up githubprcollector

.PHONY: jiraissuecalculator
jiraissuecalculator: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiraissuecalculator:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiraissuecalculator:image
	docker-compose up jiraissuecalculator

.PHONY: migrate
migrate: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dbmigration:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dbmigration:image
	docker-compose up migration
