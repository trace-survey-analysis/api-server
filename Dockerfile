FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

WORKDIR /app

# Install make
RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ internal/
COPY cmd/ cmd/
COPY Makefile ./
ARG TARGETOS TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/api-server ./bin/api-server

EXPOSE 8080

CMD ["./bin/api-server"]