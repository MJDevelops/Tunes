import { createApp } from "vue";
import router from "./router";
import App from "./App.vue";
import vuetify from "./vuetify";

createApp(App).use(vuetify).use(router).mount("#app");
