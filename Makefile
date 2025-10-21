## Cross-platform Makefile
## Works on Linux/macOS (POSIX shell) and Windows PowerShell (user's default shell is PowerShell)

# Detect GOOS to switch behaviors for Windows vs Unix-like
GOOS := $(shell go env GOOS)

# Color helpers: only enabled on Unix-like systems where tput exists
ifeq ($(OS),Windows_NT)
	# Running on Windows (cmd/powershell)
	CRED :=
	CGREEN :=
	CYELLOW :=
	CEND :=
else
	CRED := $(shell tput setaf 1 2>/dev/null || echo)
	CGREEN := $(shell tput setaf 2 2>/dev/null || echo)
	CYELLOW := $(shell tput setaf 3 2>/dev/null || echo)
	CEND := $(shell tput sgr0 2>/dev/null || echo)
endif

# Determine if platform is Windows (either GOOS or OS indicates Windows)
IS_WINDOWS := $(strip $(filter Windows_NT windows,$(OS) $(GOOS)))

ifeq ($(IS_WINDOWS),)
	# Unix-like commands
	BUILD_CMD = mkdir -p bin && go build -trimpath -o bin/tableconvert ./cmd/tableconvert
	CLEAN_CMD = rm -rf bin/ release/ feature/
else
	# Windows: run via PowerShell so commands work in PowerShell/cmd environments
	BUILD_CMD = powershell -NoProfile -Command "& { if (-not (Test-Path -Path bin)) { New-Item -ItemType Directory -Path bin | Out-Null }; go build -trimpath -o bin/tableconvert.exe ./cmd/tableconvert }"
	CLEAN_CMD = powershell -NoProfile -Command "Remove-Item -Recurse -Force bin,release,feature -ErrorAction SilentlyContinue"
endif

.PHONY: all build fmt clean test cover test-cli

all: build

# Build binary files
build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@$(BUILD_CMD)
	@echo "build success!"

# Code format
fmt:
	@echo "$(CGREEN)Code formatting...$(CEND)"
	@go fmt ./...

# Clean up build artifacts
clean:
	@echo "Cleaning..."
	@$(CLEAN_CMD)

# Run golang test cases
test: fmt
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	@go test -timeout 10m ./...
	@echo "test Success!"

# Code Coverage
# colorful coverage numerical >=90% GREEN, <80% RED, Other YELLOW
cover: test
	@echo "$(CGREEN)Run test cover check ...$(CEND)"
	@go test ./... -coverprofile=test/coverage.data
	@go tool cover -html=test/coverage.data -o test/coverage.html
	@go tool cover -func=test/coverage.data -o test/coverage.txt
	@tail -n 1 test/coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "$(CRED)"$$0"%$(CEND)"} \
		else if ($$NF >= 90) \
			{print "$(CGREEN)"$$0"%$(CEND)"} \
		else \
			{print "$(CYELLOW)"$$0"%$(CEND)"}}'

# Run random output test cases, human check result
test-cli: build
	@echo "\n$(CGREEN)Run Case 1: convert json to mysql$(CEND)"
	@./bin/tableconvert --from json -t mysql --file test/mysql.json -v
	@echo "\n$(CGREEN)Run Case 2: convert mysql to xlsx$(CEND)"
	@./bin/tableconvert --from mysql -t xlsx --file test/mysql.txt --result test/mysql.xlsx -v
	@./bin/tableconvert --from xlsx -t mysql --file test/mysql.xlsx
	@echo "\n$(CGREEN)Run Case 3: convert mysql use template$(CEND)"
	@./bin/tableconvert --from mysql -t template --file test/mysql.txt --template test/jsonlines.tmpl -v