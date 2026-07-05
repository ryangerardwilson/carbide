import { stat, readFile } from "node:fs/promises";
import { extname, join, normalize, sep } from "node:path";

const port = Number(process.env.PORT || 8080);
const siteRoot = join(import.meta.dir, "site");
const apiURL = process.env.API_URL || "http://api:8080";

const contentTypes = {
  ".css": "text/css; charset=utf-8",
  ".html": "text/html; charset=utf-8",
  ".ico": "image/x-icon",
  ".js": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".svg": "image/svg+xml",
  ".txt": "text/plain; charset=utf-8",
};

const routeAliases = {
  "/initial-user-experience": "/create-your-first-app",
  "/initial-user-experience.html": "/create-your-first-app",
};

function sitePath(pathname) {
  let requestPath = decodeURIComponent(pathname);
  if (requestPath === "/") requestPath = "/index.html";
  if (!extname(requestPath)) requestPath = `${requestPath}.html`;

  const candidate = normalize(join(siteRoot, requestPath));
  if (!candidate.startsWith(siteRoot + sep) && candidate !== siteRoot) {
    return null;
  }
  return candidate;
}

function canonicalDocsPath(pathname) {
  if (routeAliases[pathname]) {
    return routeAliases[pathname];
  }
  if (pathname === "/index" || pathname === "/index.html") {
    return "/";
  }
  if (pathname.endsWith(".html")) {
    return pathname.slice(0, -".html".length) || "/";
  }
  return "";
}

function redirectToCanonical(request, pathname) {
  const target = new URL(request.url);
  return new Response(null, {
    status: 308,
    headers: {
      location: `${pathname}${target.search}`,
    },
  });
}

async function proxy(request, pathname) {
  const upstream = new URL(pathname, apiURL);
  upstream.search = new URL(request.url).search;

  return fetch(upstream, {
    method: request.method,
    headers: request.headers,
    body: request.body,
  });
}

async function serveStatic(pathname) {
  const path = sitePath(pathname);
  if (!path) return new Response("not found", { status: 404 });

  try {
    const info = await stat(path);
    if (!info.isFile()) return new Response("not found", { status: 404 });
    const body = await readFile(path);
    const type = contentTypes[extname(path)] || "application/octet-stream";
    const cache = pathname.startsWith("/assets/")
      ? "public, max-age=31536000, immutable"
      : "no-cache";
    return new Response(body, {
      headers: {
        "cache-control": cache,
        "content-type": type,
      },
    });
  } catch (error) {
    if (error && error.code === "ENOENT") {
      return new Response("not found", { status: 404 });
    }
    throw error;
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
    return serveStatic(url.pathname);
  },
});

console.log(`Carbide docs web listening on container port ${port}`);
console.log(`public URL is ${process.env.PUBLIC_URL || `http://localhost:${port}`}`);
