// (c) Bernhard Tittelbach, 2016
package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
)

var (
	socket_path_ string
)

const DEFAULT_TUER_DOORCMD_SOCKETPATH string = "/run/tuer/door_cmd.unixpacket"

func init() {
	flag.StringVar(&socket_path_, "socketpath", "/tmp/updatetrigger.socket", "rpc command socket path")
	flag.Parse()
}

func ConnectToRPCServer(socketpath string) (c *rpc.Client) {
	var err error
	c, err = rpc.Dial("unixpacket", socketpath)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	rpcc := ConnectToRPCServer(socket_path_)
	var reply bool
	if err := rpcc.Call("RPCStruct.TriggerUpdate", true, &reply); err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("!!!", err.Error())
	}
}
