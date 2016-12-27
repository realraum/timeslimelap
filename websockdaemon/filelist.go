package main

import (
	"encoding/json"
	"path/filepath"

	"github.com/btittelbach/pubsub"
)

type FutureFileList struct {
	backchan chan []byte
}

func ListImages(path string) []string {
	imglist, _ := filepath.Glob(filepath.Join(path, "frame*.jpg"))
	return imglist
}

func MakeTimeLapseStructJSON() []byte {
	fl := ListImages(ImagePath_)
	basefl := make([]string, len(fl))
	for idx, imgpath := range fl {
		basefl[idx] = filepath.Base(imgpath)
	}
	tld := new(wsTimeLapseData)
	tld.Path = ImageURI_
	tld.ImgList = basefl
	tld.Interval_ms = TimeLapseIntervallMS_

	replydata, err := json.Marshal(wsMessageOut{Ctx: ws_ctx_update, Data: tld})
	if err != nil {
		LogWS_.Print(err)
		return nil
	}
	return replydata
}

func GoWaitOnTrigger(ps *pubsub.PubSub) {
	shutdown_chan := ps.SubOnce("shutdown")
	triggerupdate_chan := ps.Sub(PS_TRIGGERUPDATE)
	requestlatest_chan := ps.Sub(PS_REQLATESTFILES)
	latest_json_filelist := MakeTimeLapseStructJSON()
	for {
		select {
		case <-shutdown_chan:
			ps.Unsub(triggerupdate_chan, PS_TRIGGERUPDATE)
			return
		case <-triggerupdate_chan:
			latest_json_filelist = MakeTimeLapseStructJSON()
			ps.Pub(latest_json_filelist, PS_JSONTOALL)
		case future := <-requestlatest_chan:
			switch ffl := future.(type) {
			case FutureFileList:
				ffl.backchan <- latest_json_filelist
			}
		}
	}
}
