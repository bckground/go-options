test:
	go clean -i .
	go generate .
	go install .
	go generate ./...
	go test -tags testing ./...
	diff test/config_options.go test/golden/config_options.go.txt
	diff test/configWithNoError_options.go test/golden/configWithNoError_options.go.txt
	diff test/configWithBuild_options.go test/golden/configWithBuild_options.go.txt

generate:
	go generate .

golden:
	mkdir -p test/golden
	for file in test/*_options.go; do \
		cp "$$file" "test/golden/$$(basename $$file).txt"; \
	done

.PHONY: golden test

