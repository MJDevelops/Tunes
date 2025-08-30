import { createApp } from "vue";
import router from "./router";
import App from "./App.vue";
import vuetify from "./vuetify";

import "@fontsource/roboto/100.css";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import "@fontsource/roboto/900.css";

createApp(App).use(vuetify).use(router).mount("#app");
