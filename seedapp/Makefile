# simple makefile
NAME=main
COVERAGE_MIN=$(if $(CI_COVERAGE_MIN),$(CI_COVERAGE_MIN),50)

## check: Check quality code in sonarqube.
check:
	@sonar-scanner -Dsonar.projectKey="$SONAR_PROJECT_KEY" -Dsonar.sources=. -Dsonar.host.url="$SONAR_HOST" -Dsonar.login="$SONAR_LOGIN"

## test: Run test and enforce go coverage
test:
	@go generate ./...
	@echo "Checking test coverage with minimum $(COVERAGE_MIN)%..."
	@go test -v ./... -coverprofile=cp.out -coverpkg=$$(go list ./... | grep -v -f exclude_test.txt | tr '\n' ',' | sed 's/,$$//')

## coverage: Show go coverage
coverage: test
	@echo "coverage details:";
	@echo "Package coverage breakdown:"
	@go tool cover -func=cp.out | sort -k3 -n
	@echo "-----------------------"
	@go tool cover -func=cp.out | grep total | awk -v min=$(COVERAGE_MIN) '{ \
		coverage = substr($$3, 1, length($$3)-1) + 0; \
		min += 0; \
		print "Coverage: " coverage "%" ; \
		if (coverage < min) { \
			exit 1; \
		} \
	}'

## coverage-web: Show go coverage in web
coverage-web: test
	@go tool cover -html=cp.out -o coverage.html

## bench: Run benchmark test
bench:
	go test -bench=.

## watch: development with air
watch:
	air -c .air.toml

## build: Build binary applications
build:
	@go generate ./...
	@echo building binary to ./bin/${NAME}
	@go build -o ./bin/${NAME} .

## run: Run the application
run: build
	@./bin/${NAME} service

## deploy: Deploy binary to server using ansible-playbook
deploy:
	@echo "Command to deploy script distribute atrifacts to cloud, on-prem or kubernetes clusters "

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
