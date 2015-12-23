FROM golang

WORKDIR /app/gopath
ADD . /app/gopath/
ENV GOPATH /app/gopath

RUN go get github.com/everfore/rpcsv
RUN go get github.com/everfore/rpcsv/clt
RUN go build -o clnt ./clt/main.go

EXPOSE 80

CMD ["/app/gopath/clnt"]
