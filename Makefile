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
metabase: db
	docker-compose --env-file ${ENV} up -d metabasedb metabase

.PHONY: db
db: 
	docker-compose --env-file ${ENV} up -d db

.PHONY: jiracollector
jiracollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
	echo ${ENV}
	docker-compose --env-file ${ENV} up jiracollector

.PHONY: init
init: build-all db migrate
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/configloader:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/configloader:image
	docker-compose --env-file ${ENV} up configloader

.PHONY: githubprcollector
githubprcollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/githubprcollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/githubprcollector:image
	docker-compose --env-file ${ENV} up githubprcollector

.PHONY: googlesheetcollector
googlesheetcollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/googlesheetcollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/googlesheetcollector:image
	docker-compose --env-file ${ENV} up googlesheetcollector

.PHONY: pagerdutyoncallcollector
pagerdutyoncallcollector: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/pagerdutyoncallcollector:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/pagerdutyoncallcollector:image
	docker-compose --env-file ${ENV} up pagerdutyoncallcollector

.PHONY: jiraissuecalculator
jiraissuecalculator: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiraissuecalculator:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiraissuecalculator:image
	docker-compose --env-file ${ENV} up jiraissuecalculator

.PHONY: migrate
migrate: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dbmigration:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dbmigration:image
	docker-compose --env-file ${ENV} up migration

.PHONY: api
api: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/api:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/api:image
	docker-compose --env-file ${ENV} up api	

.PHONY: dashboardcli
dashboardcli: build-all db
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dashboardcli:image
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/dashboardcli:image
	docker-compose --env-file ${ENV} up dashboardcli	

.PHONY: report
report: init jiracollector jiraissuecalculator dashboardcli
