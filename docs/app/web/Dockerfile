FROM oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f

WORKDIR /app
COPY app/web/package.json app/web/bun.lock ./
RUN bun install --frozen-lockfile

COPY app/web/src ./src
COPY app/web/tsconfig.json ./
COPY site ./site
RUN bun run typecheck && bun run tailwind:build

EXPOSE 8080
CMD ["bun", "run", "start"]
