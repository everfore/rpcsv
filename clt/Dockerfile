FROM golang

COPY . /app/gopath/
WORKDIR /app/gopath/clt
ENV GOPATH /app/gopath

RUN go get github.com/everfore/rpcsv
RUN go get github.com/everfore/rpcsv/clt
#RUN go build -o clnt main.go

EXPOSE 80

CMD ["/app/gopath/bin/clt"]
