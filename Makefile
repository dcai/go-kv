BIN=kv

build:
	@go build -o $(BIN)

test-set: build
	@rm data.json
	./kv set NODE_ENV production
	./kv set testkey hello
	./kv set aws_key ssiikkuu_uuyy13
	./kv set jira_url https://fake_org.atlassian.net
	./kv set jira_user user@user.com
	./kv set jira_token aGVsbG8=
	@echo "=========================================="
	cat data.json | jq

test-getall: build
	@./kv get --all

test-getone: build
	@./kv get NODE_ENV

run: build
	go run .

release:
	go build -ldflags "-w -s"
