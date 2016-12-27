// (c) Bernhard Tittelbach, 2015
package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/btittelbach/pubsub"
)

type RPCStruct struct {
	ps *pubsub.PubSub
}

func (r *RPCStruct) TriggerUpdate(arg bool, reply *bool) error {
	r.ps.Pub(arg, PS_TRIGGERUPDATE)
	return nil
}

func (r *RPCStruct) LEDSwitch(onoff bool, reply *bool) error {
	if onoff {
		r.ps.Pub("1", PS_SERIAL)
	} else {
		r.ps.Pub("0", PS_SERIAL)
	}
	return nil
}

func StartRPCServer(ps *pubsub.PubSub, socketpath string) {
	r := &RPCStruct{ps}
	rpc.Register(r)
	l, e := net.Listen("unixpacket", socketpath)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	rpc.Accept(l) //this blocks forever
	log.Panic("rpc socket lost")
}
