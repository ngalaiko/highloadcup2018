FROM golang:1.11.3-alpine as builder

WORKDIR ${GOPATH}/src/github.com/ngalayko/highloadcup/tester

COPY . .

RUN go build -o /tester ./cmd/tester/main.go

FROM alpine:3.8

COPY --from=builder /tester /tester

ENTRYPOINT [ "/tester" ]
