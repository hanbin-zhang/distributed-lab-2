package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, o := rpc.Dial("tcp", *server)
	defer func(client *rpc.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(client)
	dat, _ := os.Open("wordlist")
	datReader := bufio.NewReader(dat)
LOOP:
	for {
		readString, err := datReader.ReadString('\n')
		if err != nil {
			break LOOP
		}
		request := stubs.Request{Message: readString}
		response := new(stubs.Response)
		errClientCall := client.Call(stubs.PremiumReverseHandler, request, response)
		if errClientCall != nil {
			fmt.Println(errClientCall)
			return
		}
		fmt.Println("Responded:" + response.Message)
	}

	//TODO: connect to the RPC server and send the request(s)
}
