package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/btittelbach/pubsub"
	"github.com/codegangsta/martini"
	"github.com/gorilla/websocket"
)

const (
	ws_ctx_update = "update"
)

type wsMessage struct {
	Ctx  string                 `json:"ctx"`
	Data map[string]interface{} `json:"data"`
}

type wsMessageOut struct {
	Ctx  string      `json:"ctx"`
	Data interface{} `json:"data"`
}

type wsTimeLapseData struct {
	Path        string   `json:"path"`
	ImgList     []string `json:"imglist"`
	Interval_ms uint64   `json:"interval_ms"`
}

const (
	ws_ping_period_      = time.Duration(58) * time.Second
	ws_read_timeout_     = time.Duration(70) * time.Second // must be > than ws_ping_period_
	ws_write_timeout_    = time.Duration(9) * time.Second
	ws_max_message_size_ = int64(4096)
)

var wsupgrader = websocket.Upgrader{} // use default options with Origin Check

func wsWriteMessage(ws *websocket.Conn, mt int, data []byte) error {
	ws.SetWriteDeadline(time.Now().Add(ws_write_timeout_))
	if err := ws.WriteMessage(mt, data); err != nil {
		LogWS_.Println("wsWriteTextMessage", ws.RemoteAddr(), "Error", err)
		return err
	}
	return nil
}

func goWriteToClient(ws *websocket.Conn, toclient_chan chan []byte, ps *pubsub.PubSub) {
	shutdown_c := ps.SubOnce("shutdown")
	udpate_c := ps.Sub(PS_JSONTOALL)
	ticker := time.NewTicker(ws_ping_period_)
	defer ps.Unsub(udpate_c, PS_JSONTOALL)

	inital_filelist_chan := make(chan []byte, 1)
	ps.Pub(FutureFileList{inital_filelist_chan}, PS_REQLATESTFILES)

WRITELOOP:
	for {
		var err error
		select {
		case <-shutdown_c:
			LogWS_.Println("goWriteToClient", ws.RemoteAddr(), "Shutdown")
			break WRITELOOP
		case jsonbytes, isopen := <-udpate_c:
			if !isopen {
				break WRITELOOP
			}
			err = wsWriteMessage(ws, websocket.TextMessage, jsonbytes.([]byte))
		case replybytes, isopen := <-toclient_chan:
			if !isopen {
				break WRITELOOP
			}
			err = wsWriteMessage(ws, websocket.TextMessage, replybytes)
		case initalbytes := <-inital_filelist_chan:
			if initalbytes != nil {
				err = wsWriteMessage(ws, websocket.TextMessage, initalbytes)
			}
		case <-ticker.C:
			err = wsWriteMessage(ws, websocket.PingMessage, []byte{})
		}
		if err != nil {
			LogWS_.Printf("goWriteToClient Error: %s", err)
			return
		}
	}
	wsWriteMessage(ws, websocket.CloseMessage, []byte{})
}

func convertToInt(vif interface{}) int {
	LogWS_.Printf("converToInt %+v", vif)
	switch v := vif.(type) {
	case string:
		vidint, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			LogWS_.Print(err)
			return -1
		}
		return int(vidint)
	case float64:
		return int(v)
	case int64:
		return int(v)
	case uint64:
		return int(v)
	case int:
		return v
	case uint:
		return int(v)
	default:
		return -1
	}
}

func convertToFloat32(vif interface{}) float32 {
	LogWS_.Printf("converToFloat %+v", vif)
	switch v := vif.(type) {
	case string:
		flt, err := strconv.ParseFloat(v, 32)
		if err != nil {
			LogWS_.Print(err)
			return -1
		}
		return float32(flt)
	case float32:
		return v
	case float64:
		return float32(v)
	case int64:
		return float32(v)
	case uint64:
		return float32(v)
	case int:
		return float32(v)
	case uint:
		return float32(v)
	default:
		return -1
	}
}

func goTalkWithClient(w http.ResponseWriter, r *http.Request, ps *pubsub.PubSub) {
	ws, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		LogWS_.Println(err)
		return
	}
	// client := ws.RemoteAddr()
	LogWS_.Println("Client connected", ws.RemoteAddr())

	// NOTE: no call to ws.WriteMessage in this function after this call
	// ONLY goWriteToClient writes to client
	toclient_chan := make(chan []byte, 10)
	defer close(toclient_chan)
	go goWriteToClient(ws, toclient_chan, ps)

	// logged_in := false
	ws.SetReadLimit(ws_max_message_size_)
	ws.SetReadDeadline(time.Now().Add(ws_read_timeout_))
	// the PongHandler will set the read deadline for next messages if pings arrive
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(ws_read_timeout_)); return nil })
WSREADLOOP:
	for {
		var v wsMessage
		err := ws.ReadJSON(&v)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				LogWS_.Printf("webHandleWebSocket Error: %v", err)
			}
			break WSREADLOOP
		}

		switch v.Ctx {
		case ws_ctx_update:
		}
	}
}

func RunMartini(ps *pubsub.PubSub) {
	m := martini.Classic()
	//m.Use(martini.Static("/var/lib/cloud9/static/"))
	m.Get("/websock", func(w http.ResponseWriter, r *http.Request) {
		goTalkWithClient(w, r, ps)
	})

	/*	if false {
			if err := http.ListenAndServeTLS(common.HTTPSDebugListenInterface, common.HTTPSTLSCertFilepath, common.HTTPSTLSKeyFilepath, m); err != nil {
				LogWS_.Fatal(err)
			}
		} else {
	*/if err := http.ListenAndServe("127.0.0.1:5000", m); err != nil {
		LogWS_.Fatal(err)
	}
	//	}

}
