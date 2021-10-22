package main

import (
	//"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"

	"time"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var FinishHandler = "Pass.DealFinish"
var ConnectionHandler = "Pass.DealConnection"
var nextAddr string
var isFinished = false

type Pass struct {
	Number int
}

var Client *rpc.Client

func (p *Pass) DealConnection(p1 Pass, p2 *Pass) (err error) {
	client, err := rpc.Dial("tcp", nextAddr)
	defer client.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p1.Number, " bottles of beer on the wall, ", p1.Number, "bottles of beer. Take one down, pass it around")
	if p1.Number-1 < 1 {
		client.Call(FinishHandler, Pass{Number: 0}, &Pass{Number: 0})
		os.Exit(2)
	} else {
		client.Go(ConnectionHandler, Pass{Number: p1.Number - 1}, &Pass{Number: 0}, nil)
		return
	}

	return

}

func (p *Pass) DealFinish(p1 Pass, p2 *Pass) (err error) {
	if isFinished {
		os.Exit(2)
		return
	}
	fmt.Println("all has been taken")
	client, err := rpc.Dial("tcp", nextAddr)
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	isFinished = true
	client.Call(FinishHandler, Pass{Number: 0}, &Pass{Number: 0})
	client.Call(FinishHandler, Pass{Number: 0}, &Pass{Number: 0})

	time.Sleep(5 * time.Second)
	os.Exit(2)
	return
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8030", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")

	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	pass := new(Pass)
	rpc.Register(pass)

	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
	fmt.Println(*thisPort)
	fmt.Println(nextAddr)
	listener, err := net.Listen("tcp", ":"+*thisPort)
	if err != nil {
		fmt.Println(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	if *bottles != 0 {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("sleep for turn:", i)

		}

		go rpc.Accept(listener)

		Client, _ = rpc.Dial("tcp", nextAddr)
		Client.Call(ConnectionHandler, Pass{Number: *bottles}, &Pass{Number: 0})

	}

	rpc.Accept(listener)
}
