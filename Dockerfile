FROM golang

WORKDIR /app/gopath
ENV GOPATH /app/gopath
ADD . /app/gopath/

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go get github.com/shaalx/goutils
RUN go build -o rpcsv ./svr/main.go

EXPOSE 8800

CMD ["/app/gopath/rpcsv"]