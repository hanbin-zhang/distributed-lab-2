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

var PassHandler = "Pass.DealPass"
var FinishHandler = "Pass.DealFinish"
var nextAddr string

type Pass struct {
	Number int
}

var Client *rpc.Client

func (p *Pass) DealPass(p1 Pass, p2 *Pass) (err error) {
	if p1.Number > 0 {
		fmt.Println(p1.Number, " bottles of beer on the wall,", p1.Number, " bottles of beer. Take one down, Pass it around...")
		e1 := Client.Call(PassHandler, Pass{Number: p1.Number - 1}, nil)

		if e1 != nil {
			fmt.Println(e1)
			return
		}
	} else {
		fmt.Println("all has been taken")
		Client.Call(FinishHandler, nil, nil)
		os.Exit(2)
	}
	return
}
func (p *Pass) DealFinish(p1 Pass, p2 *Pass) (err error) {
	fmt.Println("all has been taken")
	os.Exit(2)
	return
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8030", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	target := flag.String("target", "127.0.0.1:8030", "the ip address for the target process")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	pass := new(Pass)
	rpc.Register(pass)
	//err := rpc.Register(&PassOperation{})
	//if err != nil {
	//	fmt.Println(err)
	//		return
	//	}

	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
	fmt.Println(*thisPort)
	fmt.Println(*target)
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(listener)

	if *bottles != 0 {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("sleep for turn:", i)

		}
		fmt.Println(*target)

		Client, _ = rpc.Dial("tcp", *target)

	}
	fmt.Println("aaa")
	rpc.Accept(listener)
	fmt.Println("bbb")
	if *bottles == 0 {
		Client, _ = rpc.Dial("tcp", *target)
		defer func(client *rpc.Client) {
			err := client.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(Client)
	}

	if *bottles != 0 {
		fmt.Println("11331qwdfAFWAS")
		errC := Client.Call(PassHandler, Pass{Number: 0}, nil)
		if errC != nil {
			fmt.Println(errC)
			return
		}
	}
}
