FROM golang:1.16-alpine
ADD . /go/src/seabattle
WORKDIR /go/src/seabattle
RUN go mod download
RUN go build -o seabattle seabattle/cmd

EXPOSE 3000

CMD ["/go/src/seabattle/seabattle"]