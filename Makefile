INPUT_DIR := example/enumgo/input
OUTPUT_DIR := example/enumgo/generated

tidy:
	@go mod tidy

test:
	@clear
	@gotestsum --junitfile-hide-empty-pkg --format testname

lint:
	golangci-lint run

release:
	goreleaser release --clean --config=./.goreleaser.yml

gen_example: clean_example
	@echo "Generating example..."
	@mkdir -p $(OUTPUT_DIR)
	@go run . enumgo -i $(INPUT_DIR) -o $(OUTPUT_DIR) -p example
	@gofmt -w .

clean_example:
	@rm -rf $(OUTPUT_DIR)/*
