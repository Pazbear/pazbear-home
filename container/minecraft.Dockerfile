
# 1. Build Image
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY backend ./

RUN go mod download

RUN go build -o minecraft ./cmd/minecraft

WORKDIR /dist

RUN cp /build/minecraft .
RUN cp /build/cmd/minecraft/config/config.conf ./config.conf

# 2. Production Image
FROM scratch

COPY --from=builder /dist/config.conf ./config/config.conf
COPY --from=builder /dist/minecraft .

EXPOSE 32541
#./dbaas-api.init start
ENTRYPOINT ["./minecraft"]