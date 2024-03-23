set:
	go build
	rm data.json
	./kv set NODE_ENV production
	./kv set testkey hello
	cat data.json | jq
getall:
	go build
	./kv get
getone:
	go build
	./kv get NODE_ENV
run:
	go run .

release:
	go build -ldflags "-w -s"
