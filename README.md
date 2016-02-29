# monitor-go

A simple client server app written in google go. This project serves no purpose other than to have a bit of everything in it: some file handling, string parsing, networking, database, concurrency, marshalling to json, basic authentication etc, while trying out golang. 

The server reads cpu and kernel status from /proc/stat and /proc/cpuinfo and exposes this information with a simple REST service over https

Services available:

/					Welcome page with some examples

/cpu				Contents of /proc/cpuinfo 

/cpu/0/bogomips		Bogomips of CPU core 0


/stat				Contents of /proc/stat

/stat/cpu			See http://www.linuxhowtos.org/System/procstat.htm for the meaning of these numbers


/history			Returns the deltas per interval (10 seconds) of /stat/cpu over the past hour
					You could use this to draw a graph of cpu usage. The interesting fields will usually
					be user, idle and iowait. Temperature and cpu frequency are not implemented yet. 


The server listens to commands posted in json form: {"please":"command"}

Currently the commands Die! and Restart are supported. The only difference is that the process exits with status code 0 or 1, which a shell script could use to decide to restart the process or not. 


# Usage 

Server 

 go run monitor.go -port 8888 

Client

 go run client -host localhost -port 8888 -user username -password password cpu 0

Or run ./bin/start to start with default settings and ./bin/kill to stop 


# Configuration

User accounts are in a file in json format under $GOROOT/users/users in the form {username : hash of password}
Example: {"santa":"a0121a95b0bea2bf240dcbbcea9abfc12ffda9fb"}

Generate a hash with tool.go 

Database credentials are expected to be in json format in a file in the same directory called postgres. 
Example: {"user":"dbuser", "password":"dbpassword", "database":"dbname", "host":"localhost"}


# Makefile 

Run make static to generate statically linked binaries (default). These can be very large (several mb). Run make dynamic to generate dynamically linked binaries. Note that there are no performance benefits to this, it's purely to satisfy the feeling that such a simple program shouldn't be 9 Mb. 


# TODO

- Implement remote restart 
- Implement system memory status and tracking

