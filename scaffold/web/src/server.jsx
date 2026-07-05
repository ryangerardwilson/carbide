import { stat } from 'node:fs/promises';
import { extname, join, normalize, sep } from 'node:path';

const port = Number(Bun.env.FRONTEND_PORT || Bun.env.PORT || 8080);
const apiUrl = Bun.env.API_URL || 'http://api:8080';
const publicUrl = Bun.env.PUBLIC_URL || '';
const apiOrigin = new URL(apiUrl);
const publicRoot = join(import.meta.dir, '..', 'public');
const shellRoutes = new Set(['/', '/login', '/register', '/dashboard']);

const contentTypes = {
  '.css': 'text/css; charset=utf-8',
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.map': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml'
};

function proxyToAPI(request) {
  const incomingUrl = new URL(request.url);
  const upstreamUrl = new URL(`${incomingUrl.pathname}${incomingUrl.search}`, apiOrigin);
  const headers = new Headers(request.headers);
  headers.delete('host');
  headers.set('x-forwarded-host', incomingUrl.host);
  headers.set('x-forwarded-proto', incomingUrl.protocol.replace(':', ''));

  const options = {
    method: request.method,
    headers,
    redirect: 'manual'
  };

  if (request.method !== 'GET' && request.method !== 'HEAD') {
    options.body = request.body;
  }

  return fetch(upstreamUrl, options);
}

function safePublicPath(pathname) {
  let decoded;
  try {
    decoded = decodeURIComponent(pathname);
  } catch {
    return null;
  }

  const candidate = normalize(join(publicRoot, decoded));
  if (candidate !== publicRoot && !candidate.startsWith(publicRoot + sep)) {
    return null;
  }
  return candidate;
}

function cacheControlFor(pathname) {
  if (pathname.startsWith('/assets/')) {
    return 'public, max-age=31536000, immutable';
  }
  return 'no-store';
}

async function servePublicFile(request, pathname) {
  const path = safePublicPath(pathname);
  if (!path) {
    return new Response('Not found', {
      status: 404,
      headers: { 'Content-Type': 'text/plain; charset=utf-8' }
    });
  }

  try {
    const info = await stat(path);
    if (!info.isFile()) {
      throw new Error('not a file');
    }
  } catch {
    return new Response('Not found', {
      status: 404,
      headers: { 'Content-Type': 'text/plain; charset=utf-8' }
    });
  }

  const headers = {
    'Cache-Control': cacheControlFor(pathname),
    'Content-Type': contentTypes[extname(path)] || 'application/octet-stream'
  };
  return new Response(request.method === 'HEAD' ? null : Bun.file(path), { headers });
}

const server = Bun.serve({
  port,
  development: Bun.env.NODE_ENV !== 'production',
  async fetch(request) {
    const url = new URL(request.url);
    if (url.pathname === '/api' || url.pathname.startsWith('/api/')) {
      return proxyToAPI(request);
    }
    if (url.pathname === '/health') {
      return proxyToAPI(request);
    }
    if (shellRoutes.has(url.pathname)) {
      return servePublicFile(request, '/index.html');
    }
    if (url.pathname.startsWith('/assets/') || url.pathname === '/asset-manifest.json') {
      return servePublicFile(request, url.pathname);
    }
    return new Response('Not found', {
      status: 404,
      headers: { 'Content-Type': 'text/plain; charset=utf-8' }
    });
  }
});

console.log(`Carbide Bun frontend listening inside container on :${server.port}`);
if (publicUrl) {
  console.log(`browser entrypoint ${publicUrl}`);
}
console.log(`proxying /api and /health to api service ${apiUrl}`);
