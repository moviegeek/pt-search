IMAGENAME = gcr.io/laputa/pt-search-backend

ENVS = $(shell cat .env | tr '\n' ',')

.PHONY: build build-front deploy deploy-front deploy-back

build:
	gcloud builds submit --project laputa --tag $(IMAGENAME)

build-front:
	npm run build

deploy-front: build-front
	gsutil cp -r ./build/* gs://pt-search/

deploy-back:
	gcloud beta run \
		deploy pt-search-backend \
		--region us-central1 \
		--platform managed \
		--concurrency 1 \
		--max-instances 1 \
		--memory 2Gi \
		--timeout 15m \
		--image $(IMAGENAME) \
		--project pt-search-9c319

deploy: deploy-front deploy-back
