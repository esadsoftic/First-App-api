FROM golang:onbuild

LABEL maintainer="Esad from go"

RUN mkdir /go/src/firstapp

WORKDIR /go/src/firstapp

ADD . .

RUN go get

RUN go build -o main .

CMD ["./main"]

