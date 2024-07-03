#ARG BASE_IMAGE=gcr.io/distroless/static-debian11:latest
ARG BASE_IMAGE=cbioportal/cbioportal:6.0.5

# Build
FROM golang:1.21.0-alpine3.18 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""
ARG APP_VERSION
ARG COMMIT_HASH
ARG GIT_REF
ARG BUILD_DATE
ARG BUILD_BY=docker
ARG GOPROXY
RUN apk add --no-cache --update ca-certificates git upx
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w -X main.version=${APP_VERSION} -X main.commit=${COMMIT_HASH} -X main.date=${BUILD_DATE} -X main.builtBy=${BUILD_BY} -X main.gitRef=${GIT_REF} -linkmode internal -extldflags -static" -o cbioportal-apiserver main.go
# Compress go binary
# https://linux.die.net/man/1/upx
RUN upx -7 -qq ./cbioportal-apiserver && \
    upx -t ./cbioportal-apiserver

# Image
FROM ${BASE_IMAGE}
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Copy module files for CVE scanning / dependency analysis.
COPY --from=builder /app/go.mod /app/go.sum /app/
COPY --from=builder /app/cbioportal-apiserver /app/
ENTRYPOINT ["/app/cbioportal-apiserver"]
