FROM golang

ENV GOPATH /app/gopath
WORKDIR /app/gopath/clt
COPY ./clt /app/gopath/clt

RUN go get github.com/everfore/rpcsv
RUN go get github.com/everfore/rpcsv/clt
RUN go build -o clt main.go

EXPOSE 80

CMD ["/app/gopath/clt/clt"]