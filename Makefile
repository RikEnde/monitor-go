# Use this only for creating smaller binaries


static:
	go build monitor.go
	go build client.go

dynamic:
	go build -gccgoflags "-s" -compiler gccgo monitor.go
	go build -gccgoflags "-s" -compiler gccgo client.go

clean:
	@rm -f client monitor

