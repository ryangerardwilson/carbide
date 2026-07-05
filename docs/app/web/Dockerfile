FROM oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f

WORKDIR /app
COPY app/web/package.json ./package.json
COPY app/web/server.jsx ./server.jsx
COPY site ./site

EXPOSE 8080
CMD ["bun", "server.jsx"]
