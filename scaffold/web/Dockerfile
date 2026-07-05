FROM oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f

WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY index.html ./
COPY src ./src
COPY tsconfig.json ./
RUN bun run typecheck && bun run assets:build

EXPOSE 8080

CMD ["bun", "run", "dev"]
