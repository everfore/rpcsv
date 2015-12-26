FROM golang

ENV GOPATH /app/gopath
WORKDIR /svr

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go get github.com/everfore/rpcsv/svr

EXPOSE 8800

CMD ["/app/gopath/bin/svr"]
