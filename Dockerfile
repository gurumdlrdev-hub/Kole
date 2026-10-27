# ==========================================
# Stage 1: Build Environment
# ==========================================
FROM golang:bookworm AS builder


ENV GOTOOLCHAIN=auto


RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    libx11-dev \
    xorg-dev \
    libxtst-dev \
    libpng-dev \
    libxinerama-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxi-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/kole-server ./cmd/server

# ==========================================
# Stage 2: Minimal Runtime Environment
# ==========================================
FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y --no-install-recommends \
    libx11-6 \
    libxtst6 \
    libpng16-16 \
    libxinerama1 \
    libxcursor1 \
    libxrandr2 \
    libxi6 \
    xvfb \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY --from=builder /app/kole-server /app/kole-server


EXPOSE 9000/udp


CMD ["sh", "-c", "Xvfb :99 -screen 0 1920x1080x24 & export DISPLAY=:99 && ./kole-server"]