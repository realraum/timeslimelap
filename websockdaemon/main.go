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
	PS_SERIAL         = "toserial"
)

var (
	DebugFlags_           string
	ImagePath_            string
	ImageURI_             string
	RPCSocketPath_        string
	TimeLapseIntervallMS_ uint64
	LedFlashTTYDev_       string
)

func init() {
	flag.StringVar(&ImagePath_, "imgpath", "./timepics/", "Path to Image")
	flag.StringVar(&ImageURI_, "imguri", "/timepics/", "URI where images are exported")
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
	flag.StringVar(&RPCSocketPath_, "socketpath", "/tmp/updatetrigger.socket", "List of DebugFlags separated by ,")
	flag.Uint64Var(&TimeLapseIntervallMS_, "tlintvms", 750, "ms between frames")
	flag.StringVar(&LedFlashTTYDev_, "leddev", "/dev/ttyUSB0", "Led Flash Serial Device")
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

func GoForwardToSerial(ps *pubsub.PubSub, serwr_chan chan string) {
	shutdown_chan := ps.SubOnce("shutdown")
	serpub_chan := ps.Sub(PS_SERIAL)
	for {
		select {
		case <-shutdown_chan:
			ps.Unsub(serpub_chan, PS_SERIAL)
			return
		case toser := <-serpub_chan:
			serwr_chan <- toser.(string)
		}
	}
}

func MainThatReallyIsTheRealMain() {
	//prepare clean shutdown
	defer ps.Pub(true, "shutdown")

	serwr_chan, _, err := OpenAndHandleSerial(LedFlashTTYDev_, 57600)
	if err != nil {
		LogMain_.Printf("OpenAndHandleSerial Error: %s", err)
	} else {
		go GoForwardToSerial(ps, serwr_chan)
	}
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
