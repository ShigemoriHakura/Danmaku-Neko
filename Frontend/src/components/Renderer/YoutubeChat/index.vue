<template>
  <yt-live-chat-renderer class="style-scope yt-live-chat-app" style="--scrollbar-width:11px;" hide-timestamps
    hide-medal>
    Ticker放这里
    <yt-live-chat-item-list-renderer class="style-scope animated">
      <component v-for="item in danmakuFrame" :key="item.id" v-bind:is="item.component" :danmaku="item.danmaku"
        :getShowMedal="getShowMedal(item.medalInfo)" :medalDisplayColorLV="getMedalLVType(item.medalInfo)">
      </component>
    </yt-live-chat-item-list-renderer>
  </yt-live-chat-renderer>
</template>

<script>
//后面需要重新处理平滑

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
      "danmakuFrame": [
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
        this.danmakuFrame.push({
          "danmaku": danmaku,
          "component": "Message"
        })
        console.log(danmaku)
      }
    })
  },
  methods: {
    getShowMedal(medal) {
      if (medal && medal.ClubName != '') {
        if (this.showEqualMedal) {
          if (medal.UperID == this.roomID) {
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
    processDanmaku() {
    }
  }
}
</script>



<style src="./assets/css/youtube/yt-html.css"></style>
<style src="./assets/css/youtube/yt-live-chat-renderer.css"></style>
<style src="./assets/css/youtube/yt-live-chat-item-list-renderer.css"></style>