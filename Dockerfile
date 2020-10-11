FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /build

COPY . .
RUN go mod download

RUN go build -o main ./cmd/shorturl-serve


WORKDIR /dist

RUN cp /build/main . 
RUN cp /build/config-docker.json config.json


# Build a small image
FROM scratch
EXPOSE 8080
COPY --from=builder /dist/ /

ENTRYPOINT ["/main"]