# syntax=docker/dockerfile:1

FROM node:22-alpine AS frontend

ARG USER_WEB_REPO=https://github.com/dujiao-next/user.git
ARG ADMIN_WEB_REPO=https://github.com/dujiao-next/admin.git
ARG USER_WEB_REF=main
ARG ADMIN_WEB_REF=main

WORKDIR /web

RUN apk --no-cache add git \
    && corepack enable \
    && corepack prepare pnpm@10.34.3 --activate

RUN git clone --depth 1 --branch "$USER_WEB_REF" "$USER_WEB_REPO" user
WORKDIR /web/user
RUN pnpm install --frozen-lockfile \
    && pnpm run build

WORKDIR /web
RUN git clone --depth 1 --branch "$ADMIN_WEB_REF" "$ADMIN_WEB_REPO" admin
WORKDIR /web/admin
RUN pnpm install --frozen-lockfile \
    && pnpm run build:fullstack

FROM golang:1.26.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG APP_VERSION=v1.0.0
RUN echo "Building for $TARGETOS/$TARGETARCH$TARGETVARIANT"

WORKDIR /src

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /web/user/dist ./internal/web/dist/user
COPY --from=frontend /web/admin/dist ./internal/web/dist/admin
RUN set -eux; \
    export GOOS="$TARGETOS" GOARCH="$TARGETARCH"; \
    if [ "$TARGETARCH" = "arm" ] && [ -n "$TARGETVARIANT" ]; then export GOARM="${TARGETVARIANT#v}"; fi; \
    if [ "$TARGETARCH" = "amd64" ] && [ -n "$TARGETVARIANT" ]; then export GOAMD64="${TARGETVARIANT#v}"; fi; \
    go build -trimpath -tags "release fullstack" -ldflags="-s -w -X github.com/dujiao-next/internal/version.Version=${APP_VERSION}" -o /out/dujiao-api ./cmd/server

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata \
    && mkdir -p /app/db /app/uploads /app/logs

COPY --from=builder /out/dujiao-api /app/dujiao-api
COPY config.yml.example /app/config.yml.example

EXPOSE 8080

CMD ["./dujiao-api"]
