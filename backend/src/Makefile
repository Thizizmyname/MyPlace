test:
	go test data handler myplaceutils requests_responses
	make clean

clean:
	rm -f data/rooms data/users rooms users server clients_simulator/clients_simulator

run:
	go build server.go
	./server

testrun:
	go build server.go
	./server .
