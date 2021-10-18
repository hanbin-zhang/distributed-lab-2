package bottlesofbeer

import (
	//"bufio"
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var PassHandler = "passOperation.DealPass"
var FinishHandler = " passOperation.DealFinish"
var nextAddr string

type passOperation struct{ targetIP string }
type pass struct {
	Number int
}

func (p *passOperation) DealPass(p1 pass) (err error) {
	if p1.Number > 0 {
		fmt.Println(p1.Number, " bottles of beer on the wall,", p1.Number, " bottles of beer. Take one down, pass it around...")
		e1 := client.Call(PassHandler, pass{Number: p1.Number - 1}, nil)

		if e1 != nil {
			fmt.Println(e1)
			return
		}
	} else {
		fmt.Println("all has been taken")
		client.Call(FinishHandler, nil, nil)
		os.Exit(2)
	}
	return
}
func (p *passOperation) DealFinish(p1 pass) (err error) {
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
	err := rpc.Register(&passOperation{targetIP: *target})
	if err != nil {
		fmt.Println(err)
		return
	}

	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.

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
			client, _ = rpc.Dial("tcp", *target)
		}
	}

	rpc.Accept(listener)

	if *bottles == 0 {
		client, _ = rpc.Dial("tcp", *target)
		defer func(client *rpc.Client) {
			err := client.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(client)
	}

	if *bottles != 0 {
		errC := client.Call(PassHandler, pass{Number: 0}, nil)
		if errC != nil {
			fmt.Println(err)
			return
		}
	}
}

var client *rpc.Client
