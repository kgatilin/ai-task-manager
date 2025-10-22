.PHONY: build build-cli clean run run-cli install test-cli

build:
	@mkdir -p bin
	go build -o bin/myplugin ./cmd/myplugin
	go build -o bin/myplugin-cli ./cmd/myplugin-cli

build-cli:
	@mkdir -p bin
	go build -o bin/myplugin-cli ./cmd/myplugin-cli

clean:
	rm -rf bin/

run: build
	./bin/myplugin

run-cli: build-cli
	./bin/myplugin-cli help

install: build
	cp bin/myplugin /usr/local/bin/
	cp bin/myplugin-cli /usr/local/bin/

test-cli: build-cli
	@echo "=== Testing CLI commands ==="
	@echo "\n1. Plugin Info:"
	@./bin/myplugin-cli info
	@echo "\n2. Entity Types:"
	@./bin/myplugin-cli types
	@echo "\n3. List Items:"
	@./bin/myplugin-cli list
	@echo "\n4. Get Item:"
	@./bin/myplugin-cli get item-1
	@echo "\n5. Update Item (check stderr for event):"
	@./bin/myplugin-cli update item-1 name "Updated Name"
	@echo "\nAll CLI tests completed!"
