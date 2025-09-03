import { createApp } from "vue";
import router from "./router";
import App from "./App.vue";
import PrimeVue from "primevue/config";

createApp(App).use(PrimeVue).use(router).mount("#app");
