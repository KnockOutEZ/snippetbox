run:
	go run ./cmd/web

run-log:
	# go run ./cmd/web >>tmp/info.log 2>>tmp/error.log
	go run ./cmd/web -log

help:
	go run ./cmd/web -help