<template>
  <div>
    <div v-for="item in danmakuFrame" :key="item.id" style="position:absolute">
      <p v-bind:style="{ 
        'margin-top':item.top + 'px', 
        'margin-left':item.left + 'px',
        'white-space' : 'nowrap',
        'font-size': '100px',
        //'color': 'white',
        }">
        {{item.content}}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FullScreen',
  data() {
    return {
      "size": null,
      "lastY": 0,
      "danmakuFrame": [
      ]
    }
  },
  created() {
    this.size = this.getViewportSize()
    console.log("DanmakuNeko - FullScreen: Started at X: " + this.size.width + " ,Y: " + this.size.height + " !")
    this.$store.watch((state) => state.roomInfo.danmakuList, async (newValue, oldValue) => {
      if (newValue.length > 0) {
        var danmaku = newValue[0]
        var height = parseInt(Math.random() * (this.size.height + 1), 10);
        console.log(height)
        var width = this.size.width + 10
        this.danmakuFrame.push({
          "top": height,
          "left": width,
          "content": danmaku.content
        })
        console.log(danmaku)
      }
    })

    this.processDanmaku()
  },
  methods: {
    processDanmaku() {
      //console.log(this.danmakuFrame)
      for (var i = 0; i < this.danmakuFrame.length; i++) {
        var element = this.danmakuFrame[i]
        element.left = element.left - 10
        if (element.left < -1000) {
          this.danmakuFrame.splice(i, 1); // 将使后面的元素依次前移，数组长度减1
          i--; // 如果不减，将漏掉一个元素
        }
      }

      setTimeout(this.processDanmaku, 20)
    }
  }
}
</script>
