import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";
import { stat, readFile } from "node:fs/promises";
import { extname, join, normalize, sep } from "node:path";
import { docsResponseHeaders, rewriteDocsHtml } from "./component/l3";

const port = Number(process.env.PORT || 8080);
const siteRoot = join(import.meta.dir, "..", "site");
const apiURL = process.env.API_URL || "http://api:8080";

const contentTypes: Record<string, string> = {
  ".css": "text/css; charset=utf-8",
  ".html": "text/html; charset=utf-8",
  ".ico": "image/x-icon",
  ".js": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".map": "application/json; charset=utf-8",
  ".md": "text/markdown; charset=utf-8",
  ".svg": "image/svg+xml",
  ".txt": "text/plain; charset=utf-8",
};

const versionedAssetPaths = new Map<string, string>([
  ["assets/intro.js", versionedAssetPath("assets/intro.js")],
  ["assets/styles.css", versionedAssetPath("assets/styles.css")],
]);

const routeAliases: Record<string, string> = {
  "/initial-user-experience": "/create-your-first-app",
  "/initial-user-experience.html": "/create-your-first-app",
};

function sitePath(pathname: string): string | null {
  let requestPath;
  try {
    requestPath = decodeURIComponent(pathname);
  } catch {
    return null;
  }
  if (requestPath === "/for/agents") requestPath = "/for/agents/index.md";
  if (requestPath === "/") requestPath = "/index.html";
  if (!extname(requestPath)) requestPath = `${requestPath}.html`;

  const candidate = normalize(join(siteRoot, requestPath));
  if (!candidate.startsWith(siteRoot + sep) && candidate !== siteRoot) {
    return null;
  }
  return candidate;
}

function canonicalDocsPath(pathname: string): string {
  if (routeAliases[pathname]) {
    return routeAliases[pathname];
  }
  if (pathname === "/index" || pathname === "/index.html") {
    return "/";
  }
  if (pathname === "/for/agents/index.md") {
    return "/for/agents";
  }
  if (pathname.endsWith(".html")) {
    return pathname.slice(0, -".html".length) || "/";
  }
  return "";
}

function redirectToCanonical(request: Request, pathname: string): Response {
  const target = new URL(request.url);
  return new Response(null, {
    status: 308,
    headers: {
      location: `${pathname}${target.search}`,
    },
  });
}

async function proxy(request: Request, pathname: string): Promise<Response> {
  const upstream = new URL(pathname, apiURL);
  upstream.search = new URL(request.url).search;

  return fetch(upstream, {
    method: request.method,
    headers: request.headers,
    body: request.body,
  });
}

async function serveStatic(pathname: string): Promise<Response> {
  const path = sitePath(pathname);
  if (!path) return new Response("not found", { status: 404 });

  try {
    const info = await stat(path);
    if (!info.isFile()) return new Response("not found", { status: 404 });
    const type = contentTypes[extname(path)] || "application/octet-stream";
    const rawBody = await readFile(path);
    const body = type.startsWith("text/html")
      ? cacheBustHtml(rewriteDocsHtml(rawBody.toString("utf8")))
      : rawBody;
    const cache = cacheControlFor(pathname);
    return new Response(body, {
      headers: docsResponseHeaders({ cache, type }),
    });
  } catch (error) {
    if (error && typeof error === "object" && "code" in error && error.code === "ENOENT") {
      return new Response("not found", { status: 404 });
    }
    throw error;
  }
}

function cacheBustHtml(html: string): string {
  let output = html;
  for (const [assetPath, versionedPath] of versionedAssetPaths) {
    output = output.replaceAll(`"${assetPath}"`, `"${versionedPath}"`);
    output = output.replaceAll(`"/${assetPath}"`, `"/${versionedPath}"`);
  }
  return output;
}

function cacheControlFor(pathname: string): string {
  if (pathname === "/assets/intro.js" || pathname === "/assets/styles.css") {
    return "no-cache";
  }
  return "no-store";
}

function versionedAssetPath(assetPath: string): string {
  try {
    const content = readFileSync(join(siteRoot, assetPath));
    const hash = createHash("sha256").update(content).digest("hex").slice(0, 12);
    return `${assetPath}?v=${hash}`;
  } catch (_error) {
    return assetPath;
  }
}

Bun.serve({
  port,
  async fetch(request) {
    const url = new URL(request.url);
    if (url.pathname === "/health" || url.pathname.startsWith("/api/")) {
      return proxy(request, url.pathname);
    }
    const canonicalPath = canonicalDocsPath(url.pathname);
    if (canonicalPath) {
      return redirectToCanonical(request, canonicalPath);
    }
    const docsResponse = await serveStatic(url.pathname);
    if (docsResponse.status !== 404) {
      return docsResponse;
    }
    return docsResponse;
  },
});

console.log(`Carbide docs web listening on container port ${port}`);
console.log(`public URL is ${process.env.PUBLIC_URL || `http://localhost:${port}`}`);
