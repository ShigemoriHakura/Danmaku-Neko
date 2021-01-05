import Vue from "vue";
import VueRouter from "vue-router";
import VueI18n from "vue-i18n";
import Vuex from "vuex";
import App from "./App.vue";

import zh from "./lang/zh";
import en from "./lang/en";

import Room from "./views/Room.vue";
import NotFound from "./views/NotFound.vue";

Vue.config.productionTip = false;

Vue.config.ignoredElements = [
  /^yt-/
]

Vue.prototype.getViewportSize = function () {
  return {
    width:
      window.innerWidth ||
      document.documentElement.clientWidth ||
      document.body.clientWidth,
    height:
      window.innerHeight ||
      document.documentElement.clientHeight ||
      document.body.clientHeight,
  };
};

Vue.use(VueRouter);
Vue.use(VueI18n);
Vue.use(Vuex);

const router = new VueRouter({
  mode: "history",
  routes: [
    { path: "/room/:roomId", name: "room", component: Room },
    { path: "*", component: NotFound },
  ],
});

const store = new Vuex.Store({
  state: {
    roomInfo: {
      danmakuList: [],
    },
  },
});

let locale = window.localStorage.lang;
if (!locale) {
  let lang = navigator.language;
  if (lang.startsWith("zh")) {
    locale = "zh";
  } else {
    locale = "en";
  }
}

const i18n = new VueI18n({
  locale,
  fallbackLocale: "en",
  messages: {
    zh,
    en,
  },
});

new Vue({
  store: store,
  router,
  i18n,
  render: (h) => h(App),
}).$mount("#app");
