FROM oven/bun:1.3.14-debian

WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY index.html ./
COPY src ./src
RUN bun run tailwind:build

EXPOSE 8080

CMD ["bun", "run", "dev"]
