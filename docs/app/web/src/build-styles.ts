import { existsSync, mkdirSync } from "node:fs";
import { dirname, join } from "node:path";

const outputs = [
  join(process.cwd(), "..", "..", "site", "assets", "styles.css"),
  join(process.cwd(), "site", "assets", "styles.css"),
];

const output = outputs.find((candidate) => existsSync(dirname(candidate))) || outputs[0];
if (!output) {
  throw new Error("missing Tailwind output path");
}
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
