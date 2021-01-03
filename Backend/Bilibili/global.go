package main

import (
	"sync"
)

const defaultAvatar = "https://i0.hdslb.com/bfs/face/168aea8d4b02b6cd2235c11a0394a7d25f7350bd.jpg@42w_42h_1c.webp"

var BackendVersion = "0.0.1"
var EnableTranslate = false
var FollowText = "关注了主播"
var JoinText = "加入直播间"
var AvatarRefreshRate = 86400
var BanString []string
var UserMarks = make(map[string]string)

var ConnMap struct {
	sync.Mutex
	hubMap map[int]*Hub
}

var PhotoMap struct {
	sync.Mutex
	photoMap map[int64]*PhotoStruct
}

var RoomMap struct {
	sync.Mutex
	roomMap map[int]struct{}
}

var MessageQ MessageQueue
var RoomQ MessageQueue
