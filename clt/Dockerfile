FROM golang

ENV GOPATH /app/gopath
WORKDIR /app/gopath
COPY . /app/gopath/

RUN go get github.com/everfore/rpcsv
RUN go get github.com/everfore/rpcsv/clt
RUN go build -o clnt ./rpcsv/clt/main.go

EXPOSE 80

CMD ["/app/gopath/clnt"]