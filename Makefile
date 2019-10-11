.PHONY: all

all: vendor dockerbuild containerup

vendor:
	@go mod downlad
	@go mod vendor

dockerbuild:
	@echo "...building image iata-finder"
	@docker build -t abelgoodwin1988/iata-finder:latest . --no-cache

containerup:
	@echo "...composing container iata-finder"
	@sh -c  "docker-compose up"