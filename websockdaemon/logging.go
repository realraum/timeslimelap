// (c) Bernhard Tittelbach, 2013

package main

import "os"
import "log"

type NullWriter struct{}

func (n *NullWriter) Write(p []byte) (int, error) { return len(p), nil }

var (
	LogMain_ *log.Logger
	LogWS_   *log.Logger
	LogRPC_  *log.Logger
)

func init() {
	LogMain_ = log.New(&NullWriter{}, "", 0)
	LogWS_ = log.New(&NullWriter{}, "", 0)
	LogRPC_ = log.New(&NullWriter{}, "", 0)
}

func LogEnable(logtypes ...string) {
	for _, logtype := range logtypes {
		switch logtype {
		case "MAIN":
			LogMain_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "WS":
			LogWS_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "RPC":
			LogRPC_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "ALL":
			LogMain_ = log.New(os.Stderr, "MAIN"+" ", log.LstdFlags)
			LogWS_ = log.New(os.Stderr, "WS"+" ", log.LstdFlags)
			LogRPC_ = log.New(os.Stderr, "RPC"+" ", log.LstdFlags)
		}
	}
}
