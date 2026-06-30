FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    libpq-dev \
    libssl-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY src ./src

RUN gcc -std=c11 -Wall -Wextra -Wpedantic -O2 \
  -I/usr/include/postgresql \
  src/main.c \
  -o /usr/local/bin/sealion-app \
  -lpq -lcrypto

EXPOSE 8080

CMD ["sealion-app"]

