import vue from "@vitejs/plugin-vue";
import { defineConfig } from "vite";
import compression from "vite-plugin-compression2";

export default defineConfig({
  plugins: [
    vue(),
    compression({
      algorithm: "gzip",
      threshold: 1024,
    }),
  ],
  resolve: {
    alias: {
      "@/": "/src/",
    },
  },
  build: {
    outDir: "internal/server/dist",
    emptyOutDir: true,
    minify: "terser",
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {},
    },
    chunkSizeWarningLimit: 1500,
  },
  server: {
    port: 3000,
    host: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
