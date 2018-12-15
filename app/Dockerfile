FROM golang:1.11.3-alpine as builder

WORKDIR ${GOPATH}/src/github.com/ngalayko/highloadcup/app

COPY . .

RUN go build -o /highload ./cmd/highloadcup/main.go

FROM alpine:3.8

COPY --from=builder /highload /highload

ENTRYPOINT [ "/highload", "--addr=:80", "--data_path=/tmp/data/data.zip" ]
