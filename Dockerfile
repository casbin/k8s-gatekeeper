FROM golang:1.18 as builder
WORKDIR /webhook
COPY ./ /webhook
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o webhook cmd/webhook/main.go

FROM debian:latest as webhook
WORKDIR /workspace
COPY --from=builder /webhook/webhook .
COPY --from=builder /webhook/config ./config
CMD cd /workspace && ./webhook --externalClient=false


