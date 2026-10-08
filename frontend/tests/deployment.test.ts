// @vitest-environment node
import { afterEach, expect, it } from "vitest";
import { spawn, type ChildProcess } from "node:child_process";
import { mkdtemp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { createRequire } from "node:module";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { once } from "node:events";

let directory = "";
let server: ChildProcess | undefined;
afterEach(async () => {
  if (server && server.exitCode === null) {
    const exited = once(server, "exit");
    server.kill();
    await exited;
  }
  if (directory) await rm(directory, { recursive: true, force: true });
});

it("本番の起動コマンドはPORTに従い、SPAと静的ファイルをViteなしで配信する", async () => {
  const manifest = JSON.parse(await readFile(new URL("../package.json", import.meta.url), "utf8"));
  const [command, ...args] = manifest.scripts.start.split(" ");
  expect(command).toBe("serve");
  expect(manifest.dependencies.serve).toBeTruthy();
  directory = await mkdtemp(join(tmpdir(), "ducks-deploy-"));
  await mkdir(join(directory, "dist"));
  await writeFile(join(directory, "dist/index.html"), '<html lang="ja"><body>Ducks deployment check</body></html>');
  await writeFile(join(directory, "dist/app.js"), "console.log('Ducks');");
  const require = createRequire(import.meta.url);
  const executable = join(dirname(require.resolve("serve/package.json")), "build/main.js");
  let output = "";
  server = spawn(process.execPath, [executable, ...args], {
    cwd: directory,
    env: { ...process.env, NODE_ENV: "production", PORT: "0", NO_UPDATE_CHECK: "1" },
    stdio: ["ignore", "pipe", "pipe"],
  });
  const port = await new Promise<string>((resolve, reject) => {
    server!.on("error", reject);
    server!.on("exit", () => reject(new Error(`Server stopped: ${output}`)));
    server!.stdout!.on("data", chunk => {
      output += String(chunk);
      const address = output.match(/http:\/\/localhost:(\d+)/);
      if (address) resolve(address[1]);
    });
    server!.stderr!.on("data", chunk => { output += String(chunk); });
  });
  const base = `http://127.0.0.1:${port}`;
  for (const path of ["/", "/event/festival-id", "/poster/detail/poster-id", "/sales/cashier"]) {
    const response = await fetch(base + path);
    expect(response.status).toBe(200);
    expect(response.headers.get("content-type")).toContain("text/html");
    expect(await response.text()).toContain("Ducks deployment check");
  }
  const asset = await fetch(`${base}/app.js`);
  expect(asset.headers.get("content-type")).toContain("javascript");
  expect(await asset.text()).toBe("console.log('Ducks');");
  const head = await fetch(`${base}/sales/orders`, { method: "HEAD" });
  expect(head.status).toBe(200);
  expect(await head.text()).toBe("");
}, 10000);
