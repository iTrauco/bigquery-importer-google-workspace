include tools/ko/rules.mk

revision := $(shell git rev-parse --verify HEAD)
dirty := $(shell git diff --quiet || echo '-dirty')
docker_tag := $(revision)$(dirty)
gcp_project := einride-github
gcp_project_number := 1069625570948
app := github.com/einride/bigquery-importer-google-workspace
docker_repo := eu.gcr.io/$(gcp_project)
docker_image := $(docker_repo)/$(app):$(docker_tag)
local_docker_image := ko.local/$(app):$(docker_tag)

.PHONY: ko-publish-local
ko-publish-local: $(ko)
	$(info [$@] publishing local Docker image $(local_docker_image)...)
	@$(ko) publish --local --preserve-import-paths -t $(docker_tag) $(app)

.PHONY: ko-publish
ko-publish: $(ko)
	$(info [$@] publishing Docker image $(docker_image)...)
	KO_DOCKER_REPO=$(docker_repo) \
		$(ko) publish --preserve-import-paths -t $(docker_tag) $(app)

.PHONY: cloud-build-local-run
cloud-build-local-run: .cloudbuild/run.yaml ko-publish-local
	$(info [$@] starting local run...)
	@cloud-build-local \
		--no-source \
		--dryrun=false \
		--config=$< \
		--substitutions=_IMAGE_NAME=$(local_docker_image)
