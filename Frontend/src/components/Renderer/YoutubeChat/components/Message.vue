<template>
  <yt-live-chat-text-message-renderer>
    <img-shadow id="author-photo" height="24" width="24" class="style-scope yt-live-chat-text-message-renderer"
      :imgUrl="danmaku.avatarUrl"></img-shadow>
    <div id="content" class="style-scope yt-live-chat-text-message-renderer">
      <span id="timestamp" class="style-scope yt-live-chat-text-message-renderer">{{timeText}}</span>
      <span v-if="getShowMedal" :id="`medalName-${medalDisplayColorLV}`"
        class="style-scope yt-live-chat-text-message-renderer">{{danmaku.medalInfo.ClubName}}</span>
      <span v-if="getShowMedal" :id="`medalLevel-${medalDisplayColorLV}`"
        class="style-scope yt-live-chat-text-message-renderer">({{danmaku.medalInfo.Level}})</span>
      <span v-if="danmaku.userMark !== ''" id="usermark"
        class="style-scope yt-live-chat-text-message-renderer">{{danmaku.userMark}}</span>
      <yt-live-chat-author-chip class="style-scope yt-live-chat-text-message-renderer">
        <span id="author-name" dir="auto" class="style-scope yt-live-chat-author-chip" :type="authorTypeText">{{
          danmaku.authorName
          }}
          <!-- 这里是已验证勋章 -->
          <span id="chip-badges" class="style-scope yt-live-chat-author-chip"></span>
        </span>
        <span id="chat-badges" class="style-scope yt-live-chat-author-chip">
          <author-badge class="style-scope yt-live-chat-author-chip" :isAdmin="danmaku.authorType === 2"
            :privilegeType="danmaku.privilegeType"></author-badge>
        </span>
      </yt-live-chat-author-chip>
      <span id="message" class="style-scope yt-live-chat-text-message-renderer">{{
        danmaku.content
        }}<div :value="repeated" :max="99" v-show="repeated > 1" class="style-scope yt-live-chat-text-message-renderer"
          :style="{'--repeated-mark-color': repeatedMarkColor}">
        </div>
      </span>
    </div>
  </yt-live-chat-text-message-renderer>
</template>

<script>
import ImgShadow from './ImgShadow.vue'
import AuthorBadge from './AuthorBadge.vue'
import * as constants from './constants'

// HSL
const REPEATED_MARK_COLOR_START = [210, 100.0, 62.5]
const REPEATED_MARK_COLOR_END = [360, 87.3, 69.2]

export default {
  name: 'Message',
  components: {
    ImgShadow,
    AuthorBadge
  },
  props: {
    danmaku: Object,
    repeated: Number,
    medalDisplayColorLV: {
      type: Number,
      default: 1
    },
    getShowMedal: {
      type: Boolean,
      default: false
    },
  },
  data() {
    return {
    }
  },
  created() {
  },
  computed: {
    timeText() {
      let date = new Date(this.danmaku.timestamp * 1000)
      let hour = date.getHours()
      let min = ('00' + date.getMinutes()).slice(-2)
      return `${hour}:${min}`
    },
    authorTypeText() {
      return constants.AUTHOR_TYPE_TO_TEXT[this.authorType]
    },
    repeatedMarkColor() {
      let color
      if (this.repeated <= 2) {
        color = REPEATED_MARK_COLOR_START
      } else if (this.repeated >= 10) {
        color = REPEATED_MARK_COLOR_END
      } else {
        color = [0, 0, 0]
        let t = (this.repeated - 2) / (10 - 2)
        for (let i = 0; i < 3; i++) {
          color[i] = REPEATED_MARK_COLOR_START[i] + (REPEATED_MARK_COLOR_END[i] - REPEATED_MARK_COLOR_START[i]) * t
        }
      }
      return `hsl(${color[0]}, ${color[1]}%, ${color[2]}%)`
    }
  }
}
</script>


<style>
yt-live-chat-text-message-renderer > #content > #message > .el-badge {
  margin-left: 10px;
}

yt-live-chat-text-message-renderer
  > #content
  > #message
  > .el-badge
  .el-badge__content {
  font-size: 12px !important;
  line-height: 18px !important;
  text-shadow: none !important;
  font-family: sans-serif !important;
  color: #fff !important;
  background-color: var(--repeated-mark-color) !important;
  border: none;
}
</style>

<style src="../assets/css/youtube/yt-live-chat-text-message-renderer.css"></style>
<style src="../assets/css/youtube/yt-live-chat-author-chip.css"></style>