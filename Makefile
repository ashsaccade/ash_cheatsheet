run:
	go run cmd/app/main.go

sandbox:
	go run cmd/sandbox/main.go

peg:
	peg -switch -inline -strict -output internal/grammar/grammar.go internal/grammar/grammar.peg