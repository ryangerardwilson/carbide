import { readdir, readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const root = join(import.meta.dir, '..');
const publicRoot = join(root, 'public');
const assetsRoot = join(publicRoot, 'assets');

const files = (await readdir(assetsRoot)).sort();
const scripts = files.filter((file) => /^main-[a-z0-9]+\.js$/.test(file));
const stylesheets = files.filter((file) => /^main-[a-z0-9]+\.css$/.test(file));

if (scripts.length !== 1) {
  throw new Error(`expected one hashed main script, found ${scripts.length}`);
}

let html = await readFile(join(root, 'index.html'), 'utf8');
const styleTags = stylesheets.map((file) => `    <link rel="stylesheet" href="/assets/${file}">`).join('\n');
const scriptTag = `    <script type="module" src="/assets/${scripts[0]}"></script>`;
const replacement = [styleTags, scriptTag].filter(Boolean).join('\n');

html = html.replace(/\s*<script type="module" src="\.\/src\/main\.jsx"><\/script>/, `\n${replacement}`);
if (!html.includes(`/assets/${scripts[0]}`)) {
  throw new Error('failed to write hashed browser asset references');
}

await writeFile(join(publicRoot, 'index.html'), html);
await writeFile(
  join(publicRoot, 'asset-manifest.json'),
  `${JSON.stringify(
    {
      entry: `/assets/${scripts[0]}`,
      stylesheets: stylesheets.map((file) => `/assets/${file}`)
    },
    null,
    2
  )}\n`
);
