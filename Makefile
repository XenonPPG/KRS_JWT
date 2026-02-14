GOPROXY_OFF := GOPROXY=direct
GOSUMDB_OFF := GOSUMDB=off
ENV_VARS := set $(GOPROXY_OFF) && set $(GOSUMDB_OFF)

update-deps:
	cd app && $(ENV_VARS) && go clean -modcache
	cd app && $(ENV_VARS) && go get -u github.com/XenonPPG/KRS_CONTRACTS@master
	cd app && $(ENV_VARS) && go mod tidy
	echo "Contracts updated"

.PHONY: update-deps