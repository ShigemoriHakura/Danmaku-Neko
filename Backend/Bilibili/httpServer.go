package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func startHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r)
	})
	r.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r)
	})
	r.HandleFunc("/room/{key}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	r.HandleFunc("/stylegen", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	r.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	r.HandleFunc("/server_info", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version": "` + BackendVersion + `", "config": {"enableTranslate": ` + strconv.FormatBool(EnableTranslate) + `}}`))
	})
	r.HandleFunc("/room_info", func(w http.ResponseWriter, r *http.Request) {
		var roomStr = ""
		var i = 0
		ConnMap.Lock()
		for _, v := range ConnMap.hubMap {
			roomStr += strconv.Itoa(v.roomId) + ","
			i ++
		}
		ConnMap.Unlock()
		roomStr = trimLastChar(roomStr)
		w.Write([]byte(`{"version": "` + BackendVersion + `", "rooms": "` + roomStr + `", "roomCount": "` + strconv.Itoa(i) + `"}`))
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))
	http.Handle("/", r)

	log.Println("[Main]", "等待用户连接")
	err := http.ListenAndServe("0.0.0.0:23456", nil)
	if err != nil {
		log.Println("[Main]", "发生主端口监听错误: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WS Home]", "发生处理错误: ", err)
	} else {
		log.Println("[WS Home]", "新的前端WS连接：", fmt.Sprintf("%s", conn.RemoteAddr().String()))
		go serveWS(conn)
	}
}

func serveWS(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WS Server]", "发生连接错误: ", err)
			conn.Close()
			break
		} else {
			any := jsoniter.Get(msg)
			var cmd = any.Get("cmd").ToString()
			//log.Println("Conn: ", 1, string(msg))
			//log.Println("Conn cmd: ", cmd)
			switch cmd {
			case "0":
				conn.WriteMessage(1, []byte(`{}`))
				break
			case "1":
				var roomID = any.Get("data", "roomId").ToInt()
				var roomIDString = any.Get("data", "roomId").ToString()
				var isfirstLoad = any.Get("data", "isfirstLoad").ToBool()
				var FrontendV = any.Get("data", "version").ToString()
				log.Println("[WS Server]", "请求房间ID：", roomID)
				ConnMap.Lock()
				ConnM, ok := ConnMap.hubMap[roomID]
				ConnMap.Unlock()
				if !ok {
					var data = new(Message)
					data.RoomID = roomID
					RoomQ.Enqueue(data)
					ConnMap.Lock()
					ConnMap.hubMap[roomID] = newHub()
					ConnMap.hubMap[roomID].roomId = roomID
					go ConnMap.hubMap[roomID].run()
					ConnMap.Unlock()
					conn.Close()
					return
				}
				if(isfirstLoad){
					conn.WriteMessage(1, []byte(`{"cmd":2,"data":{"id":0,"avatarUrl":"https://i0.hdslb.com/bfs/face/168aea8d4b02b6cd2235c11a0394a7d25f7350bd.jpg@42w_42h_1c.webp","timestamp":1601641021,"authorName":"弹幕姬","authorType":0,"privilegeType":0,"translation":"","content":"房间 (` + roomIDString + `) - (` + FrontendV + `) 连接成功~","userMark":"","medalInfo":{"UperID":0,"ClubName":"","Level":0}}}`))			
				}
				client := &Client{hub: ConnM, conn: conn, send: make(chan []byte, 8192)}
				client.hub.register <- client
				go client.readPump()
				return
			}
		}
	}
}

func parseVersion(s string, width int) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", width)
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		fmt.Printf("[Parse Version] parseVersion(%s): error=%s\n", s, err);
		return 0
	}
	//fmt.Printf("parseVersion: [%s] => [%s] => [%d]\n", s, v, result);
	return result;
}