
run_ping:
	go run cmd/main.go
run_race_ping:
	go run -race cmd/main.go

run_udp_echo_server:
	go run udpserver/udpserver.go --addr=0.0.0.0:9000
run_race_udp_echo_server:
	go run -race udpserver/udpserver.go --addr=0.0.0.0:9000

build_udp_server:
	go build -o udpserver udpserver/udpserver.go