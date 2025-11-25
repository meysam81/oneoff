import { createApp } from "vue";
import { createPinia } from "pinia";
import hljs from "highlight.js/lib/core";
import json from "highlight.js/lib/languages/json";
import bash from "highlight.js/lib/languages/bash";
import router from "./router";
import App from "./App.vue";

hljs.registerLanguage("json", json);
hljs.registerLanguage("bash", bash);

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

app.provide("hljs", hljs);

app.mount("#app");
