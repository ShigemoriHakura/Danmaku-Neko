<template>
  <component v-bind:is="currentRenderer" :config="config"></component>
</template>

<script>
//引入核心
import { Base64 } from 'js-base64';
import { v4 as uuidv4 } from 'uuid';

//引入组件
import * as chatConfig from '@/config/chatConfig.js'
import ChatRoom from '@/components/Renderer/ChatRoom'
import FullScreen from '@/components/Renderer/FullScreen'
import YoutubeChat from '@/components/Renderer/YoutubeChat'

const COMMAND_HEARTBEAT = 0
const COMMAND_JOIN_ROOM = 1
const COMMAND_ADD_TEXT = 2
const COMMAND_ADD_GIFT = 3
//const COMMAND_ADD_MEMBER = 4
//const COMMAND_ADD_SUPER_CHAT = 5
//const COMMAND_DEL_SUPER_CHAT = 6
//const COMMAND_UPDATE_TRANSLATION = 7
const COMMAND_ADD_LOVE = 8
const COMMAND_QUIT_ROOM = 9
const COMMAND_ADD_FOLLOW = 10
const COMMAND_ADD_JOIN_GROUP = 11

export default {
  name: 'Room',
  components: {
    ChatRoom,
    FullScreen,
    YoutubeChat
  },
  data() {
    return {
      //渲染器选择
      "currentRenderer": "YoutubeChat",

      //WS链接
      "webSocketConnection": null,
      "timer": {
        "heartbeat": null
      },
      "status": {
        "server": 0,
        "client": 0,
        "lost": 0,
        "retry": 0,
        "destroy": false,
      },

      //配置文件档
      config: { ...chatConfig.DEFAULT_CONFIG },
    }
  },
  beforeDestroy() {
    this.status.destroy = true
    this.webSocketConnection.close()
  },
  created() {
    console.log("DanmakuNeko Started")
    if (this.$route.query.renderer != "") {
      this.currentRenderer = this.$route.query.renderer
    }
    //可能还是需要兼容老样式，方便调试
    if (this.$route.query.config != "") {
      //var b64 = Base64.decode(link)
      //this.jsonstr = JSON.parse(b64);
    }
    this.config.roomID = this.$route.params.roomId
    //处理WS链接
    this.linkStart()
  },
  methods: {
    sendHeartbeat() {
      if (this.webSocketConnection.readyState === 1) {
        this.webSocketConnection.send(JSON.stringify({
          cmd: COMMAND_HEARTBEAT
        }))
      }

      this.status.client = Date.now()

      if (this.status.client - this.status.server > 2 * 1000) {
        window.console.log(`无心跳 ${++this.status.lost}`)
      } else {
        this.status.lost = 0
      }

      if (this.status.lost > 2) {
        window.console.log(`无心跳重连`)
        this.status.lost = 0
        this.webSocketConnection.close()
      }
    },
    linkStart() {
      const url = `wss://danmaku.loli.ren/chat` //A（稳定运行，感谢Orzogc的后端支持）
      //const url = `wss://danmu.loli.ren/chat` //B（修复中）
      this.webSocketConnection = new WebSocket(url)
      this.webSocketConnection.onopen = this.onWsOpen
      this.webSocketConnection.onclose = this.onWsClose
      this.webSocketConnection.onmessage = this.onWsMessage
      this.timer.heartbeat = window.setInterval(this.sendHeartbeat, 1 * 1000)
    },
    onWsOpen() {
      this.status.retry = 0
      this.status.server = Date.now()
      this.webSocketConnection.send(JSON.stringify({
        cmd: COMMAND_JOIN_ROOM,
        data: {
          roomId: parseInt(this.$route.params.roomId),
          version: "1.0.0",
          isfirstLoad: true,
          config: {
            autoTranslate: false
          }
        }
      }))
    },
    onWsClose() {
      if (this.status.destroy) {
        return
      }

      if (this.timer.heartbeat) {
        window.clearInterval(this.timer.heartbeat)
        this.timer.heartbeat = null
      }

      if (this.isDestroying) {
        return
      }
      window.console.log(`掉线重连中: ${++this.status.retry}`)
      this.linkStart()
    },
    //此处进行优先处理
    onWsMessage(event) {
      let { cmd, data } = JSON.parse(event.data)
      //console.log(data)
      //console.log(cmd)
      let message = null
      //优先进行处理后保存至store，再由渲染器处理渲染
      switch (cmd) {
        case COMMAND_HEARTBEAT:
          this.status.server = Date.now()
          break
        case COMMAND_JOIN_ROOM:
          if (!this.config.danmaku.showJoin || this.mergeSimilar(data)) {
            break
          }
          message = data
          break
        case COMMAND_QUIT_ROOM:
          if (!this.config.danmaku.showQuit || this.mergeSimilar(data)) {
            break
          }
          message = data
          break
        case COMMAND_ADD_TEXT:
          if (!this.config.danmaku.showDanmaku || !this.messageFilter(data) || this.mergeSimilar(data)) {
            break
          }
          message = data
          break
        case COMMAND_ADD_GIFT:
          /*if (!this.config.showGift) {
            break
          }
          let price = (data.totalCoin / 1000)
          if (this.mergeSimilarGift(data.authorName, price, data.giftName, data.num)) {
            break
          }
          if (price < this.config.minGiftPrice) { // 丢人
            break
          }*/
          break
        case COMMAND_ADD_LOVE:
          if (!this.config.danmaku.showStar || this.mergeSimilar(data)) {
            break
          }
          message = data
          break
        case COMMAND_ADD_FOLLOW:
          if (!this.config.danmaku.showFollow || this.mergeSimilar(data)) {
            break
          }
          message = data
          break
        case COMMAND_ADD_JOIN_GROUP:
          if (!this.config.danmaku.showJoinGroup) {
            break
          }
          message = data
          break
      }
      if (message) {
        message.type = cmd
        this.$store.state.roomInfo.danmakuList.unshift(message)
      }
    },
    //否 要推，是 不要推
    mergeSimilar(data) {
      return false
    },
    messageFilter(data) {
      return true
    }
  }
}
</script>