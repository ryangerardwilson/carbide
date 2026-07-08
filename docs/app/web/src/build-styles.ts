import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";

const output = join(process.cwd(), "site", "assets", "styles.css");
mkdirSync(dirname(output), { recursive: true });

const result = Bun.spawnSync([
  "tailwindcss",
  "-i",
  "./src/styles.css",
  "-o",
  output,
  "--minify",
], {
  stdout: "inherit",
  stderr: "inherit",
});

process.exit(result.exitCode ?? 1);
