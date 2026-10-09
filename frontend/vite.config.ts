import { fileURLToPath, URL } from "node:url";
import vue from "@vitejs/plugin-vue";
import { loadEnv } from "vite";
import { defineConfig } from "vitest/config";

export default defineConfig(({ command, mode }) => {
    const env = loadEnv(mode, process.cwd(), "");
    // 旧.envのNODE_ENV=developmentを本番ビルドへ持ち込まない。
    if (command === "build") {
        process.env.NODE_ENV = "production";
        process.env.VITE_USER_NODE_ENV = "";
    }
    // 既存の環境変数でも移行直後に同じAPIへ接続できるようにする。
    const apiUrl = env.VITE_API_URL || env.NEXT_PUBLIC_API_URL || "/api/v1";
    return {
        plugins: [vue()],
        resolve: {
            alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
        },
        define: { __API_URL__: JSON.stringify(apiUrl) },
        server: {
            port: 3000,
            proxy: {
                "/api": {
                    target: env.API_PROXY_TARGET || "http://localhost:8080",
                    changeOrigin: true,
                },
            },
        },
        preview: { port: 3000 },
        test: { environment: "jsdom", clearMocks: true },
    };
});
