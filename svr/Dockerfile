FROM golang

COPY . /app/gopath/
WORKDIR /app/gopath/svr
ENV GOPATH /app/gopath

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go get github.com/everfore/rpcsv/svr
#RUN go build -o rpcsv main.go

EXPOSE 8800

CMD ["/app/gopath/bin/svr"]
