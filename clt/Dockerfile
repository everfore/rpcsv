FROM golang

COPY . /app/gopath/
ENV GOPATH /app/gopath
WORKDIR /app/gopath/clt

RUN go get github.com/everfore/rpcsv
RUN go get github.com/everfore/rpcsv/clt
RUN go build -o clt main.go

EXPOSE 80

CMD ["/app/gopath/clt/clt"]
