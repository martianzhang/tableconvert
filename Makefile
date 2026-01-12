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

# Unified colored echo helpers per platform
ifeq ($(IS_WINDOWS),)
define ECHO_GREEN
	@printf "%s%s%s\n" "$(CGREEN)" "$(1)" "$(CEND)"
endef
define ECHO_YELLOW
	@printf "%s%s%s\n" "$(CYELLOW)" "$(1)" "$(CEND)"
endef
define ECHO_RED
	@printf "%s%s%s\n" "$(CRED)" "$(1)" "$(CEND)"
endef
else
define ECHO_GREEN
	@powershell -NoProfile -Command "Write-Host '$(1)' -ForegroundColor Green"
endef
define ECHO_YELLOW
	@powershell -NoProfile -Command "Write-Host '$(1)' -ForegroundColor Yellow"
endef
define ECHO_RED
	@powershell -NoProfile -Command "Write-Host '$(1)' -ForegroundColor Red"
endef
endif

ifeq ($(IS_WINDOWS),)
	# Unix-like commands
	BUILD_CMD = mkdir -p bin && go build -trimpath -o bin/tableconvert ./cmd/tableconvert
	CLEAN_CMD = rm -rf bin/ release/ feature/
	COVER_SUMMARY_CMD = tail -n 1 test/coverage.txt | awk '{sub(/%/, "", $$NF); if($$NF < 70) {print "$(CRED)"$$0"%$(CEND)"} else if ($$NF >= 80) {print "$(CGREEN)"$$0"%$(CEND)"} else {print "$(CYELLOW)"$$0"%$(CEND)"}}'
else
	# Windows: run via PowerShell so commands work in PowerShell/cmd environments
	BUILD_CMD = powershell -NoProfile -Command "& { if (-not (Test-Path -Path bin)) { New-Item -ItemType Directory -Path bin | Out-Null }; go build -trimpath -o bin/tableconvert.exe ./cmd/tableconvert }"
	CLEAN_CMD = powershell -NoProfile -Command "Remove-Item -Recurse -Force bin,release,feature -ErrorAction SilentlyContinue"
	# Use Write-Host -ForegroundColor for colored coverage output
	COVER_SUMMARY_CMD = powershell -NoProfile -Command "$$line = (Get-Content test/coverage.txt | Select-Object -Last 1); if ($$line -match '(\d+(\.\d+)?)%') { $$p=[double]$$Matches[1]; if($$p -lt 80){ Write-Host $$line -ForegroundColor Red } elseif($$p -ge 90){ Write-Host $$line -ForegroundColor Green } else { Write-Host $$line -ForegroundColor Yellow } } else { Write-Host $$line }"
endif

.PHONY: all build fmt clean test cover test-cli release release-clean release-checksums release-zip release-notes

all: build

# Build binary files
build: fmt
	$(call ECHO_GREEN,Building ...)
	@$(BUILD_CMD)
	@echo "build success!"

# Release management - builds all platform binaries and creates checksums
release: release-clean
	$(call ECHO_GREEN,Starting release build...)
	@echo "Building binaries for all platforms..."

	@# Linux
	$(call ECHO_YELLOW,Building Linux amd64...)
	@GOOS=linux GOARCH=amd64 go build -trimpath -o release/tableconvert-linux-amd64 ./cmd/tableconvert
	$(call ECHO_YELLOW,Building Linux arm64...)
	@GOOS=linux GOARCH=arm64 go build -trimpath -o release/tableconvert-linux-arm64 ./cmd/tableconvert

	@# macOS (Darwin)
	$(call ECHO_YELLOW,Building macOS amd64...)
	@GOOS=darwin GOARCH=amd64 go build -trimpath -o release/tableconvert-darwin-amd64 ./cmd/tableconvert
	$(call ECHO_YELLOW,Building macOS arm64...)
	@GOOS=darwin GOARCH=arm64 go build -trimpath -o release/tableconvert-darwin-arm64 ./cmd/tableconvert

	@# Windows
	$(call ECHO_YELLOW,Building Windows amd64...)
	@GOOS=windows GOARCH=amd64 go build -trimpath -o release/tableconvert-windows-amd64.exe ./cmd/tableconvert

	@# Generate checksums
	$(call ECHO_GREEN,Generating checksums...)
	@cd release && sha256sum tableconvert-* > checksums.txt

	@# Create release info
	$(call ECHO_GREEN,Creating release info...)
	@echo "Version: $(shell git describe --tags --always 2>/dev/null || echo 'dev')" > release/RELEASE_INFO.txt
	@echo "Built: $(shell date -u +'%Y-%m-%d %H:%M:%S UTC')" >> release/RELEASE_INFO.txt
	@echo "Go: $(shell go version)" >> release/RELEASE_INFO.txt

	@$(call ECHO_GREEN,Release build complete!)
	@echo ""
	@echo "Release artifacts in ./release/:"
	@ls -lh release/
	@echo ""
	@echo "Checksums:"
	@cat release/checksums.txt

# Clean release directory
release-clean:
	@echo "Cleaning release directory..."
	@rm -rf release/

# Just generate checksums for existing binaries
release-checksums:
	$(call ECHO_GREEN,Generating checksums for existing binaries...)
	@cd release && sha256sum tableconvert-* > checksums.txt
	@cat release/checksums.txt

# Create zip archives for each platform (requires zip command)
release-zip: release
	$(call ECHO_GREEN,Creating zip archives...)
	@if ! command -v zip &> /dev/null; then \
		$(call ECHO_RED,Error: 'zip' command not found. Please install zip to create archives.); \
		echo "On Ubuntu/Debian: sudo apt-get install zip"; \
		echo "On macOS: brew install zip"; \
		echo "On Windows: Use WSL or install zip via Chocolatey"; \
		exit 1; \
	fi
	@cd release && \
		for platform in linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64; do \
			if [ -f "tableconvert-$${platform}" ] || [ -f "tableconvert-$${platform}.exe" ]; then \
				echo "Zipping $${platform}..."; \
				zip -q tableconvert-$${platform}.zip tableconvert-$${platform}* checksums.txt RELEASE_INFO.txt 2>/dev/null || true; \
			fi; \
		done
	@$(call ECHO_GREEN,Zip archives created in release/)

# Generate release notes from git history
release-notes:
	$(call ECHO_GREEN,Generating release notes...)
	@echo "# Release Notes" > release/RELEASE_NOTES.md
	@echo "" >> release/RELEASE_NOTES.md
	@echo "## Changes since last tag" >> release/RELEASE_NOTES.md
	@git log --oneline --decorate --no-merges -n 20 2>/dev/null || echo "No git history available" >> release/RELEASE_NOTES.md
	@echo "" >> release/RELEASE_NOTES.md
	@echo "## Build Info" >> release/RELEASE_NOTES.md
	@echo "- Built: $(shell date -u +'%Y-%m-%d %H:%M:%S UTC')" >> release/RELEASE_NOTES.md
	@echo "- Go: $(shell go version | cut -d' ' -f3)" >> release/RELEASE_NOTES.md
	@echo "" >> release/RELEASE_NOTES.md
	@echo "## Platform Support" >> release/RELEASE_NOTES.md
	@echo "- Linux (x64, ARM64)" >> release/RELEASE_NOTES.md
	@echo "- macOS (x64, ARM64)" >> release/RELEASE_NOTES.md
	@echo "- Windows (x64)" >> release/RELEASE_NOTES.md
	@cat release/RELEASE_NOTES.md

# Code format
fmt:
	$(call ECHO_GREEN,Code formatting...)
	@go fmt ./...

# Clean up build artifacts
clean:
	@echo "Cleaning..."
	@$(CLEAN_CMD)

# Run golang test cases
test: fmt
	$(call ECHO_GREEN,Run all test cases ...)
	@go test -timeout 10m ./...
	@echo "test Success!"

# Code Coverage
# colorful coverage numerical >=90% GREEN, <80% RED, Other YELLOW
cover: test
	$(call ECHO_GREEN,Run test cover check ...)
	@go test ./... -coverprofile=test/coverage.data
	@go tool cover -html=test/coverage.data -o test/coverage.html
	@go tool cover -func=test/coverage.data -o test/coverage.txt
	@$(COVER_SUMMARY_CMD)

# Run random output test cases, human check result
test-cli: build
	@echo ""
	$(call ECHO_GREEN,Run Case 1: convert json to mysql)
	@./bin/tableconvert --from json -t mysql --file test/mysql.json -v
	@echo ""
	$(call ECHO_GREEN,Run Case 2: convert mysql to xlsx)
	@./bin/tableconvert --from mysql -t xlsx --file test/mysql.txt --result test/mysql.xlsx -v
	@./bin/tableconvert --from xlsx -t mysql --file test/mysql.xlsx
	@echo ""
	$(call ECHO_GREEN,Run Case 3: convert mysql use template)
	@./bin/tableconvert --from mysql -t template --file test/mysql.txt --template test/jsonlines.tmpl -v
