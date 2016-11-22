package main

import (
	"ZDTalk/queue/queuehandler"
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	jsonStr := "{\"classroomId\":1,\"cmdSet\":{\"cmdDrawLine\":[{\"index\":1,\"pointAx\":101,\"pointAy\":21,\"pointBx\":72,\"pointBy\":232,\"color\":\"blue\"},{\"index\":2,\"pointAx\":23,\"pointAy\":235,\"pointBx\":59,\"pointBy\":74,\"color\":\"black\"}],\"cmdDrawText\":[{\"index\":3,\"pointX\":456,\"pointY\":353,\"text\":\"abc\",\"textWidth\":200},{\"index\":4,\"pointX\":23,\"pointY\":43,\"text\":\"abc\",\"textWidth\":200}]},\"fileUrl\":\"http://test2.clcong.com/nothing\",\"fontSize\":12,\"height\":768,\"width\":1024}"
	draw := queuehandler.DrawContent{}
	err := json.Unmarshal([]byte(jsonStr), &draw)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(draw)
	var s string = "1000000050"
	v, _ := strconv.ParseInt(s, 10, 32)
	fmt.Println(v)
}
