set:
	go build
	rm data.json
	./kv set NODE_ENV production
	./kv set testkey hello
	./kv set aws_key ssiikkuu_uuyy13
	./kv set jira_url https://fake_org.atlassian.net
	./kv set jira_user user@user.com
	./kv set jira_token aGVsbG8=
	cat data.json | jq
getall:
	go build
	./kv get --all
getone:
	go build
	./kv get NODE_ENV
run:
	go run .

release:
	go build -ldflags "-w -s"
