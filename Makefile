# colors compatible settingfmt
CRED:=$(shell tput setaf 1 2>/dev/null)
CGREEN:=$(shell tput setaf 2 2>/dev/null)
CYELLOW:=$(shell tput setaf 3 2>/dev/null)
CEND:=$(shell tput sgr0 2>/dev/null)

# Build binary files
.PHONY: build
build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@mkdir -p bin
	@ret=0 && for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		go build -trimpath -o bin/$${b} $$d || ret=$$? ; \
	done ; exit $$ret
	@echo "build Success!"

# Code format
.PHONY: fmt
fmt:
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret	

# Clean up build artifacts
.PHONY: clean
clean:
	git clean -x -f
	rm -rf bin/
	rm -rf release/

# Run golang test cases
.PHONY: test
test: fmt
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	@go test -timeout 10m -race ./...
	@echo "test Success!"

# Code Coverage
# colorful coverage numerical >=90% GREEN, <80% RED, Other YELLOW
.PHONY: cover
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

.PHONY: test-cli
test-cli: build
	@echo "$(CGREEN)Run Case 1: convert mysql to markdown$(CEND)"
	@./bin/tableconvert --from mysql -t markdown --file test/mysql.txt --key value -v
	@echo "$(CGREEN)Run Case 2: convert markdown to mysql$(CEND)"
	@./bin/tableconvert --from markdown -t mysql --file test/mysql.md -v
	@echo "$(CGREEN)Run Case 3: convert mysql to csv$(CEND)"
	@./bin/tableconvert --from mysql -t csv --file test/mysql.txt --delimiter=SEMICOLON -v
	@echo "$(CGREEN)Run Case 4: convert csv to mysql$(CEND)"
	@./bin/tableconvert --from csv -t mysql --file test/mysql.csv -v
	@echo "$(CGREEN)Run Case 5: convert mysql to json$(CEND)"
	@./bin/tableconvert --from mysql -t json --file test/mysql.txt --parsing-json -v
	@echo "\n$(CGREEN)Run Case 6: convert mysql to json format:2d$(CEND)"
	@./bin/tableconvert --from mysql -t json --file test/mysql.txt --parsing-json --format=2d --minify -v
	@echo "\n$(CGREEN)Run Case 7: convert mysql to json format:column$(CEND)"
	@./bin/tableconvert --from mysql -t json --file test/mysql.txt --parsing-json --format=column --minify -v
	@echo "\n$(CGREEN)Run Case 8: convert json to mysql$(CEND)"
	@./bin/tableconvert --from json -t mysql --file test/mysql.json -v