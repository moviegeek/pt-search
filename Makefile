IMAGENAME = gcr.io/laputa/pt-search-backend

ENVS = $(shell cat .env | tr '\n' ',')

.PHONY: build deploy docker-run

build:
	gcloud builds submit --project laputa --tag $(IMAGENAME)

deploy:
	gcloud beta run \
		deploy pt-search-backend \
		--region us-central1 \
		--platform managed \
		--concurrency 1 \
		--max-instances 1 \
		--memory 2Gi \
		--timeout 15m \
		--image $(IMAGENAME) \
		--project movie-221500
