package main

import (
	"bufio"
	"log"
	"os"
	//"strconv"
	"time"
	
	"github.com/satori/go.uuid"
	"github.com/akkuman/parseConfig"
	jsoniter "github.com/json-iterator/go"
	"github.com/ShigemoriHakura/BilibiliDanmu"
)

func main() {
	defer func() {
		log.Println("[Main]", "请按回车关闭。。。")
		for {
			consoleReader := bufio.NewReaderSize(os.Stdin, 1)
			_, _ = consoleReader.ReadByte()
			os.Exit(0)
		}
	}()

	ConnMap.hubMap = make(map[int]*Hub)
	PhotoMap.photoMap = make(map[int64]*PhotoStruct)
	RoomMap.roomMap = make(map[int]struct{})

	log.Println("[Main]", "读取配置文件中")
	importConfig()
	log.Println("[Main]", "启动中，DanmakuChat，", BackendVersion)
	log.Println("[Main]", "头像缓存时间：", AvatarRefreshRate, "秒")
	startMessageQueue()
	startRoomQueue()
	go processMessageQueue()
	go processRoomQueue()
	go processRoomRetryQueue()
	startHttpServer()
}

func importConfig() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[Main]", "发生配置文件错误：", r)
			log.Println("[Main]", "跳过配置文件使用默认值")
		}
	}()

	var config = parseConfig.New("config.json")
	var BanWords = config.Get("BanWords").([]interface{})
	var UserMark = config.Get("UserMarks").(map[string]interface{})
	for _, v := range BanWords {
		BanString = append(BanString, v.(string))
	}
	for k, v := range UserMark {
		UserMarks[k] = v.(string)
	}

	FollowText = config.Get("FollowText").(string)
	JoinText = config.Get("JoinText").(string)
	AvatarRefreshRate = int(config.Get("AvatarRefreshRate").(float64))
}

func startMessageQueue() {
	MessageQ := initMessageQueue()
	log.Println("[Message Queue]", "初始化成功，当前队列长度：", MessageQ.Size())
}

func startRoomQueue() {
	RoomQ := initRoomQueue()
	log.Println("[Room Queue]", "初始化成功，当前队列长度：", RoomQ.Size())
}

func processMessageQueue() {
	for {
		for !MessageQ.IsEmpty() {
			tmp := MessageQ.Dequeue()
			log.Println("[Message Queue]", tmp.RoomID, "处理消息")
			ConnMap.Lock()
			connHub, ok := ConnMap.hubMap[tmp.RoomID]
			ConnMap.Unlock()
			if ok {
				json := jsoniter.ConfigCompatibleWithStandardLibrary
				ddata, err := json.Marshal(tmp.Data)
				if err == nil {
					//log.Println("Sent: ", 1, string(ddata))
					connHub.broadcast <- ddata
					//time.Sleep(time.Duration(80) * time.Millisecond)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func processRoomQueue() {
	for {
		for !RoomQ.IsEmpty() {
			tmp := RoomQ.Dequeue()
			log.Println("[Room Queue]", tmp.RoomID, "处理房间")
			RoomMap.Lock()
			_, ok := RoomMap.roomMap[tmp.RoomID]
			RoomMap.Unlock()
			if !ok {
				log.Println("[Room Queue]", tmp.RoomID, "建立WS链接")
				go startBWS(tmp.RoomID)
			} else {
				log.Println("[Room Queue]", tmp.RoomID, "已存在，不新建")
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func processRoomRetryQueue() {
	for {
		time.Sleep(10 * time.Second)
		//log.Println("[Room Retry Queue]", "检查存在Hub但是不存在弹幕服务的房间")
		ConnMap.Lock()
		for _, v := range ConnMap.hubMap {
			//log.Println("[Room Retry Queue]", "检查", v.roomId)
			RoomMap.Lock()
			if _, ok := RoomMap.roomMap[v.roomId]; !ok {
				log.Println("[Room Retry Queue]", v.roomId, "建立WS链接")
				RoomMap.roomMap[v.roomId] = struct{}{}
				go startBWS(v.roomId)
			}
			RoomMap.Unlock()
		}
		ConnMap.Unlock()
		//log.Println("[Room Retry Queue]", "检查完成")
	}
}

func startBWS(roomID int) {
	RoomMap.Lock()
	RoomMap.roomMap[roomID] = struct{}{}
	RoomMap.Unlock()
	defer func() {
		log.Println("[Danmaku]", roomID, "结束")
		RoomMap.Lock()
		delete(RoomMap.roomMap, roomID)
		RoomMap.Unlock()
	}()
	log.Println("[Danmaku]", roomID, "WS监听服务启动中")

	c, err := BilibiliDanmu.NewClient(uint32(roomID))
	if err != nil {
		log.Println("models.NewClient err: ", err)
		return
	}
	pool, err := c.Start();
	if err != nil {
		log.Println("c.Start err :", err)
		return
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	for {
		u1 := uuid.Must(uuid.NewV4())
		select {
		case uc := <-pool.MsgUncompressed:
			// 目前只处理未压缩数据的关注数变化信息
			if cmd := json.Get([]byte(uc), "cmd").ToString(); BilibiliDanmu.CMD(cmd) == BilibiliDanmu.CMDRoomRealTimeMessageUpdate {
				fans := json.Get([]byte(uc), "data", "fans").ToInt()
				log.Println("当前房间关注数变动：", fans)
			}
		case src := <-pool.UserMsg:
			m := BilibiliDanmu.NewDanmu()
			m.GetDanmuMsg([]byte(src))
			//log.Println(src)
			if !checkComments(m.Text) {
				var Medal = new(MedalInfo)
				Medal.UperID = int64(m.MedalUperId)
				Medal.UserID = int64(m.UID)
				Medal.ClubName = m.MedalName
				Medal.Level = int(m.MedalLevel)

				var data = new(dataUserStruct)
				data.Cmd = 2
				data.Data.Id = u1.String()
				data.Data.UserId = int64(m.UID)
				data.Data.AvatarUrl = ""
				data.Data.Timestamp = time.Now().Unix()
				data.Data.AuthorName = m.Uname
				data.Data.AuthorType = 0
				data.Data.PrivilegeType = 0
				data.Data.Content = m.Text
				data.Data.UserMark = getUserMark(int64(m.UID))
				data.Data.Medal = *Medal
				var dataQ = new(Message)
				dataQ.RoomID = roomID
				dataQ.Data = data
				MessageQ.Enqueue(dataQ)
			}
			log.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.Ulevel, m.Uname, m.Text)

		case src := <-pool.UserGift:
			g := BilibiliDanmu.NewGift()
			g.GetGiftMsg([]byte(src))

			log.Println(src)
			var data = new(dataGiftStruct)
			data.Cmd = 3
			data.Data.Id = u1.String()
			//data.Data.UserId = int64(g.UID)
			data.Data.AvatarUrl = ""
			data.Data.Timestamp = time.Now().Unix()
			data.Data.AuthorName = g.UUname
			data.Data.AuthorType = 0
			//data.Data.UserMark = getUserMark(d.UserID)
			data.Data.GiftName = g.GiftName
			data.Data.TotalCoin = int(g.Price)
			var dataQ = new(Message)
			dataQ.RoomID = roomID
			dataQ.Data = data
			MessageQ.Enqueue(dataQ)
			//log.Println("Conn Gift", data)

			log.Printf("%s %s 价值 %d 的 %s\n", g.UUname, g.Action, g.Price, g.GiftName)
		case src := <-pool.UserEnter:
			name := json.Get([]byte(src), "data", "uname").ToString()
			log.Printf("欢迎VIP %s 进入直播间", name)
		case src := <-pool.UserGuard:
			name := json.Get([]byte(src), "data", "username").ToString()
			log.Printf("欢迎房管 %s 进入直播间", name)
		case src := <-pool.UserEntry:
			cw := json.Get([]byte(src), "data", "copy_writing").ToString()
			log.Printf("%s", cw)
		}
	}
}
