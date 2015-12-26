FROM golang

ENV GOPATH /app/gopath
ADD . /wkdir/
WORKDIR /wkdir/svr

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go get github.com/everfore/rpcsv
RUN go build main.go

EXPOSE 8800

CMD ["/wkdir/svr/svr"]
