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
	op_trigger_  bool
	op_led_      string
)

const DEFAULT_TUER_DOORCMD_SOCKETPATH string = "/run/tuer/door_cmd.unixpacket"

func init() {
	flag.StringVar(&socket_path_, "socketpath", "/tmp/updatetrigger.socket", "rpc command socket path")
	flag.BoolVar(&op_trigger_, "updatefilelist", false, "updatefilelist update trigger")
	flag.StringVar(&op_led_, "led", "", "led on/off")

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
	if op_trigger_ {
		if err := rpcc.Call("RPCStruct.TriggerUpdate", true, &reply); err == nil {
			fmt.Println("websockdaemon triggered")
		} else {
			fmt.Println("!!!", err.Error())
		}
	}
	switch op_led_ {
	case "on", "ON", "1", "On":
		if err := rpcc.Call("RPCStruct.LEDSwitch", true, &reply); err == nil {
			fmt.Println("LED on")
		} else {
			fmt.Println("!!!", err.Error())
		}
	case "off", "OFF", "0", "Off":
		if err := rpcc.Call("RPCStruct.LEDSwitch", false, &reply); err == nil {
			fmt.Println("LED on")
		} else {
			fmt.Println("!!!", err.Error())
		}
	default:
	}
}
