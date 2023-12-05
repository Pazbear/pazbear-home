# 1. Build Image
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY backend ./

RUN go mod download

RUN go build -o mainapp ./cmd/mainapp

WORKDIR /dist

RUN cp /build/mainapp .
RUN cp /build/cmd/mainapp/config/config.conf ./config.conf

# 2. Production Image
FROM scratch

COPY --from=builder /dist/config.conf ./config/config.conf
COPY --from=builder /dist/mainapp .

EXPOSE 32541
#./dbaas-api.init start
ENTRYPOINT ["./mainapp"]