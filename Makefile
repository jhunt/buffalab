all: local pi
local:
	go build -o buffalabd .
pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o buffalabd-pi

.PHONY: all local pi
