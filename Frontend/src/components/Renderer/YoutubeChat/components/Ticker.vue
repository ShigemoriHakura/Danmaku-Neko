<template>
  <yt-live-chat-ticker-renderer :hidden="messages.length === 0">
    <div id="container" dir="ltr" class="style-scope yt-live-chat-ticker-renderer">
      <div id="items" class="style-scope yt-live-chat-ticker-renderer">
        <yt-live-chat-ticker-paid-message-item-renderer v-for="message in messages" :key="message.raw.id" tabindex="0"
          class="style-scope yt-live-chat-ticker-renderer" style="overflow: hidden;" @click="onItemClick(message.raw)">
          <div id="container" dir="ltr" class="style-scope yt-live-chat-ticker-paid-message-item-renderer" :style="{
            background: message.bgColor,
          }">
            <div id="content" class="style-scope yt-live-chat-ticker-paid-message-item-renderer" :style="{
              color: message.color
            }">
              <img-shadow id="author-photo" height="24" width="24"
                class="style-scope yt-live-chat-ticker-paid-message-item-renderer" :imgUrl="message.raw.avatarUrl">
              </img-shadow>
              <span id="text" dir="ltr"
                class="style-scope yt-live-chat-ticker-paid-message-item-renderer">{{message.text}}</span>
            </div>
          </div>
        </yt-live-chat-ticker-paid-message-item-renderer>
      </div>
    </div>
    <template v-if="pinnedMessage">
      <membership-item :key="pinnedMessage.id" v-if="pinnedMessage.type === MESSAGE_TYPE_MEMBER"
        class="style-scope yt-live-chat-ticker-renderer" :avatarUrl="pinnedMessage.avatarUrl"
        :authorName="pinnedMessage.authorName" :privilegeType="pinnedMessage.privilegeType" :title="pinnedMessage.title"
        :time="pinnedMessage.time"></membership-item>
      <paid-message :key="pinnedMessage.id" v-else class="style-scope yt-live-chat-ticker-renderer"
        :price="pinnedMessage.price" :avatarUrl="pinnedMessage.avatarUrl" :authorName="pinnedMessage.authorName"
        :time="pinnedMessage.time" :content="pinnedMessageShowContent"></paid-message>
    </template>
  </yt-live-chat-ticker-renderer>
</template>

<script>
export default {
  name: 'Ticker',
  data() {
    return {
      pinnedMessage: null
    }
  },
  props: {
    messages: Array,
  },
  created() {
  },
  methods: {
    onItemClick(message) {
      if (this.pinnedMessage == message) {
        this.pinnedMessage = null
      } else {
        this.pinnedMessage = message
      }
    }
  }
}
</script>
