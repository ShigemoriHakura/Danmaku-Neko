<template>
  <yt-live-chat-renderer class="style-scope yt-live-chat-app" style="--scrollbar-width:11px;" hide-timestamps
    hide-medal>
    <Ticker :messages="tickerMessage"></Ticker>
    <yt-live-chat-item-list-renderer class="style-scope animated">
      <component v-for="item in danmakuFrame" :key="item.id" v-bind:is="item.component" :danmaku="item.danmaku"
        :getShowMedal="getShowMedal(item.medalInfo)" :medalDisplayColorLV="getMedalLVType(item.medalInfo)">
      </component>
    </yt-live-chat-item-list-renderer>
  </yt-live-chat-renderer>
</template>

<script>
//后面需要重新处理平滑

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

import Gift from './components/Gift'
import Message from './components/Message'
import Ticker from './components/Ticker'

export default {
  name: 'YoutubeChat',
  components: {
    Gift,
    Message,
    Ticker
  },
  data() {
    return {
      danmakuFrame: [
      ],
      tickerMessage: [
      ]
    }
  },
  props: {
    config: Object,
  },
  created() {
    this.size = this.getViewportSize()
    console.log("DanmakuNeko - YoutubeChat: Started at X: " + this.size.width + " ,Y: " + this.size.height + " !")
    this.$store.watch((state) => state.roomInfo.danmakuList, async (newValue, oldValue) => {
      if (newValue.length > 0) {
        var danmaku = newValue[0]
        switch (danmaku.type) {
          //选择对应渲染器渲染
          case COMMAND_JOIN_ROOM:
          case COMMAND_ADD_TEXT:
          case COMMAND_ADD_LOVE:
          case COMMAND_QUIT_ROOM:
          case COMMAND_ADD_FOLLOW:
          case COMMAND_ADD_JOIN_GROUP:
            this.danmakuFrame.push({
              "danmaku": danmaku,
              "component": "Message"
            })
            break
          case COMMAND_ADD_GIFT:
            this.danmakuFrame.push({
              "danmaku": danmaku,
              "component": "Gift"
            })
            break
        }
        console.log(danmaku)
      }
    })

    this.processDanmaku()
  },
  methods: {
    getShowMedal(medal) {
      if (medal && medal.ClubName != '') {
        if (this.showEqualMedal) {
          if (medal.UperID == this.config.roomID) {
            return true;
          }
        } else {
          return true;
        }
      }
      return false;
    },
    getMedalLVType(medal) {
      if (medal && medal.Level) {
        if (medal.Level <= 3) {
          return 1;
        } else if (medal.Level > 3 && medal.Level <= 6) {
          return 2;
        } else if (medal.Level > 6 && medal.Level <= 9) {
          return 3;
        } else if (medal.Level > 9 && medal.Level <= 12) {
          return 4;
        } else if (medal.Level > 12 && medal.Level) {
          return 5;
        }
      }
      return 0;
    },
    //处理弹幕消失，ticker计算，处理文字显示转换，总之就是核心函数！需要做trycatch处理并且上报错误
    processDanmaku() {

    }
  }
}
</script>



<style src="./assets/css/youtube/yt-html.css"></style>
<style src="./assets/css/youtube/yt-live-chat-renderer.css"></style>
<style src="./assets/css/youtube/yt-live-chat-item-list-renderer.css"></style>