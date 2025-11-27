import { defineConfig } from "astro/config";
import sitemap from "@astrojs/sitemap";
import vue from "@astrojs/vue";
import tailwind from "@astrojs/tailwind";

export default defineConfig({
  site: "https://tryoneoff.com",
  integrations: [
    vue(),
    tailwind(),
    sitemap({
      filter: (page) => !page.includes("/404"),
      changefreq: "weekly",
      priority: 0.7,
      lastmod: new Date(),
    }),
  ],
  output: "static",
  server: {
    port: 3000,
  },
});
