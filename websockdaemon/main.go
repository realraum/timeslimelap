package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/btittelbach/pubsub"
)

var (
	ps *pubsub.PubSub
)

const (
	PS_JSONTOALL      = "jsontoall"
	PS_TRIGGERUPDATE  = "triggerupdate"
	PS_REQLATESTFILES = "getlastestfiles"
)

var (
	DebugFlags_           string
	ImagePath_            string
	ImageURI_             string
	RPCSocketPath_        string
	TimeLapseIntervallMS_ uint64
)

func init() {
	flag.StringVar(&ImagePath_, "imgpath", "./timepics/", "Path to Image")
	flag.StringVar(&ImageURI_, "imguri", "/timepics/", "URI where images are exported")
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
	flag.StringVar(&RPCSocketPath_, "socketpath", "/tmp/updatetrigger.socket", "List of DebugFlags separated by ,")
	flag.Uint64Var(&TimeLapseIntervallMS_, "tlintvms", 750, "ms between frames")
}

func main() {
	flag.Parse()
	if len(DebugFlags_) > 0 {
		LogEnable(strings.Split(DebugFlags_, ",")...)
	}

	ps = pubsub.New(10)

	// call rest of main is submain func, thus give submain() defers time to do their work @exit
	MainThatReallyIsTheRealMain()

	LogMain_.Print("Exiting..")
}

func MainThatReallyIsTheRealMain() {
	//prepare clean shutdown
	defer ps.Pub(true, "shutdown")

	go StartRPCServer(ps, RPCSocketPath_) //let frontendd connect to me
	go GoWaitOnTrigger(ps)
	go RunMartini(ps)

	// wait on Ctrl-C or sigInt or sigKill
	func() {
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		fmt.Println("SIGINT received, exiting gracefully ...")
		ps.Pub(true, "shutdown")
	}()

	//return to cleanup in main-Main
}
