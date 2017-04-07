all: local pi test
local:
	go build -o buffalabd .
pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o buffalabd-pi
test:
	./tests

.PHONY: all local pi test
