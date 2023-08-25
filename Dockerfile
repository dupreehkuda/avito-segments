FROM golang:alpine AS builder

WORKDIR /build

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "$LDFLAGS" \
    -v -o ./build ./cmd/app/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build /build

ENTRYPOINT ["./build"]