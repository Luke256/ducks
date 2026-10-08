// @vitest-environment node
import { expect, it } from "vitest";
import { resolveConfig } from "vite";
import viteConfig from "../vite.config";

it("旧環境変数が開発モードでも本番用の設定を生成する", async () => {
  const originalNodeEnv = process.env.NODE_ENV;
  const originalUserEnv = process.env.VITE_USER_NODE_ENV;
  try {
    process.env.NODE_ENV = "development";
    process.env.VITE_USER_NODE_ENV = "development";
    if (typeof viteConfig !== "function")
      throw new Error("設定関数が必要です。");
    const config = await viteConfig({ command: "build", mode: "production" });
    const resolved = await resolveConfig(
      { ...config, configFile: false },
      "build",
      "production",
      "production",
    );
    expect(resolved.isProduction).toBe(true);
    expect(process.env.VITE_USER_NODE_ENV).toBe("");
  } finally {
    if (originalNodeEnv === undefined) delete process.env.NODE_ENV;
    else process.env.NODE_ENV = originalNodeEnv;
    if (originalUserEnv === undefined) delete process.env.VITE_USER_NODE_ENV;
    else process.env.VITE_USER_NODE_ENV = originalUserEnv;
  }
});
