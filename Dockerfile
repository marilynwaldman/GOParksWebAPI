FROM golang:1.7.3-alpine

ENV SOURCES /go/src/github.com/user/GOParksWebAPI/
COPY . ${SOURCES}

WORKDIR ${SOURCES}
RUN go build ./server.go

CMD ["./server"]

EXPOSE 8080
