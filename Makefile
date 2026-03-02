update-deps:
	cd app && set "GOPROXY=direct" && set "GOSUMDB=off" && go get -u github.com/XenonPPG/KRS_CONTRACTS@master
	cd app && go mod tidy

swagger-gen:
	docker compose --profile tools run --rm swagger-gen

.PHONY: update-deps swagger-gen