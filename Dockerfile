FROM golang

ADD . /go/src/github.com/Betank/go-todo

RUN go install github.com/Betank/go-todo

ENTRYPOINT /go/bin/go-todo

EXPOSE 8080