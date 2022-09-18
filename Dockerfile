FROM golang:alpine as builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /go/src/gitlab.com/httpserver/

RUN mkdir -p bin/amd64

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/httpserver .

FROM ubuntu as prod

LABEL multi.language="go" multi.usage="test"

COPY --from=builder /go/src/gitlab.com/httpserver/bin/amd64/httpserver .

EXPOSE 8080

CMD ["./httpserver"]